package postgresql

import "time"

type Config struct {
	Port               int           `toml:"port"`
	Host               string        `toml:"host"`
	Login              string        `toml:"login"`
	Password           string        `toml:"password"`
	Database           string        `toml:"database"`
	SSLMode            string        `toml:"ssl_mode"`
	PgDriver           string        `toml:"pg_driver"`
	AutoMigrate        bool          `toml:"auto_migrate"`
	QueryExecMode      string        `toml:"query_exec_mode"`
	MaxOpenConnections int           `toml:"max_open_connections"`
	MaxIdleConnections int           `toml:"max_idle_connections"`
	ConnMaxLifetime    time.Duration `toml:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `toml:"conn_max_idle_time"`
}
