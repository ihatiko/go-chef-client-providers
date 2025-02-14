package redis

import "time"

type Config struct {
	Host               string        `json:"host"`
	Login              string        `toml:"login"`
	Password           string        `toml:"password"`
	Database           int           `toml:"database"`
	SentinelHosts      []string      `toml:"sentinel_hosts"`
	MasterName         string        `toml:"master_name"`
	DialTimeout        time.Duration `toml:"dial_timeout"`
	ReadTimeout        time.Duration `toml:"read_timeout"`
	WriteTimeout       time.Duration `toml:"write_timeout"`
	ConnMaxIdleTime    time.Duration `toml:"conn_max_idle_time"`
	ConnMaxLifetime    time.Duration `toml:"conn_max_lifetime"`
	MaxIdleConnections int           `toml:"max_idle_connections"`
}
