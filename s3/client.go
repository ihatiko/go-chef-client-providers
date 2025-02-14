package s3

import (
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	key = "s3"
)

type Client struct {
	initError error
	config    *Config
	Cli       *minio.Client
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
	client.Cli, client.initError = minio.New(c.Host, opts)

	return client
}
