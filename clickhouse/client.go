package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"net"
	"strings"
)

type Client struct {
	config    Config
	Db        clickhouse.Conn
	initError error
}

func (c *Client) Error() error {
	return c.initError
}

func (c *Client) HasError() bool {
	return c.initError != nil
}

const (
	key = "clickhouse"
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
	opts := &clickhouse.Options{
		Addr: c.Hosts,
		Auth: clickhouse.Auth{
			Database: c.Database,
			Username: c.Login,
			Password: c.Password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	}
	if c.MaxExecTime > 0 {
		opts.Settings = clickhouse.Settings{
			"max_execution_time": c.MaxExecTime,
		}
	}
	if c.DialTimeout > 0 {
		opts.DialTimeout = c.DialTimeout
	}
	if c.MaxOpenConns > 0 {
		opts.MaxOpenConns = c.MaxOpenConns
	}
	if c.MaxIdleConns > 0 {
		opts.MaxIdleConns = c.MaxIdleConns
	}
	if c.ConnMaxLifetime > 0 {
		opts.ConnMaxLifetime = c.ConnMaxLifetime
	}
	if c.BlockBufferSize > 0 {
		opts.BlockBufferSize = c.BlockBufferSize
	}
	if c.MaxCompressionBuffer > 0 {
		opts.MaxCompressionBuffer = c.MaxCompressionBuffer
	}
	conn, err := clickhouse.Open(opts)
	client.initError = err
	client.Db = conn
	return client
}
