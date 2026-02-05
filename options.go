package kraken

import (
	"io"
	"net/http"
	"time"
)

// RateLimitTier represents the Kraken account tier for rate limiting.
type RateLimitTier int

const (
	// TierStarter is the default tier with lowest limits.
	TierStarter RateLimitTier = iota
	// TierIntermediate is the intermediate tier with higher limits.
	TierIntermediate
	// TierPro is the professional tier with highest limits.
	TierPro
)

// String returns the string representation of the tier.
func (t RateLimitTier) String() string {
	switch t {
	case TierStarter:
		return "Starter"
	case TierIntermediate:
		return "Intermediate"
	case TierPro:
		return "Pro"
	default:
		return "Unknown"
	}
}

// Logger defines the interface for logging.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Options configures the Kraken client.
type Options struct {
	// APIKey is the public API key.
	APIKey string

	// APISecret is the private API secret (base64 encoded).
	APISecret string

	// BaseURL is the base URL for the REST API.
	// Defaults to DefaultBaseURL.
	BaseURL string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, a default client with sensible timeouts is used.
	HTTPClient *http.Client

	// Timeout is the default timeout for requests.
	// Defaults to 30 seconds.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts for failed requests.
	// Defaults to 3.
	MaxRetries int

	// RetryWaitMin is the minimum wait time between retries.
	// Defaults to 500ms.
	RetryWaitMin time.Duration

	// RetryWaitMax is the maximum wait time between retries.
	// Defaults to 30s.
	RetryWaitMax time.Duration

	// PublicWSURL is the WebSocket URL for public channels.
	// Defaults to PublicWSURL constant.
	PublicWSURL string

	// PrivateWSURL is the WebSocket URL for private channels.
	// Defaults to PrivateWSURL constant.
	PrivateWSURL string

	// RateLimitTier is the account tier for rate limiting.
	// Defaults to TierStarter.
	RateLimitTier RateLimitTier

	// DisableRateLimiting disables the built-in rate limiter.
	DisableRateLimiting bool

	// UserAgent is the User-Agent header value.
	// Defaults to "go-kraken-sdk/1.0".
	UserAgent string

	// OTPGenerator is an optional function that returns a one-time password
	// for 2FA authentication. If set, the OTP will be included in private requests.
	OTPGenerator func() string

	// Logger is an optional logger for debugging.
	Logger Logger

	// DebugWriter is an optional writer for debug output.
	// If set, request/response bodies will be written here.
	DebugWriter io.Writer

	// OnError is an optional callback for handling errors.
	OnError func(error)

	// WSReconnectInterval is the interval between WebSocket reconnection attempts.
	// Defaults to 5 seconds.
	WSReconnectInterval time.Duration

	// WSPingInterval is the interval between WebSocket ping messages.
	// Defaults to 30 seconds.
	WSPingInterval time.Duration

	// WSReadTimeout is the read timeout for WebSocket connections.
	// Defaults to 60 seconds.
	WSReadTimeout time.Duration

	// WSWriteTimeout is the write timeout for WebSocket connections.
	// Defaults to 10 seconds.
	WSWriteTimeout time.Duration
}

// Option is a functional option for configuring the client.
type Option func(*Options)

// WithAPIKey sets the API key and secret for authentication.
func WithAPIKey(key, secret string) Option {
	return func(o *Options) {
		o.APIKey = key
		o.APISecret = secret
	}
}

// WithBaseURL sets a custom base URL for the REST API.
func WithBaseURL(url string) Option {
	return func(o *Options) {
		o.BaseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(retries int) Option {
	return func(o *Options) {
		o.MaxRetries = retries
	}
}

// WithRetryWait sets the retry wait times.
func WithRetryWait(min, max time.Duration) Option {
	return func(o *Options) {
		o.RetryWaitMin = min
		o.RetryWaitMax = max
	}
}

// WithRateLimitTier sets the rate limit tier.
func WithRateLimitTier(tier RateLimitTier) Option {
	return func(o *Options) {
		o.RateLimitTier = tier
	}
}

// WithoutRateLimiting disables the built-in rate limiter.
func WithoutRateLimiting() Option {
	return func(o *Options) {
		o.DisableRateLimiting = true
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(userAgent string) Option {
	return func(o *Options) {
		o.UserAgent = userAgent
	}
}

// With2FA sets the OTP generator function for 2FA authentication.
func With2FA(generator func() string) Option {
	return func(o *Options) {
		o.OTPGenerator = generator
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithDebugWriter sets a debug writer for request/response logging.
func WithDebugWriter(w io.Writer) Option {
	return func(o *Options) {
		o.DebugWriter = w
	}
}

// WithOnError sets the error callback.
func WithOnError(fn func(error)) Option {
	return func(o *Options) {
		o.OnError = fn
	}
}

// WithWSReconnectInterval sets the WebSocket reconnection interval.
func WithWSReconnectInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.WSReconnectInterval = interval
	}
}

// WithWSPingInterval sets the WebSocket ping interval.
func WithWSPingInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.WSPingInterval = interval
	}
}

// WithWSTimeouts sets the WebSocket read and write timeouts.
func WithWSTimeouts(read, write time.Duration) Option {
	return func(o *Options) {
		o.WSReadTimeout = read
		o.WSWriteTimeout = write
	}
}

// DefaultOptions returns the default client options.
func DefaultOptions() *Options {
	return &Options{
		BaseURL:             DefaultBaseURL,
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryWaitMin:        500 * time.Millisecond,
		RetryWaitMax:        30 * time.Second,
		PublicWSURL:         PublicWSURL,
		PrivateWSURL:        PrivateWSURL,
		RateLimitTier:       TierStarter,
		UserAgent:           "go-kraken-sdk/1.0",
		WSReconnectInterval: 5 * time.Second,
		WSPingInterval:      30 * time.Second,
		WSReadTimeout:       60 * time.Second,
		WSWriteTimeout:      10 * time.Second,
	}
}

// applyOptions applies the given options to the default options.
func applyOptions(opts []Option) *Options {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}
