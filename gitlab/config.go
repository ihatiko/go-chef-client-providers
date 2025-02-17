package gitlab

import "time"

type Config struct {
	Host          string        `toml:"host"`
	Token         string        `toml:"token"`
	HealthTimeout time.Duration `toml:"health_timeout"`
}
