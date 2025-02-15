package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"resty.dev/v3"
)

type Client struct {
	config    Config
	Client    *resty.Client
	initError error
}

func (c *Client) Error() error {
	return c.initError
}

func (c *Client) HasError() bool {
	return c.initError != nil
}

const (
	key = "http client resty"
)

func (c *Client) Live(ctx context.Context) error {
	return c.Health(ctx)
}
func (c *Client) AfterShutdown() error {
	return c.Client.Close()
}
func (c *Client) Name() string {
	return fmt.Sprintf(
		"name: %s hosts:%s",
		key,
		c.config.Host,
	)
}

func (c Config) New() *Client {
	client := new(Client)
	defer store.PackageStore.Load(client)
	client.config = c
	restyClient := resty.New().
		SetBaseURL(c.Host).
		SetHeaders(c.Headers)
	client.Client = restyClient
	client.initError = client.Health(context.TODO())
	return client
}

func (c *Client) Health(ctx context.Context) error {
	response, err := c.Client.
		NewRequest().
		SetContext(ctx).
		Get(c.config.HealthHost)
	if err != nil {
		return err
	}
	if response.StatusCode() > 400 {
		return errors.New(fmt.Sprintf("error check component status code: %d", response.StatusCode()))
	}
	return nil
}
