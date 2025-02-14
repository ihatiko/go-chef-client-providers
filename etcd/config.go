package etcd

import "time"

type Config struct {

	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	AutoSyncInterval time.Duration `toml:"auto_sync_interval"`

	// DialTimeout is the timeout for failing to establish a connection.
	DialTimeout time.Duration `toml:"dial_timeout"`

	// DialKeepAliveTime is the time after which client pings the server to see if
	// transport is alive.
	DialKeepAliveTime time.Duration `toml:"dial_keep_alive_time"`

	// DialKeepAliveTimeout is the time that the client waits for a response for the
	// keep-alive probe. If the response is not received in this time, the connection is closed.
	DialKeepAliveTimeout time.Duration `toml:"dial_keep_alive_timeout"`

	// MaxCallSendMsgSize is the client-side request send limit in bytes.
	// If 0, it defaults to 2.0 MiB (2 * 1024 * 1024).
	// Make sure that "MaxCallSendMsgSize" < server-side default send/recv limit.
	// ("--max-request-bytes" flag to etcd or "embed.Config.MaxRequestBytes").
	MaxCallSendMsgSize int `json:"max_call_send_msg_size"`

	// MaxCallRecvMsgSize is the client-side response receive limit.
	// If 0, it defaults to "math.MaxInt32", because range response can
	// easily exceed request send limits.
	// Make sure that "MaxCallRecvMsgSize" >= server-side default send/recv limit.
	// ("--max-request-bytes" flag to etcd or "embed.Config.MaxRequestBytes").
	MaxCallRecvMsgSize int `toml:"max_call_recv_msg_size"`
	// Endpoints is a list of URLs.
	Hosts []string `toml:"hosts"`
	// Login is a user name for authentication.
	Login string `toml:"login"`

	// Password is a password for authentication.
	Password string `toml:"password"`
	// TLS cert
	PEM string `toml:"pem"`
	// PermitWithoutStream when set will allow client to send keepalive pings to server without any active streams(RPCs).
	PermitWithoutStream bool `toml:"permit_without_stream"`

	// MaxUnaryRetries is the maximum number of retries for unary RPCs.
	MaxUnaryRetries uint `toml:"max_unary_retries"`

	// BackoffWaitBetween is the wait time before retrying an RPC.
	BackoffWaitBetween time.Duration `toml:"backoff_wait_between"`

	// BackoffJitterFraction is the jitter fraction to randomize backoff wait time.
	BackoffJitterFraction float64 `toml:"backoff_jitter_fraction"`
}
