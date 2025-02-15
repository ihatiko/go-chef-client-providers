package s3

import (
	"context"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

const (
	key = "s3"
)

type Client struct {
	initError error
	config    *Config
	Db        *minio.Client
}

func (c *Client) Live(ctx context.Context) error {
	_, err := c.Db.HealthCheck(c.config.HealthTimeout)
	return err
}
func (c *Client) Error() error {
	return c.initError
}

func (c *Client) Name() string {
	return fmt.Sprintf(
		"name: %s host:%s",
		key,
		c.config.Host,
	)
}
func (c *Client) HasError() bool {
	return c.initError != nil
}

func (c Config) New() *Client {
	client := new(Client)
	defer store.PackageStore.Load(client)
	opts := new(minio.Options)
	if c.Login != "" {
		opts.Creds = credentials.NewStaticV4(c.Login, c.Password, c.Token)
	}
	opts.Secure = c.SSL
	if opts.MaxRetries > 0 {
		opts.MaxRetries = c.MaxRetries
	}
	if c.HealthTimeout == 0 {
		c.HealthTimeout = time.Second * 5
	}
	client.config = &c
	client.Db, client.initError = minio.New(c.Host, opts)
	if client.initError == nil {
		client.initError = client.Live(context.TODO())
	}
	return client
}
