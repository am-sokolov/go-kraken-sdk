// Package websocket provides a WebSocket client for Kraken API v2.
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client for Kraken.
type Client struct {
	// URL is the WebSocket endpoint URL.
	url string
	// Token is the authentication token for private channels.
	token string

	// conn is the underlying WebSocket connection.
	conn *websocket.Conn
	// mu protects conn, token, and connection state (connected, reconnecting, closed).
	mu sync.RWMutex

	// Options
	reconnectInterval    time.Duration
	maxReconnectInterval time.Duration
	maxReconnectAttempts int // 0 means unlimited
	pingInterval         time.Duration
	readTimeout          time.Duration
	writeTimeout         time.Duration

	// Request ID counter
	reqID int64

	// Subscription management
	subscriptions map[string]*Subscription
	subMu         sync.RWMutex

	// Pending requests waiting for response
	pendingRequests map[int64]chan *BaseMessage
	pendingMu       sync.RWMutex

	// Event handlers
	onTicker     func([]TickerData)
	onBook       func([]BookData)
	onTrade      func([]TradeData)
	onOHLC       func([]OHLCData)
	onInstrument func([]InstrumentData)
	onStatus     func(*StatusData)
	onHeartbeat  func(*HeartbeatData)
	onExecution  func([]ExecutionData)
	onBalance    func([]BalanceData)
	onError      func(error)
	onConnect    func()
	onDisconnect func()
	// handlerMu protects event handlers from concurrent modification.
	handlerMu sync.RWMutex

	// State (protected by mu)
	connected    bool
	reconnecting bool
	closed       bool
	closeCh      chan struct{}

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// Option is a functional option for configuring the client.
type Option func(*Client)

// WithReconnectInterval sets the initial reconnect interval.
// The interval will increase exponentially up to maxReconnectInterval.
func WithReconnectInterval(d time.Duration) Option {
	return func(c *Client) {
		c.reconnectInterval = d
	}
}

// WithMaxReconnectInterval sets the maximum reconnect interval for exponential backoff.
func WithMaxReconnectInterval(d time.Duration) Option {
	return func(c *Client) {
		c.maxReconnectInterval = d
	}
}

// WithMaxReconnectAttempts sets the maximum number of reconnection attempts.
// Set to 0 for unlimited attempts (default).
func WithMaxReconnectAttempts(n int) Option {
	return func(c *Client) {
		c.maxReconnectAttempts = n
	}
}

// WithPingInterval sets the ping interval.
func WithPingInterval(d time.Duration) Option {
	return func(c *Client) {
		c.pingInterval = d
	}
}

// WithReadTimeout sets the read timeout.
func WithReadTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.readTimeout = d
	}
}

// WithWriteTimeout sets the write timeout.
func WithWriteTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.writeTimeout = d
	}
}

// New creates a new WebSocket client.
func New(url string, opts ...Option) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		url:                  url,
		reconnectInterval:    5 * time.Second,
		maxReconnectInterval: 60 * time.Second,
		maxReconnectAttempts: 0, // 0 means unlimited
		pingInterval:         30 * time.Second,
		readTimeout:          60 * time.Second,
		writeTimeout:         10 * time.Second,
		subscriptions:        make(map[string]*Subscription),
		pendingRequests:      make(map[int64]chan *BaseMessage),
		closeCh:              make(chan struct{}),
		ctx:                  ctx,
		cancel:               cancel,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewAuthenticated creates a new authenticated WebSocket client.
func NewAuthenticated(url, token string, opts ...Option) *Client {
	c := New(url, opts...)
	c.token = token
	return c
}

// Connect establishes the WebSocket connection.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	c.connected = true
	c.closed = false

	// Set up pong handler to track connection health
	conn.SetPongHandler(func(appData string) error {
		// Reset read deadline on pong received
		return conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	})

	// Start background goroutines
	go c.readLoop()
	go c.pingLoop()

	c.handlerMu.RLock()
	onConnect := c.onConnect
	c.handlerMu.RUnlock()
	if onConnect != nil {
		onConnect()
	}

	return nil
}

// Close closes the WebSocket connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	c.connected = false
	c.cancel()

	close(c.closeCh)

	// Clean up all pending requests
	c.cleanupPendingRequests()

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// cleanupPendingRequests signals all pending request channels to unblock waiters.
// Uses non-blocking send to avoid deadlock, then closes channels.
func (c *Client) cleanupPendingRequests() {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	for reqID, ch := range c.pendingRequests {
		// Non-blocking send of nil to signal closure, then close channel
		select {
		case ch <- nil:
		default:
		}
		close(ch)
		delete(c.pendingRequests, reqID)
	}
}

