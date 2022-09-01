package etcenter
import (
	"context"
	"crypto/tls"
	redis "github.com/go-redis/redis/v8"
	"net"
	"time"
)

const (
	client = "client"
	clusterClient = "clusterClient"
)
type Client struct {
	*redis.Client //normal client or sentinelClietn
	*redis.ClusterClient // culster
	Type string
}

type Options struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	//// Hook that is called when new connection is established.
	//OnConnect func(ctx context.Context, cn *Conn) error

	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher overhead compared to lifo.
	PoolFIFO bool
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration


	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	//// Limiter interface used to implemented circuit breaker or rate limiter.
	//Limiter Limiter
}

type ClusterOptions struct {
	// A seed list of host:port addresses of cluster nodes.
	Addrs []string

	// NewClient creates a cluster node client with provided name and options.
	NewClient func(opt *Options) *Client

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 3 retries.
	MaxRedirects int

	// Enables read-only commands on slave nodes.
	ReadOnly bool
	// Allows routing read-only commands to the closest master or slave node.
	// It automatically enables ReadOnly.
	RouteByLatency bool
	// Allows routing read-only commands to the random master or slave node.
	// It automatically enables ReadOnly.
	RouteRandomly bool

	// Following options are copied from Options struct.

	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)


	Username string
	Password string

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool

	// PoolSize applies per cluster node and not for the whole cluster.
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration

	TLSConfig *tls.Config
}



//Addr:     "49.232.190.114:63790",
func NewClient(opt Options)*Client {

	return &Client{
		redis.NewClient(&redis.Options{
			Network: opt.Network,
			Addr: opt.Addr,
			Username: opt.Username,
			Password: opt.Password, // no password set
			DB:       opt.DB,  // use default DB
			MaxRetries: opt.MaxRetries,
			DialTimeout: opt.DialTimeout,
			ReadTimeout: opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
			MinIdleConns: opt.MinIdleConns,
			IdleTimeout: opt.IdleTimeout,
			TLSConfig: opt.TLSConfig,
		}),nil, client,
	}
}

func NewClusterClient(opt ClusterOptions)*Client {
	return &Client{
		nil,redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: opt.Addrs,
			Username: opt.Username,
			Password: opt.Password, // no password set
			MaxRetries: opt.MaxRetries,
			DialTimeout: opt.DialTimeout,
			ReadTimeout: opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
			MinIdleConns: opt.MinIdleConns,
			IdleTimeout: opt.IdleTimeout,
			TLSConfig: opt.TLSConfig,
		}), clusterClient,
	}
}