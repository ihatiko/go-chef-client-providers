package clickhouse

import (
	"time"
)

type Config struct {
	Database             string        `toml:"database"`
	Login                string        `toml:"login"`
	Password             string        `toml:"password"`
	Hosts                []string      `toml:"hosts"`
	DialTimeout          time.Duration `toml:"dial_timeout"`
	MaxOpenConnections   int           `toml:"max_open_connections"`
	MaxIdleConnections   int           `toml:"max_idle_connections" `
	ConnMaxLifetime      time.Duration `toml:"conn_max_lifetime"`
	BlockBufferSize      uint8         `toml:"block_buffer_size"`
	MaxCompressionBuffer int           `toml:"max_compression_buffer"`
	MaxExecTime          int           `toml:"max_exec_time"`
}
