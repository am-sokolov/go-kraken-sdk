package kraken

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/am-sokolov/go-kraken-sdk/auth"
	"github.com/am-sokolov/go-kraken-sdk/rest"
	"github.com/am-sokolov/go-kraken-sdk/websocket"
)

// Version is the SDK version.
const Version = "1.0.0"

// Client is the main entry point for the Kraken SDK.
// It provides access to both REST and WebSocket APIs.
type Client struct {
	options *Options

	// REST client
	restClient *rest.Client

	// REST API services
	Public     *rest.PublicService
	Account    *AccountService
	Trading    *TradingService
	Funding    *FundingService
	Earn       *EarnService
	Subaccount *SubaccountService

	// Authentication
	auth *auth.Authenticator

	// Internal state
	mu     sync.RWMutex
	closed bool
}

// New creates a new Kraken client with the given options.
func New(opts ...Option) (*Client, error) {
	options := applyOptions(opts)

	// Create HTTP client if not provided
	httpClient := options.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: options.Timeout,
		}
	}

	c := &Client{
		options: options,
	}

	// Create authenticator if credentials provided
	var clientOpts []rest.ClientOption
	if options.APIKey != "" && options.APISecret != "" {
		authenticator, err := auth.NewAuthenticator(options.APIKey, options.APISecret)
		if err != nil {
			return nil, err
		}
		c.auth = authenticator
		clientOpts = append(clientOpts, rest.WithAuth(authenticator))
	}

	// Add other client options
	if options.UserAgent != "" {
		clientOpts = append(clientOpts, rest.WithUserAgent(options.UserAgent))
	}
	if options.OnError != nil {
		clientOpts = append(clientOpts, rest.WithOnError(options.OnError))
	}
	if options.DebugWriter != nil {
		clientOpts = append(clientOpts, rest.WithDebugWriter(options.DebugWriter))
	}

	// Create REST client
	c.restClient = rest.NewClient(options.BaseURL, httpClient, clientOpts...)

	// Initialize REST services
	c.Public = rest.NewPublicService(c.restClient)
	c.Account = newAccountService(c)
	c.Trading = newTradingService(c)
	c.Funding = newFundingService(c)
	c.Earn = newEarnService(c)
	c.Subaccount = newSubaccountService(c)

	return c, nil
}

// Close gracefully shuts down the client and releases resources.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true

	return nil
}

// Options returns the client options.
func (c *Client) Options() *Options {
	return c.options
}

// RESTClient returns the underlying REST client.
func (c *Client) RESTClient() *rest.Client {
	return c.restClient
}

// IsAuthenticated returns true if the client has API credentials configured.
func (c *Client) IsAuthenticated() bool {
	return c.auth != nil
}

// WebSocket creates a new public WebSocket client.
// The returned client connects to the public WebSocket endpoint for market data.
// Call Connect() on the returned client to establish the connection.
func (c *Client) WebSocket() *websocket.Client {
	return websocket.New(
		c.options.PublicWSURL,
		websocket.WithReconnectInterval(c.options.WSReconnectInterval),
		websocket.WithPingInterval(c.options.WSPingInterval),
		websocket.WithReadTimeout(c.options.WSReadTimeout),
		websocket.WithWriteTimeout(c.options.WSWriteTimeout),
	)
}

// WebSocketAuth creates a new authenticated WebSocket client.
// The client will first obtain a WebSocket token via REST API,
// then return a configured client. Call Connect() on the returned client
// to establish the connection.
func (c *Client) WebSocketAuth(ctx context.Context) (*websocket.Client, error) {
	if !c.IsAuthenticated() {
		return nil, &APIError{
			Category: "EAuth",
			Messages: []string{"API credentials required for authenticated WebSocket"},
		}
	}

	// Get WebSocket token
	token, err := c.Account.GetWebSocketsToken(ctx)
	if err != nil {
		return nil, err
	}

	return websocket.NewAuthenticated(
		c.options.PrivateWSURL,
		token.Token,
		websocket.WithReconnectInterval(c.options.WSReconnectInterval),
		websocket.WithPingInterval(c.options.WSPingInterval),
		websocket.WithReadTimeout(c.options.WSReadTimeout),
		websocket.WithWriteTimeout(c.options.WSWriteTimeout),
	), nil
}

// AccountService provides access to private account REST API endpoints.
type AccountService struct {
	client *Client
}

func newAccountService(c *Client) *AccountService {
	return &AccountService{client: c}
}

// TradingService provides access to private trading REST API endpoints.
type TradingService struct {
	client *Client
}

func newTradingService(c *Client) *TradingService {
	return &TradingService{client: c}
}

// FundingService provides access to private funding REST API endpoints.
type FundingService struct {
	client *Client
}

func newFundingService(c *Client) *FundingService {
	return &FundingService{client: c}
}

// EarnService provides access to private earn REST API endpoints.
type EarnService struct {
	client *Client
}

func newEarnService(c *Client) *EarnService {
	return &EarnService{client: c}
}

// SubaccountService provides access to subaccount REST API endpoints.
type SubaccountService struct {
	client *Client
}

func newSubaccountService(c *Client) *SubaccountService {
	return &SubaccountService{client: c}
}

// logDebug logs a debug message if a logger is configured.
func (c *Client) logDebug(msg string, args ...any) {
	if c.options.Logger != nil {
		c.options.Logger.Debug(msg, args...)
	}
}

// logError logs an error message if a logger is configured.
func (c *Client) logError(msg string, args ...any) {
	if c.options.Logger != nil {
		c.options.Logger.Error(msg, args...)
	}
}

// handleError invokes the error callback if configured.
func (c *Client) handleError(err error) {
	if c.options.OnError != nil {
		c.options.OnError(err)
	}
}

// now returns the current time. This is a method to allow for testing.
func (c *Client) now() time.Time {
	return time.Now()
}
