package s3

import (
	"context"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

const (
	key = "s3"
)

type Client struct {
	types.Component
	config *Config
	Db     *minio.Client
}

func (c *Client) GetKey() string {
	return key
}

type Details struct {
	Host string `json:"host"`
}

func (c *Client) Details() any {
	details := Details{Host: c.config.Host}
	return details
}

func (c *Client) Shutdown() error {
	return nil
}

func (c *Client) Live(ctx context.Context) error {
	_, err := c.Db.ListBuckets(ctx)
	return err
}

func (c *Client) Connection() *minio.Client {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c.Db
}

func (c Config) New() *Client {
	client := new(Client)
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
	client.Db, client.Init.Error = minio.New(c.Host, opts)
	if client.Init.Error != nil {
		store.PackageStore.Load(client)
		return client
	}
	client.Ping(c.HealthTimeout, client.Live)
	return client
}
