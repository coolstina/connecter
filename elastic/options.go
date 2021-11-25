package elastic

import (
	"net/http"
	"time"

	"github.com/olivere/elastic"
)

// Decoder is used to decode responses from Elasticsearch.
// Users of elastic can implement their own marshaler for advanced purposes
// and set them per Client (see WithDecoder). If none is specified,
// elastic.DefaultDecoder is used.
type Decoder = elastic.Decoder

// Logger specifies the interface for all log operations.
type Logger = elastic.Logger

// Retrier decides whether to retry a failed HTTP request with Elasticsearch.
type Retrier = elastic.Retrier

type Option interface {
	apply(*options)
}

type optionFunc func(ops *options)

func (o optionFunc) apply(ops *options) {
	o(ops)
}

type options struct {
	httpClient                *http.Client
	basicAuthUsername         *string
	basicAuthPassword         *string
	urls                      []string
	scheme                    *string
	snifferEnabled            *bool
	snifferTimeoutStartup     *time.Duration
	snifferTimeout            *time.Duration
	snifferInterval           *time.Duration
	snifferCallback           SnifferCallback
	healthCheckEnabled        *bool
	healthCheckTimeoutStartup *time.Duration
	healthCheckTimeout        *time.Duration
	healthCheckInterval       *time.Duration
	gzipEnabled               *bool
	decoder                   Decoder
	requiredPlugins           []string
	errorlog                  Logger
	infolog                   Logger
	tracelog                  Logger
	sendGetBodyAs             *string
	retrier                   Retrier
	headers                   http.Header
}

var (
	// noRetries is a retrier that does not retry.
	noRetries = elastic.NewStopRetrier()
)

// WithHttpClient can be used to specify the http.Client to use when making
// HTTP requests to Elasticsearch.
func WithHttpClient(httpClient *http.Client) Option {
	return optionFunc(func(ops *options) {
		elastic.SetURL(DefaultURL)
		ops.httpClient = httpClient
	})
}

// WithBasicAuth can be used to specify the HTTP Basic Auth credentials to
// use when making HTTP requests to Elasticsearch.
func WithBasicAuth(username, password string) Option {
	return optionFunc(func(ops *options) {
		ops.basicAuthUsername = &username
		ops.basicAuthPassword = &password
	})
}

// WithSetURL defines the URL endpoints of the Elasticsearch nodes. Notice that
// when sniffing is enabled, these URLs are used to initially sniff the
// cluster on startup.
func WithSetURL(urls ...string) Option {
	return optionFunc(func(ops *options) {
		ops.urls = urls
	})
}

// WithScheme sets the HTTP scheme to look for when sniffing (http or https).
// This is http by default.
func WithScheme(scheme string) Option {
	return optionFunc(func(ops *options) {
		ops.scheme = &scheme
	})
}

// WithSniff enables or disables the sniffer (enabled by default).
func WithSniff(enabled bool) Option {
	return optionFunc(func(ops *options) {
		ops.snifferEnabled = &enabled
	})
}

// WithSnifferTimeoutStartup sets the timeout for the sniffer that is used
// when creating a new client. The default is 5 seconds. Notice that the
// timeout being used for subsequent sniffing processes is set with
// SetSnifferTimeout.
func WithSnifferTimeoutStartup(timeout time.Duration) Option {
	return optionFunc(func(ops *options) {
		ops.snifferTimeoutStartup = &timeout
	})
}

// WithSnifferTimeout sets the timeout for the sniffer that finds the
// nodes in a cluster. The default is 2 seconds. Notice that the timeout
// used when creating a new client on startup is usually greater and can
// be set with SetSnifferTimeoutStartup.
func WithSnifferTimeout(timeout time.Duration) Option {
	return optionFunc(func(ops *options) {
		ops.snifferTimeout = &timeout
	})
}

// WithSnifferInterval sets the interval between two sniffing processes.
// The default interval is 15 minutes.
func WithSnifferInterval(interval time.Duration) Option {
	return optionFunc(func(ops *options) {
		ops.snifferInterval = &interval
	})
}