// IsConnected returns whether the client is connected.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// SetToken sets the authentication token.
func (c *Client) SetToken(token string) {
	c.mu.Lock()
	c.token = token
	c.mu.Unlock()
}

// nextReqID returns the next request ID.
func (c *Client) nextReqID() int64 {
	return atomic.AddInt64(&c.reqID, 1)
}

// sendMessage sends a message over the WebSocket connection.
func (c *Client) sendMessage(msg interface{}) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("not connected")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// sendRequest sends a request and waits for a response.
func (c *Client) sendRequest(ctx context.Context, msg interface{}, reqID int64) (*BaseMessage, error) {
	respCh := make(chan *BaseMessage, 1)

	c.pendingMu.Lock()
	c.pendingRequests[reqID] = respCh
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pendingRequests, reqID)
		c.pendingMu.Unlock()
	}()

	if err := c.sendMessage(msg); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp, ok := <-respCh:
		if !ok || resp == nil {
			return nil, fmt.Errorf("connection closed")
		}
		return resp, nil
	}
}

// readLoop reads messages from the WebSocket connection.
func (c *Client) readLoop() {
	defer func() {
		c.mu.Lock()
		c.connected = false
		shouldReconnect := !c.closed && !c.reconnecting
		c.mu.Unlock()

		// Clean up pending requests on disconnect
		c.cleanupPendingRequests()

		c.handlerMu.RLock()
		onDisconnect := c.onDisconnect
		c.handlerMu.RUnlock()
		if onDisconnect != nil {
			onDisconnect()
		}

		// Attempt reconnection if not closed
		if shouldReconnect {
			go c.reconnect()
		}
	}()

	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		if err := conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.handlerMu.RLock()
			onError := c.onError
			c.handlerMu.RUnlock()
			if onError != nil {
				onError(fmt.Errorf("failed to set read deadline: %w", err))
			}
			return
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			c.handlerMu.RLock()
			onError := c.onError
			c.handlerMu.RUnlock()
			if onError != nil {
				onError(err)
			}
			return
		}

		go c.handleMessage(data)
	}
}

// pingLoop sends periodic pings to keep the connection alive.
func (c *Client) pingLoop() {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeCh:
			return
		case <-ticker.C:
			if c.IsConnected() {
				c.sendMessage(&PingMessage{
					Method: MsgTypePing,
					ReqID:  c.nextReqID(),
				})
			}
		}
	}
}

// reconnect attempts to reconnect to the WebSocket with exponential backoff.
//
// Reconnection state machine:
// 1. Check if already reconnecting or closed - abort if so
// 2. Set reconnecting flag to prevent concurrent reconnection attempts
// 3. Loop with exponential backoff:
//   - Check if close signal received - abort if so
//   - Sleep for current interval (starts at reconnectInterval)
//   - Attempt connection with 10s timeout
//   - On success: resubscribe and return
//   - On failure: double interval (up to maxReconnectInterval), increment attempt count
//   - If maxReconnectAttempts reached (and > 0): give up
//
// 4. Clear reconnecting flag on exit
func (c *Client) reconnect() {
	c.mu.Lock()
	if c.reconnecting || c.closed {
		c.mu.Unlock()
		return
	}
	c.reconnecting = true
	maxAttempts := c.maxReconnectAttempts
	interval := c.reconnectInterval
	maxInterval := c.maxReconnectInterval
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.reconnecting = false
		c.mu.Unlock()
	}()

	attempts := 0
	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		time.Sleep(interval)
		attempts++

		ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
		err := c.Connect(ctx)
		cancel()

		if err == nil {
			// Resubscribe to all active subscriptions
			c.resubscribe()
			return
		}

		c.handlerMu.RLock()
		onError := c.onError
		c.handlerMu.RUnlock()
		if onError != nil {
			onError(fmt.Errorf("reconnection attempt %d failed: %w", attempts, err))
		}

		// Check if max attempts reached
		if maxAttempts > 0 && attempts >= maxAttempts {
			if onError != nil {
				onError(fmt.Errorf("max reconnection attempts (%d) reached, giving up", maxAttempts))
			}
			return
		}

		// Exponential backoff: double the interval up to maxInterval
		interval *= 2
		if interval > maxInterval {
			interval = maxInterval
		}
	}
}

