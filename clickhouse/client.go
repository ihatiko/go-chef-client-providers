package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	"net"
	"time"
)

type Client struct {
	types.Component
	config Config
	Db     clickhouse.Conn
}

func (c *Client) GetKey() string {
	return key
}

type Details struct {
	Database string `json:"database"`
	Host     string `json:"host"`
}

func (c *Client) Details() any {
	details := new(Details)
	details.Database = c.config.Database
	details.Host = c.config.Login
	return details
}

func (c *Client) Live(ctx context.Context) error {
	return c.Db.Ping(ctx)
}

const (
	key = "clickhouse"
)

func (c *Client) Shutdown() error {
	return c.Db.Close()
}

func (c Config) New() *Client {
	client := new(Client)
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
	if c.DialTimeout == 0 {
		c.DialTimeout = time.Second * 5
	}
	if c.DialTimeout > 0 {
		opts.DialTimeout = c.DialTimeout
	}
	if c.MaxOpenConnections > 0 {
		opts.MaxOpenConns = c.MaxOpenConnections
	}
	if c.MaxIdleConnections > 0 {
		opts.MaxIdleConns = c.MaxIdleConnections
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
	if err != nil {
		client.Init.Error = err
		store.PackageStore.Load(client)
		return client
	}
	client.Db = conn
	client.Ping(c.DialTimeout, client.Live)
	return client
}
func (c *Client) Connection() clickhouse.Conn {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c.Db
}
