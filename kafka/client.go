package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.opentelemetry.io/otel"
	"sync"
	"time"
)

type IClient interface {
	Produce(ctx context.Context, data ...any) error
	ProduceByPartitionKey(ctx context.Context, key string, data ...any) error
}
type Client struct {
	types.Component
	config Config
	writer *kafka.Writer
}

func (c *Client) GetKey() string {
	return key
}

type Details struct {
	Hosts []string `json:"hosts"`
	Topic string   `json:"topic"`
}

func (c *Client) Details() any {
	details := Details{}
	details.Hosts = c.config.Hosts
	details.Topic = c.config.Topic
	return details
}

const (
	key = "kafka-producer"
)

func (c *Client) Live(ctx context.Context) error {
	return c.checkKafkaConnectivity(ctx)
}
func (c *Client) Connection() IClient {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c
}

func (c *Client) getDialer() *kafka.Dialer {
	dialer := &kafka.Dialer{
		Timeout:   c.config.DialTimeout,
		DualStack: true,
	}
	if c.config.Login != "" {
		mechanism := &plain.Mechanism{
			Username: c.config.Login, Password: c.config.Password,
		}
		dialer.SASLMechanism = mechanism
	}
	if c.config.PEM != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(c.config.PEM))

		// Create a TLS configuration.
		tlsConfig := &tls.Config{
			RootCAs: caCertPool, // Use the CA certificate pool.
		}
		dialer.TLS = tlsConfig
	}
	return dialer
}
func (c *Client) checkKafkaConnectivity(ctx context.Context) error {
	var errorsGroup []error
	wg := &sync.WaitGroup{}
	wg.Add(len(c.config.Hosts))
	for _, addr := range c.config.Hosts {
		go func(addr string) {
			defer wg.Done()
			_, err := c.getDialer().DialContext(ctx, "tcp", addr)
			if err != nil {
				errorsGroup = append(errorsGroup, err)
			}
		}(addr)
	}
	wg.Wait()
	percent := 1 - float32(len(errorsGroup))/float32(len(c.config.Hosts))
	if percent > 0.6 {
		return nil
	}
	return fmt.Errorf("kafka errors: %s", errorsGroup)
}

func (c *Client) Shutdown() error {
	return c.writer.Close()
}
func (c *Client) Produce(ctx context.Context, data ...any) error {
	return c.innerProducer(ctx, "", data)
}

// HeaderCarrier is a custom type to adapt []kafka.Header to propagation.TextMapCarrier.
type HeaderCarrier []kafka.Header

// Get returns the value for a given key.
func (c *HeaderCarrier) Get(key string) string {
	for _, header := range *c {
		if header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}

// Set sets a key-value pair.
func (c *HeaderCarrier) Set(key, value string) {
	*c = append(*c, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

// Keys returns all keys in the carrier.
func (c *HeaderCarrier) Keys() []string {
	keys := make([]string, len(*c))
	for i, header := range *c {
		keys[i] = header.Key
	}
	return keys
}

func (c *Client) innerProducer(ctx context.Context, key string, data ...any) error {
	messages := make([]kafka.Message, len(data))
	var headers HeaderCarrier
	otel.GetTextMapPropagator().Inject(ctx, &headers)
	for i := range data {
		output, err := sonic.Marshal(&data)
		if err != nil {
			return err
		}
		messages[i] = kafka.Message{
			Value:   output,
			Key:     []byte(key),
			Headers: headers,
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return c.writer.WriteMessages(ctx, messages...)
}
func (c *Client) ProduceByPartitionKey(ctx context.Context, key string, data ...any) error {
	return c.innerProducer(ctx, key, data)
}
func (c *Config) New() *Client {
	client := new(Client)
	client.config = *c
	if len(client.config.Hosts) == 0 {
		client.Init.Error = errors.New("no hosts provided")
		store.PackageStore.Load(client)
		return client
	}
	client.writer = c.newWriter()
	client.Ping(c.DialTimeout, client.Live)
	return client
}

func (c *Config) newWriter() *kafka.Writer {
	transport := new(kafka.Transport)
	if c.DialTimeout == 0 {
		c.DialTimeout = time.Second * 5
	}
	if c.DialTimeout > 0 {
		transport.DialTimeout = c.DialTimeout
	}
	if c.IdleTimeout > 0 {
		transport.IdleTimeout = c.IdleTimeout
	}
	if c.MetadataTTL > 0 {
		transport.MetadataTTL = c.MetadataTTL
	}
	if c.Login != "" {
		mechanism := &plain.Mechanism{
			Username: c.Login, Password: c.Password,
		}
		transport.SASL = mechanism
	}
	if c.PEM != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(c.PEM))

		// Create a TLS configuration.
		tlsConfig := &tls.Config{
			RootCAs: caCertPool, // Use the CA certificate pool.
		}
		transport.TLS = tlsConfig
	}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(c.Hosts...),
		Topic:                  c.Topic,
		AllowAutoTopicCreation: c.AllowAutoTopicCreation,
		Transport:              transport,
		Compression:            kafka.Snappy,
		Balancer:               &kafka.Hash{},
	}
	if c.ReadTimeout != 0 {
		writer.ReadTimeout = c.ReadTimeout
	}
	if c.WriteTimeout != 0 {
		writer.WriteTimeout = c.WriteTimeout
	}
	if c.Async {
		writer.Async = c.Async
	}
	if c.BatchBytes > 0 {
		writer.BatchBytes = c.BatchBytes
	}
	if c.BatchSize > 0 {
		writer.BatchSize = c.BatchSize
	}
	if c.BatchTimeout > 0 {
		writer.BatchTimeout = c.BatchTimeout
	}
	if c.MaxAttempts > 0 {
		writer.MaxAttempts = c.MaxAttempts
	}
	if c.WriteBackoffMax > 0 {
		writer.WriteBackoffMax = c.WriteBackoffMax
	}
	if c.WriteBackoffMin > 0 {
		writer.WriteBackoffMin = c.WriteBackoffMin
	}
	return writer
}
