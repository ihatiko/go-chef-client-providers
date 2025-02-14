package s3

type Config struct {
	Host       string `toml:"host"`
	Login      string `toml:"login"`
	Password   string `toml:"password"`
	Token      string `toml:"token"`
	SSL        bool   `toml:"ssl"`
	MaxRetries int    `toml:"max_retries"`
}
