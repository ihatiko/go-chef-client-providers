package kafka

import "time"

type Config struct {
	AllowAutoTopicCreation bool     `toml:"allow_auto_topic_creation"`
	Hosts                  []string `toml:"hosts"` // Brokers
	// Topic is the name of the topic that the writer will produce messages to.
	//
	// Setting this field or not is a mutually exclusive option. If you set Topic
	// here, you must not set Topic for any produced Message. Otherwise, if you	do
	// not set Topic, every Message must have Topic specified.
	Topic string `toml:"topic"` // Topics to write data

	// Limit on how many attempts will be made to deliver a message.
	//
	// The default is to try at most 10 times.
	MaxAttempts int `toml:"max_attempts"`

	// WriteBackoffMin optionally sets the smallest amount of time the writer waits before
	// it attempts to write a batch of messages
	//
	// Default: 100ms
	WriteBackoffMin time.Duration `toml:"write_backoff_min"`

	// WriteBackoffMax optionally sets the maximum amount of time the writer waits before
	// it attempts to write a batch of messages
	//
	// Default: 1s
	WriteBackoffMax time.Duration `toml:"write_backoff_min"`

	// Limit on how many messages will be buffered before being sent to a
	// partition.
	//
	// The default is to use a target batch size of 100 messages.
	BatchSize int `toml:"batch_size"`

	// Limit the maximum size of a request in bytes before being sent to
	// a partition.
	//
	// The default is to use a kafka default value of 1048576.
	BatchBytes int64 `toml:"batch_bytes"`

	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	BatchTimeout time.Duration `toml:"batch_timeout"`

	// Timeout for read operations performed by the Writer.
	//
	// Defaults to 10 seconds.
	ReadTimeout time.Duration `toml:"read_timeout"`

	// Timeout for write operation performed by the Writer.
	//
	// Defaults to 10 seconds.
	WriteTimeout time.Duration `toml:"write_timeout"`

	// Number of acknowledges from partition replicas required before receiving
	// a response to a produce request, the following values are supported:
	//
	//  RequireNone (0)  fire-and-forget, do not wait for acknowledgements from the
	//  RequireOne  (1)  wait for the leader to acknowledge the writes
	//  RequireAll  (-1) wait for the full ISR to acknowledge the writes
	//
	// Defaults to RequireNone.
	RequiredAcks int `toml:"required_acks"`

	// Setting this flag to true causes the WriteMessages method to never block.
	// It also means that errors are ignored since the caller will not receive
	// the returned value. Use this only if you don't care about guarantees of
	// whether the messages were written to kafka.
	//
	// Defaults to false.
	Async bool `toml:"async"`

	// Time limit set for establishing connections to the kafka cluster. This
	// limit includes all round trips done to establish the connections (TLS
	// handshake, SASL negotiation, etc...).
	//
	// Defaults to 5s.
	DialTimeout time.Duration `toml:"dial_timeout"`

	// Maximum amount of time that connections will remain open and unused.
	// The transport will manage to automatically close connections that have
	// been idle for too long, and re-open them on demand when the transport is
	// used again.
	//
	// Defaults to 30s.
	IdleTimeout time.Duration `toml:"idle_timeout"`

	// TTL for the metadata cached by this transport. Note that the value
	// configured here is an upper bound, the transport randomizes the TTLs to
	// avoid getting into states where multiple clients end up synchronized and
	// cause bursts of requests to the kafka broker.
	//
	// Default to 6s.
	MetadataTTL time.Duration `toml:"metadata_ttl"`
	Login       string        `toml:"login"`
	Password    string        `toml:"password"`
	PEM         string        `toml:"pem"`
}