// resubscribe resubscribes to all active subscriptions after reconnection.
func (c *Client) resubscribe() {
	c.subMu.RLock()
	subs := make([]*Subscription, 0, len(c.subscriptions))
	for _, sub := range c.subscriptions {
		subs = append(subs, sub)
	}
	c.subMu.RUnlock()

	// Get token under lock
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	for _, sub := range subs {
		ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
		params := SubscribeParams{
			Channel:  sub.Channel,
			Symbol:   sub.Symbols,
			Depth:    sub.Depth,
			Interval: sub.Interval,
			Snapshot: sub.Snapshot,
		}
		if token != "" {
			params.Token = token
		}
		if err := c.Subscribe(ctx, params); err != nil {
			c.handlerMu.RLock()
			onError := c.onError
			c.handlerMu.RUnlock()
			if onError != nil {
				onError(fmt.Errorf("resubscribe failed for %s: %w", sub.Channel, err))
			}
		}
		cancel()
	}
}

// handleMessage processes an incoming WebSocket message.
func (c *Client) handleMessage(data []byte) {
	var msg BaseMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		c.handlerMu.RLock()
		onError := c.onError
		c.handlerMu.RUnlock()
		if onError != nil {
			onError(fmt.Errorf("failed to parse message: %w", err))
		}
		return
	}

	// Handle response to pending request
	if msg.ReqID != 0 {
		c.pendingMu.RLock()
		respCh, ok := c.pendingRequests[msg.ReqID]
		c.pendingMu.RUnlock()
		if ok {
			// Use non-blocking send to avoid panic if channel was closed during cleanup
			select {
			case respCh <- &msg:
			default:
				// Channel full or closed, response will be dropped
			}
			return
		}
	}

	// Handle channel messages
	switch msg.Channel {
	case ChannelTicker:
		c.handleTickerMessage(msg.Data)
	case ChannelBook:
		c.handleBookMessage(msg.Data)
	case ChannelTrade:
		c.handleTradeMessage(msg.Data)
	case ChannelOHLC:
		c.handleOHLCMessage(msg.Data)
	case ChannelInstrument:
		c.handleInstrumentMessage(msg.Data)
	case ChannelStatus:
		c.handleStatusMessage(msg.Data)
	case ChannelHeartbeat:
		c.handleHeartbeatMessage(msg.Data)
	case ChannelExecutions:
		c.handleExecutionMessage(msg.Data)
	case ChannelBalances:
		c.handleBalanceMessage(msg.Data)
	}
}

func (c *Client) handleTickerMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onTicker := c.onTicker
	onError := c.onError
	c.handlerMu.RUnlock()

	if onTicker == nil {
		return
	}
	var tickers []TickerData
	if err := json.Unmarshal(data, &tickers); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse ticker: %w", err))
		}
		return
	}
	onTicker(tickers)
}

func (c *Client) handleBookMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onBook := c.onBook
	onError := c.onError
	c.handlerMu.RUnlock()

	if onBook == nil {
		return
	}
	var books []BookData
	if err := json.Unmarshal(data, &books); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse book: %w", err))
		}
		return
	}
	onBook(books)
}

func (c *Client) handleTradeMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onTrade := c.onTrade
	onError := c.onError
	c.handlerMu.RUnlock()

	if onTrade == nil {
		return
	}
	var trades []TradeData
	if err := json.Unmarshal(data, &trades); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse trade: %w", err))
		}
		return
	}
	onTrade(trades)
}

func (c *Client) handleOHLCMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onOHLC := c.onOHLC
	onError := c.onError
	c.handlerMu.RUnlock()

	if onOHLC == nil {
		return
	}
	var ohlc []OHLCData
	if err := json.Unmarshal(data, &ohlc); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse ohlc: %w", err))
		}
		return
	}
	onOHLC(ohlc)
}

func (c *Client) handleInstrumentMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onInstrument := c.onInstrument
	onError := c.onError
	c.handlerMu.RUnlock()

	if onInstrument == nil {
		return
	}

	// Instrument data can come as either an array or a single object with "assets" field
	var instruments []InstrumentData
	if err := json.Unmarshal(data, &instruments); err == nil {
		onInstrument(instruments)
		return
	}

	// Try single object (snapshot may come as object with nested data)
	var singleInstrument InstrumentData
	if err := json.Unmarshal(data, &singleInstrument); err == nil {
		onInstrument([]InstrumentData{singleInstrument})
		return
	}

	// Try object with assets array
	var wrapper struct {
		Assets []InstrumentData `json:"assets"`
	}
	if err := json.Unmarshal(data, &wrapper); err == nil && len(wrapper.Assets) > 0 {
		onInstrument(wrapper.Assets)
		return
	}

	if onError != nil {
		onError(fmt.Errorf("failed to parse instrument: unrecognized format"))
	}
}