// SnifferCallback defines the protocol for sniffing decisions.
type SnifferCallback func(*elastic.NodesInfoNode) bool

// WithSnifferCallback allows the caller to modify sniffer decisions.
// When setting the callback, the given SnifferCallback is called for
// each (healthy) node found during the sniffing process.
// If the callback returns false, the node is ignored: No requests
// are routed to it.
func WithSnifferCallback(f SnifferCallback) Option {
	return optionFunc(func(ops *options) {
		ops.snifferCallback = f
	})
}

// WithHealthCheck enables or disables healthChecks (enabled by default).
func WithHealthCheck(enabled bool) Option {
	return optionFunc(func(ops *options) {
		ops.healthCheckEnabled = &enabled
	})
}

// WithHealthCheckTimeoutStartup sets the timeout for the initial health check.
// The default timeout is 5 seconds (see DefaultHealthCheckTimeoutStartup).
// Notice that timeouts for subsequent health checks can be modified with
// SetHealthCheckTimeout.
func WithHealthCheckTimeoutStartup(timeout time.Duration) Option {
	return optionFunc(func(ops *options) {
		ops.healthCheckTimeoutStartup = &timeout
	})
}

// WithHealthCheckTimeout sets the timeout for periodic health checks.
// The default timeout is 1 second (see DefaultHealthCheckTimeout).
// Notice that a different (usually larger) timeout is used for the initial
// healthCheck, which is initiated while creating a new client.
// The startup timeout can be modified with SetHealthCheckTimeoutStartup.
func WithHealthCheckTimeout(timeout time.Duration) Option {
	return optionFunc(func(ops *options) {
		ops.healthCheckTimeout = &timeout
	})
}

// WithHealthCheckInterval sets the interval between two health checks.
// The default interval is 60 seconds.
func WithHealthCheckInterval(timeout time.Duration) Option {
	return optionFunc(func(ops *options) {
		ops.healthCheckInterval = &timeout
	})
}

// WithGzip enables or disables gzip compression (disabled by default).
func WithGzip(enabled bool) Option {
	return optionFunc(func(ops *options) {
		ops.gzipEnabled = &enabled
	})
}

// WithDecoder sets the Decoder to use when decoding data from Elasticsearch.
// DefaultDecoder is used by default.
func WithDecoder(decoder Decoder) Option {
	return optionFunc(func(ops *options) {
		ops.decoder = decoder
	})
}

// WithRequiredPlugins can be used to indicate that some plugins are required
// before a Client will be created.
func WithRequiredPlugins(plugins ...string) Option {
	return optionFunc(func(ops *options) {
		ops.requiredPlugins = plugins
	})
}

// WithErrorLog sets the logger for critical messages like nodes joining
// or leaving the cluster or failing requests. It is nil by default.
func WithErrorLog(logger Logger) Option {
	return optionFunc(func(ops *options) {
		ops.errorlog = logger
	})
}

// WithInfoLog sets the logger for informational messages, e.g. requests
// and their response times. It is nil by default.
func WithInfoLog(logger Logger) Option {
	return optionFunc(func(ops *options) {
		ops.infolog = logger
	})
}

// WithTraceLog specifies the log.Logger to use for output of HTTP requests
// and responses which is helpful during debugging. It is nil by default.
func WithTraceLog(logger Logger) Option {
	return optionFunc(func(ops *options) {
		ops.tracelog = logger
	})
}

// WithSendGetBodyAs specifies the HTTP method to use when sending a GET request
// with a body. It is GET by default.
func WithSendGetBodyAs(httpMethod string) Option {
	return optionFunc(func(ops *options) {
		ops.sendGetBodyAs = &httpMethod
	})
}

// WithRetrier specifies the HTTP method to use when sending a GET request
// with a body. It is GET by default.
func WithRetrier(retrier Retrier) Option {
	return optionFunc(func(ops *options) {
		ops.retrier = retrier
	})
}

// WithHeaders adds a list of default HTTP headers that will be added to
// each requests executed by PerformRequest.
func WithHeaders(headers http.Header) Option {
	return optionFunc(func(ops *options) {
		ops.headers = headers
	})
}
