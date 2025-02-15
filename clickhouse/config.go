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
	MaxOpenConns         int           `toml:"max_open_conns" `
	MaxIdleConns         int           `toml:"max_idle_conns" `
	ConnMaxLifetime      time.Duration `toml:"conn_max_lifetime"`
	BlockBufferSize      uint8         `toml:"block_buffer_size"`
	MaxCompressionBuffer int           `toml:"max_compression_buffer"`
	MaxExecTime          int           `toml:"max_exec_time"`
}
