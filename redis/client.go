package redis

import (
	"context"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/store"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
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
	Db        *redis.Client
	cfg       *Config
	initError error
}

func (c Client) Name() string {
	return fmt.Sprintf("name: %s host:%s database:%d", key, c.cfg.Host, c.cfg.Database)
}

func (c Client) Live(ctx context.Context) error {
	return c.Db.Ping(ctx).Err()
}

func (c Client) Error() error {
	return c.initError
}

func (c Client) HasError() bool {
	return c.initError != nil
}

func (c Client) AfterShutdown() error {
	return c.Db.Close()
}
func (c *Config) New() Client {
	client := Client{cfg: c}
	defer store.PackageStore.Load(client)
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
		client.initError = err
		return client
	}
	client.initError = client.Db.Ping(context.Background()).Err()

	return client
}
