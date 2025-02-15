package clickhouse

type Config struct {
	Host       string            `toml:"host"`
	Headers    map[string]string `toml:"headers"`
	HealthHost string            `toml:"health_host"`
}