func (c *Client) handleStatusMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onStatus := c.onStatus
	onError := c.onError
	c.handlerMu.RUnlock()

	if onStatus == nil {
		return
	}

	// Status can come as either a single object or an array
	// Try array first (most common)
	var statusArr []StatusData
	if err := json.Unmarshal(data, &statusArr); err == nil && len(statusArr) > 0 {
		onStatus(&statusArr[0])
		return
	}

	// Try single object
	var status StatusData
	if err := json.Unmarshal(data, &status); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse status: %w", err))
		}
		return
	}
	onStatus(&status)
}

func (c *Client) handleHeartbeatMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onHeartbeat := c.onHeartbeat
	c.handlerMu.RUnlock()

	if onHeartbeat == nil {
		return
	}

	// Heartbeat might be empty or have minimal data
	if len(data) == 0 || string(data) == "null" || string(data) == "{}" {
		onHeartbeat(&HeartbeatData{})
		return
	}

	var heartbeat HeartbeatData
	if err := json.Unmarshal(data, &heartbeat); err != nil {
		// On parse error, still call handler with empty data since heartbeat was received
		onHeartbeat(&HeartbeatData{})
		return
	}
	onHeartbeat(&heartbeat)
}

func (c *Client) handleExecutionMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onExecution := c.onExecution
	onError := c.onError
	c.handlerMu.RUnlock()

	if onExecution == nil {
		return
	}
	var executions []ExecutionData
	if err := json.Unmarshal(data, &executions); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse execution: %w", err))
		}
		return
	}
	onExecution(executions)
}

func (c *Client) handleBalanceMessage(data json.RawMessage) {
	c.handlerMu.RLock()
	onBalance := c.onBalance
	onError := c.onError
	c.handlerMu.RUnlock()

	if onBalance == nil {
		return
	}
	var balances []BalanceData
	if err := json.Unmarshal(data, &balances); err != nil {
		if onError != nil {
			onError(fmt.Errorf("failed to parse balance: %w", err))
		}
		return
	}
	onBalance(balances)
}

// Event handler setters

// OnTicker sets the handler for ticker updates.
func (c *Client) OnTicker(fn func([]TickerData)) {
	c.handlerMu.Lock()
	c.onTicker = fn
	c.handlerMu.Unlock()
}

// OnBook sets the handler for order book updates.
func (c *Client) OnBook(fn func([]BookData)) {
	c.handlerMu.Lock()
	c.onBook = fn
	c.handlerMu.Unlock()
}

// OnTrade sets the handler for trade updates.
func (c *Client) OnTrade(fn func([]TradeData)) {
	c.handlerMu.Lock()
	c.onTrade = fn
	c.handlerMu.Unlock()
}

// OnOHLC sets the handler for OHLC updates.
func (c *Client) OnOHLC(fn func([]OHLCData)) {
	c.handlerMu.Lock()
	c.onOHLC = fn
	c.handlerMu.Unlock()
}

// OnInstrument sets the handler for instrument updates.
func (c *Client) OnInstrument(fn func([]InstrumentData)) {
	c.handlerMu.Lock()
	c.onInstrument = fn
	c.handlerMu.Unlock()
}

// OnStatus sets the handler for system status updates.
func (c *Client) OnStatus(fn func(*StatusData)) {
	c.handlerMu.Lock()
	c.onStatus = fn
	c.handlerMu.Unlock()
}

// OnHeartbeat sets the handler for heartbeat updates.
func (c *Client) OnHeartbeat(fn func(*HeartbeatData)) {
	c.handlerMu.Lock()
	c.onHeartbeat = fn
	c.handlerMu.Unlock()
}

// OnExecution sets the handler for execution updates.
func (c *Client) OnExecution(fn func([]ExecutionData)) {
	c.handlerMu.Lock()
	c.onExecution = fn
	c.handlerMu.Unlock()
}

// OnBalance sets the handler for balance updates.
func (c *Client) OnBalance(fn func([]BalanceData)) {
	c.handlerMu.Lock()
	c.onBalance = fn
	c.handlerMu.Unlock()
}

// OnError sets the handler for errors.
func (c *Client) OnError(fn func(error)) {
	c.handlerMu.Lock()
	c.onError = fn
	c.handlerMu.Unlock()
}

// OnConnect sets the handler for successful connection.
func (c *Client) OnConnect(fn func()) {
	c.handlerMu.Lock()
	c.onConnect = fn
	c.handlerMu.Unlock()
}

// OnDisconnect sets the handler for disconnection.
func (c *Client) OnDisconnect(fn func()) {
	c.handlerMu.Lock()
	c.onDisconnect = fn
	c.handlerMu.Unlock()
}
