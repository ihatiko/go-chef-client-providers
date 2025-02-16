package redis

import (
	"context"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"github.com/ihatiko/go-chef-core-sdk/types"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	defaultHealthTimeout      = 5
	defaultReadTimeout        = 5
	defaultWriteTimeout       = 5
	defaultConnMaxLifetime    = 120
	defaultMaxIdleConnections = 30
	defaultConnMaxIdleTime    = 20
)
const (
	key = "redis"
)

type Client struct {
	types.Component
	Db  *redis.Client
	cfg *Config
}

func (c *Client) GetKey() string {
	return key
}

type Details struct {
	Host          string   `json:"host,omitempty"`
	Database      int      `toml:"database,omitempty"`
	SentinelHosts []string `toml:"sentinel_hosts,omitempty"`
	MasterName    string   `toml:"master_name,omitempty"`
}

func (c *Client) Details() any {
	details := Details{}
	details.MasterName = c.cfg.MasterName
	details.Host = c.cfg.Host
	details.Database = c.cfg.Database
	details.SentinelHosts = c.cfg.SentinelHosts
	return details
}

func (c *Client) Shutdown() error {
	return c.Db.Close()
}

func (c *Client) Live(ctx context.Context) error {
	return c.Db.Ping(ctx).Err()
}
func (c *Client) Connection() *redis.Client {
	defer store.PackageStore.Load(c)
	c.AwaitPing()
	return c.Db
}

func (c *Config) New() *Client {
	client := new(Client)
	client.cfg = c
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = defaultConnMaxLifetime * time.Second
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = defaultConnMaxIdleTime * time.Second
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout * time.Second
	}
	if c.HealthTimeout == 0 {
		c.HealthTimeout = defaultHealthTimeout * time.Second
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = defaultMaxIdleConnections
	}
	if len(c.SentinelHosts) > 0 {
		client.Db = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:      c.MasterName,
			SentinelAddrs:   c.SentinelHosts,
			DB:              c.Database,
			WriteTimeout:    c.WriteTimeout,
			ReadTimeout:     c.ReadTimeout,
			ConnMaxIdleTime: c.ConnMaxIdleTime,
			ConnMaxLifetime: c.ConnMaxLifetime,
			MaxIdleConns:    c.MaxIdleConnections,
		})
	} else {
		client.Db = redis.NewClient(&redis.Options{
			Addr:            c.Host,
			Password:        c.Password,
			DB:              c.Database,
			Username:        c.Login,
			WriteTimeout:    c.WriteTimeout,
			ReadTimeout:     c.ReadTimeout,
			ConnMaxIdleTime: c.ConnMaxIdleTime,
			ConnMaxLifetime: c.ConnMaxLifetime,
			MaxIdleConns:    c.MaxIdleConnections,
		})
	}
	if err := redisotel.InstrumentTracing(client.Db); err != nil {
		client.Init.Error = err
		store.PackageStore.Load(client)
		return client
	}
	client.Ping(c.ReadTimeout, client.Live)
	return client
}
