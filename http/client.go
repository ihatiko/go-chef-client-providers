package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	"resty.dev/v3"
	"time"
)

const defaultTimeout = 3 * time.Second

const (
	key = "resty-http-client"
)

type Client struct {
	types.Component
	config Config
	client *resty.Client
}

func (c *Client) Live(ctx context.Context) error {
	return c.liveness(ctx)
}
func (c *Client) Shutdown() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *Client) Name() string {
	return c.Id.String()
}

type Details struct {
	Host       string        `json:"host"`
	HealthHost string        `json:"health_host"`
	Timeout    time.Duration `json:"timeout"`
}

func (c *Client) GetKey() string {
	return key
}

func (c *Client) Details() any {
	details := new(Details)
	details.Host = c.config.Host
	details.HealthHost = c.config.HealthHost
	details.Timeout = c.config.Timeout
	return details
}
func (c *Client) Connection() *resty.Client {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c.client
}

func (c Config) New() *Client {
	client := new(Client)
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}
	client.config = c
	restyClient := resty.New().
		SetBaseURL(c.Host).
		SetHeaders(c.Headers)

	client.client = restyClient
	client.Ping(c.Timeout, client.liveness)
	return client
}

func (c *Client) liveness(ctx context.Context) error {
	response, err := c.client.
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
