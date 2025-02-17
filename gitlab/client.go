package gitlab

import (
	"context"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"time"
)

const (
	key = "gitlab"
)

type Client struct {
	types.Component
	config *Config
	Db     *gitlab.Client
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
	_, _, err := c.Db.Users.ListUsers(&gitlab.ListUsersOptions{}, gitlab.WithContext(ctx))
	return err
}
func (c *Client) Connection() *gitlab.Client {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c.Db
}

func (c Config) New() *Client {
	client := new(Client)
	if c.HealthTimeout == 0 {
		c.HealthTimeout = 5 * time.Second
	}
	client.config = &c
	client.Db, client.Init.Error = gitlab.NewClient(c.Token, gitlab.WithBaseURL(c.Host))
	if client.Init.Error != nil {
		store.PackageStore.Load(client)
		return client
	}
	client.Ping(c.HealthTimeout, client.Live)
	return client
}
