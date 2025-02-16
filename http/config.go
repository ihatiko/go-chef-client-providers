package clickhouse

import "time"

type Config struct {
	Host       string            `toml:"host"`
	Headers    map[string]string `toml:"headers"`
	HealthHost string            `toml:"health_host"`
	Timeout    time.Duration     `toml:"timeout"`
}
