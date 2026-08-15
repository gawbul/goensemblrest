package goensemblrest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// Version is the current semantic version of the goensemblrest library.
	Version = "0.1.0"

	// DefaultBaseURL is the official Ensembl REST API endpoint.
	DefaultBaseURL = "https://rest.ensembl.org"

	// DefaultContentType is the default MIME type used in requests.
	DefaultContentType = "application/json"

	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 60 * time.Second

	// DefaultMaxAttempts is the default maximum number of retry attempts.
	DefaultMaxAttempts = 5

	// DefaultReqsPerSec is the default rate limit (15 requests per second).
	DefaultReqsPerSec = 15

	// DefaultWallTime is the sliding window duration for rate limiting.
	DefaultWallTime = 1 * time.Second
)

// DefaultUserAgent returns the default User-Agent string.
var DefaultUserAgent = fmt.Sprintf("goensemblrest/%s (Go 1.26; +https://github.com/gawbul/goensemblrest)", Version)

// ClientOption is a functional option for configuring a Client.
type ClientOption func(*Client) error

// Client is the Ensembl REST API client. It is safe for concurrent use across multiple goroutines.
type Client struct {
	mu           sync.RWMutex
	baseURL      *url.URL
	httpClient   *http.Client
	userAgent    string
	extraHeaders http.Header
	limiter      *rateLimiter
	maxAttempts  int
	timeout      time.Duration
	wallTime     time.Duration
	reqsPerSec   int
}

// NewClient creates and initializes a new Ensembl REST API Client with optional configuration.
func NewClient(opts ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(DefaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid default base URL: %w", err)
	}

	c := &Client{
		baseURL:      parsedURL,
		httpClient:   &http.Client{Timeout: DefaultTimeout},
		userAgent:    DefaultUserAgent,
		extraHeaders: make(http.Header),
		maxAttempts:  DefaultMaxAttempts,
		timeout:      DefaultTimeout,
		wallTime:     DefaultWallTime,
		reqsPerSec:   DefaultReqsPerSec,
	}

	for _, opt := range opts {
		if opt != nil {
			if err := opt(c); err != nil {
				return nil, fmt.Errorf("failed to apply client option: %w", err)
			}
		}
	}

	c.limiter = newRateLimiter(c.reqsPerSec, c.wallTime)

	return c, nil
}

// WithBaseURL configures a custom base URL for the client.
func WithBaseURL(rawURL string) ClientOption {
	return func(c *Client) error {
		rawURL = strings.TrimRight(rawURL, "/")
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return fmt.Errorf("invalid base URL %q: %w", rawURL, err)
		}
		c.baseURL = parsed
		return nil
	}
}

// WithHTTPClient configures a custom *http.Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("httpClient cannot be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout sets the request timeout for the client's HTTP requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be greater than zero")
		}
		c.timeout = timeout
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithRateLimit configures client-side rate limiting (requests per second and window duration).
func WithRateLimit(reqsPerSec int, window time.Duration) ClientOption {
	return func(c *Client) error {
		if reqsPerSec <= 0 {
			return fmt.Errorf("reqsPerSec must be positive")
		}
		if window <= 0 {
			return fmt.Errorf("window must be positive")
		}
		c.reqsPerSec = reqsPerSec
		c.wallTime = window
		return nil
	}
}

// WithMaxAttempts sets the maximum number of retry attempts for transient server errors.
func WithMaxAttempts(attempts int) ClientOption {
	return func(c *Client) error {
		if attempts < 1 {
			return fmt.Errorf("max attempts must be at least 1")
		}
		c.maxAttempts = attempts
		return nil
	}
}

// WithUserAgent sets a custom User-Agent header value.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithHeader adds a persistent default header to all client requests.
func WithHeader(key, value string) ClientOption {
	return func(c *Client) error {
		c.extraHeaders.Set(key, value)
		return nil
	}
}

// BaseURL returns the configured base URL string.
func (c *Client) BaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL.String()
}

// UserAgent returns the configured User-Agent string.
func (c *Client) UserAgent() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userAgent
}

// Timeout returns the configured timeout duration.
func (c *Client) Timeout() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timeout
}

// MaxAttempts returns the configured maximum retry attempts.
func (c *Client) MaxAttempts() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.maxAttempts
}

// RateLimit returns the current rate limit metadata parsed from responses.
func (c *Client) RateLimit() RateLimitInfo {
	if c.limiter == nil {
		return RateLimitInfo{}
	}
	return c.limiter.GetInfo()
}

// Close closes idle transport connections associated with the HTTP client.
func (c *Client) Close() error {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}
