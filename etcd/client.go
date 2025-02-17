package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	etcd "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type Client struct {
	types.Component
	config Config
	Db     *etcd.Client
}

func (c *Client) GetKey() string {
	return key
}

type Details struct {
	Hosts []string `json:"hosts"`
}

func (c *Client) Details() any {
	details := new(Details)
	details.Hosts = c.config.Hosts
	return details
}

func (c *Client) Shutdown() error {
	return c.Db.Close()
}

func (c *Client) Live(ctx context.Context) error {
	if c.Db == nil {
		return c.Init.Error
	}
	var errorsGroup []error

	wg := &sync.WaitGroup{}
	wg.Add(len(c.config.Hosts))
	for _, addr := range c.config.Hosts {
		go func(addr string) {
			defer wg.Done()
			_, err := c.Db.Status(ctx, addr)
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
	return fmt.Errorf("etcd errors: %s", errorsGroup)
}

const (
	key = "etcd"
)
const (
	defaultTimeout = 5 * time.Second
)

func (c Config) New() *Client {
	client := new(Client)
	client.config = c
	if c.DialTimeout == 0 {
		c.DialTimeout = defaultTimeout
	}
	config := etcd.Config{
		Endpoints:             c.Hosts,
		Username:              c.Login,
		Password:              c.Password,
		DialTimeout:           c.DialTimeout,
		AutoSyncInterval:      c.AutoSyncInterval,
		DialKeepAliveTime:     c.DialKeepAliveTime,
		DialKeepAliveTimeout:  c.DialKeepAliveTimeout,
		MaxCallSendMsgSize:    c.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:    c.MaxCallRecvMsgSize,
		PermitWithoutStream:   c.PermitWithoutStream,
		MaxUnaryRetries:       c.MaxUnaryRetries,
		BackoffJitterFraction: c.BackoffJitterFraction,
		BackoffWaitBetween:    c.BackoffWaitBetween,
	}

	if client.config.PEM != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(client.config.PEM))

		// Create a TLS configuration.
		tlsConfig := &tls.Config{
			RootCAs: caCertPool, // Use the CA certificate pool.
		}
		config.TLS = tlsConfig
	}
	cli, err := etcd.New(config)
	if err != nil {
		client.Init.Error = err
		store.PackageStore.Load(client)
		return client
	}
	client.Db = cli
	client.Ping(c.DialTimeout, client.Live)
	return client
}
func (c *Client) Connection() *etcd.Client {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c.Db
}
