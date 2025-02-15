package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	etcd "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
)

type Client struct {
	config    Config
	Db        *etcd.Client
	initError error
}

func (c *Client) Error() error {
	return c.initError
}

func (c *Client) HasError() bool {
	return c.initError != nil
}
func (c *Client) Live(ctx context.Context) error {
	var errorsGroup []error

	wg := &sync.WaitGroup{}
	wg.Add(len(c.config.Hosts))
	for _, addr := range c.config.Hosts {
		go func(addr string) {
			defer wg.Done()
			_, err := c.Db.Status(ctx, addr)
			errorsGroup = append(errorsGroup, err)
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

func (c *Client) AfterShutdown() error {
	return c.Db.Close()
}
func (c *Client) Name() string {
	return fmt.Sprintf(
		"name: %s hosts:%s",
		key,
		strings.Join(c.config.Hosts, ","),
	)
}
func (c Config) New() *Client {
	client := new(Client)
	defer store.PackageStore.Load(client)
	client.config = c
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
	client.Db = cli
	client.initError = err
	if client.initError == nil {
		client.initError = client.Live(context.TODO())
	}
	return client
}
