package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ihatiko/go-chef-core-sdk/types"
	"log"
	"time"

	"github.com/ihatiko/go-chef-core-sdk/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var defaultTimeout = 5 * time.Second

const (
	maxOpenConnections   = 60
	connMaxLifetime      = 120
	maxIdleConnections   = 30
	connMaxIdleTime      = 20
	defaultQueryExecMode = "simple_protocol"
	defaultSslMode       = "disable"
)
const (
	key           = "postgres"
	defaultDriver = "pgx"
)

var metricProvider *metric.MeterProvider

func getMetricProvider() *metric.MeterProvider {
	if metricProvider == nil {
		exporter, err := prometheus.New()
		if err != nil {
			log.Fatal(err)
		}
		metricProvider = metric.NewMeterProvider(metric.WithReader(exporter))
	}

	return metricProvider
}

func (c *Config) toPgConnection() string {
	queryExecMode := c.QueryExecMode
	if queryExecMode == "" {
		queryExecMode = defaultQueryExecMode
	}

	dataSourceName := fmt.Sprintf("host=%s port=%d login=%s dbname=%s password=%s sslmode=%s default_query_exec_mode=%s",
		c.Host,
		c.Port,
		c.Login,
		c.Database,
		c.Password,
		c.SSLMode,
		queryExecMode,
	)
	return dataSourceName
}

type Client struct {
	types.Component
	Db  *sqlx.DB
	cfg *Config
}

func (c *Client) GetKey() string {
	return key
}

type Details struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	PgDriver string `json:"pg_driver"`
}

func (c *Client) Details() any {
	details := Details{}
	details.Host = c.cfg.Host
	details.Port = c.cfg.Port
	details.Database = c.cfg.Database
	details.PgDriver = c.cfg.PgDriver
	return details
}

func (c *Client) Live(ctx context.Context) error {
	if c.Db == nil {
		return c.Init.Error
	}
	return c.Db.PingContext(ctx)
}

func (c *Client) Shutdown() error {
	return c.Db.Close()
}

func (c *Client) Connection() *sqlx.DB {
	defer store.PackageStore.Load(c)
	if c.Db == nil {
		return new(sqlx.DB)
	}
	c.AwaitPing()
	return c.Db
}

func (c *Config) New() *Client {
	client := new(Client)
	client.cfg = c
	if c.PgDriver == "" {
		c.PgDriver = defaultDriver
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = maxOpenConnections
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = maxIdleConnections
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = connMaxIdleTime * time.Second
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = connMaxLifetime * time.Second
	}
	if c.SSLMode == "" {
		c.SSLMode = defaultSslMode
	}
	connectionString := c.toPgConnection()

	db, err := otelsqlx.Connect(c.PgDriver, connectionString,
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
			attribute.KeyValue{Key: "driver", Value: attribute.StringSliceValue(sql.Drivers())},
		),
		otelsql.WithDBName(c.Database),
		otelsql.WithMeterProvider(getMetricProvider()),
	)
	client.Db = db
	if err != nil {
		client.Init.Error = err
		store.PackageStore.Load(client)
		return client
	}
	db.SetMaxOpenConns(c.MaxIdleConnections)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetMaxIdleConns(c.MaxIdleConnections)
	db.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	client.Ping(defaultTimeout, client.Live)
	return client
}
