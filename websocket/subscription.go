package websocket

import (
	"context"
	"fmt"
	"strings"
)

// Valid book depth values according to Kraken API documentation.
var validBookDepths = map[int]bool{
	10:   true,
	25:   true,
	100:  true,
	500:  true,
	1000: true,
}

// Valid OHLC interval values (in minutes) according to Kraken API documentation.
var validOHLCIntervals = map[int]bool{
	1:     true, // 1 minute
	5:     true, // 5 minutes
	15:    true, // 15 minutes
	30:    true, // 30 minutes
	60:    true, // 1 hour
	240:   true, // 4 hours
	1440:  true, // 1 day
	10080: true, // 1 week
	21600: true, // 15 days
}

// ValidateBookDepth validates that the book depth is a valid value.
func ValidateBookDepth(depth int) error {
	if !validBookDepths[depth] {
		return fmt.Errorf("invalid book depth %d: valid values are 10, 25, 100, 500, 1000", depth)
	}
	return nil
}

// ValidateOHLCInterval validates that the OHLC interval is a valid value.
func ValidateOHLCInterval(interval int) error {
	if !validOHLCIntervals[interval] {
		return fmt.Errorf("invalid OHLC interval %d: valid values are 1, 5, 15, 30, 60, 240, 1440, 10080, 21600", interval)
	}
	return nil
}

// ValidateSymbols validates that symbols are provided for channels that require them.
func ValidateSymbols(symbols []string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("at least one symbol is required")
	}
	for _, s := range symbols {
		if s == "" {
			return fmt.Errorf("symbol cannot be empty")
		}
	}
	return nil
}

// Subscription represents an active subscription.
type Subscription struct {
	Channel  Channel
	Symbols  []string
	Depth    int
	Interval int
	Snapshot bool
}

// subscriptionKey generates a unique key for a subscription.
func subscriptionKey(channel Channel, symbols []string, depth, interval int) string {
	return fmt.Sprintf("%s:%s:%d:%d", channel, strings.Join(symbols, ","), depth, interval)
}

// Subscribe subscribes to a channel.
func (c *Client) Subscribe(ctx context.Context, params SubscribeParams) error {
	reqID := c.nextReqID()

	req := SubscribeRequest{
		Method: MsgTypeSubscribe,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	if resp.Error != "" {
		return fmt.Errorf("subscribe error: %s", resp.Error)
	}

	// Track the subscription
	key := subscriptionKey(params.Channel, params.Symbol, params.Depth, params.Interval)
	c.subMu.Lock()
	c.subscriptions[key] = &Subscription{
		Channel:  params.Channel,
		Symbols:  params.Symbol,
		Depth:    params.Depth,
		Interval: params.Interval,
		Snapshot: params.Snapshot,
	}
	c.subMu.Unlock()

	return nil
}

// Unsubscribe unsubscribes from a channel.
func (c *Client) Unsubscribe(ctx context.Context, params UnsubscribeParams) error {
	reqID := c.nextReqID()

	req := UnsubscribeRequest{
		Method: MsgTypeUnsubscribe,
		ReqID:  reqID,
		Params: params,
	}

	resp, err := c.sendRequest(ctx, req, reqID)
	if err != nil {
		return fmt.Errorf("unsubscribe failed: %w", err)
	}

	if resp.Error != "" {
		return fmt.Errorf("unsubscribe error: %s", resp.Error)
	}

	// Remove the subscription
	key := subscriptionKey(params.Channel, params.Symbol, params.Depth, params.Interval)
	c.subMu.Lock()
	delete(c.subscriptions, key)
	c.subMu.Unlock()

	return nil
}

// SubscribeTicker subscribes to ticker updates for the specified symbols.
func (c *Client) SubscribeTicker(ctx context.Context, symbols []string) error {
	if err := ValidateSymbols(symbols); err != nil {
		return fmt.Errorf("ticker subscription: %w", err)
	}
	return c.Subscribe(ctx, SubscribeParams{
		Channel: ChannelTicker,
		Symbol:  symbols,
	})
}

// SubscribeBook subscribes to order book updates for the specified symbols.
// Depth options: 10, 25, 100, 500, 1000
func (c *Client) SubscribeBook(ctx context.Context, symbols []string, depth int) error {
	if err := ValidateSymbols(symbols); err != nil {
		return fmt.Errorf("book subscription: %w", err)
	}
	if err := ValidateBookDepth(depth); err != nil {
		return fmt.Errorf("book subscription: %w", err)
	}
	return c.Subscribe(ctx, SubscribeParams{
		Channel:  ChannelBook,
		Symbol:   symbols,
		Depth:    depth,
		Snapshot: true,
	})
}

// SubscribeTrade subscribes to trade updates for the specified symbols.
func (c *Client) SubscribeTrade(ctx context.Context, symbols []string) error {
	if err := ValidateSymbols(symbols); err != nil {
		return fmt.Errorf("trade subscription: %w", err)
	}
	return c.Subscribe(ctx, SubscribeParams{
		Channel:  ChannelTrade,
		Symbol:   symbols,
		Snapshot: true,
	})
}

// SubscribeOHLC subscribes to OHLC updates for the specified symbols.
// Interval options: 1, 5, 15, 30, 60, 240, 1440, 10080, 21600
func (c *Client) SubscribeOHLC(ctx context.Context, symbols []string, interval int) error {
	if err := ValidateSymbols(symbols); err != nil {
		return fmt.Errorf("ohlc subscription: %w", err)
	}
	if err := ValidateOHLCInterval(interval); err != nil {
		return fmt.Errorf("ohlc subscription: %w", err)
	}
	return c.Subscribe(ctx, SubscribeParams{
		Channel:  ChannelOHLC,
		Symbol:   symbols,
		Interval: interval,
		Snapshot: true,
	})
}

// SubscribeInstrument subscribes to instrument updates.
func (c *Client) SubscribeInstrument(ctx context.Context, snapshot bool) error {
	return c.Subscribe(ctx, SubscribeParams{
		Channel:  ChannelInstrument,
		Snapshot: snapshot,
	})
}

// SubscribeStatus subscribes to system status updates.
func (c *Client) SubscribeStatus(ctx context.Context) error {
	return c.Subscribe(ctx, SubscribeParams{
		Channel: ChannelStatus,
	})
}

// SubscribeHeartbeat subscribes to heartbeat messages.
func (c *Client) SubscribeHeartbeat(ctx context.Context) error {
	return c.Subscribe(ctx, SubscribeParams{
		Channel: ChannelHeartbeat,
	})
}

// SubscribeExecutions subscribes to execution updates (private).
// Requires authentication token.
func (c *Client) SubscribeExecutions(ctx context.Context, snapshot bool) error {
	if c.token == "" {
		return fmt.Errorf("authentication token required for executions channel")
	}
	return c.Subscribe(ctx, SubscribeParams{
		Channel:  ChannelExecutions,
		Snapshot: snapshot,
		Token:    c.token,
	})
}

// SubscribeBalances subscribes to balance updates (private).
// Requires authentication token.
func (c *Client) SubscribeBalances(ctx context.Context, snapshot bool) error {
	if c.token == "" {
		return fmt.Errorf("authentication token required for balances channel")
	}
	return c.Subscribe(ctx, SubscribeParams{
		Channel:  ChannelBalances,
		Snapshot: snapshot,
		Token:    c.token,
	})
}

// UnsubscribeTicker unsubscribes from ticker updates.
func (c *Client) UnsubscribeTicker(ctx context.Context, symbols []string) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelTicker,
		Symbol:  symbols,
	})
}

// UnsubscribeBook unsubscribes from order book updates.
func (c *Client) UnsubscribeBook(ctx context.Context, symbols []string, depth int) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelBook,
		Symbol:  symbols,
		Depth:   depth,
	})
}

// UnsubscribeTrade unsubscribes from trade updates.
func (c *Client) UnsubscribeTrade(ctx context.Context, symbols []string) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelTrade,
		Symbol:  symbols,
	})
}

// UnsubscribeOHLC unsubscribes from OHLC updates.
func (c *Client) UnsubscribeOHLC(ctx context.Context, symbols []string, interval int) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel:  ChannelOHLC,
		Symbol:   symbols,
		Interval: interval,
	})
}

// UnsubscribeInstrument unsubscribes from instrument updates.
func (c *Client) UnsubscribeInstrument(ctx context.Context) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelInstrument,
	})
}

// UnsubscribeStatus unsubscribes from system status updates.
func (c *Client) UnsubscribeStatus(ctx context.Context) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelStatus,
	})
}

// UnsubscribeHeartbeat unsubscribes from heartbeat messages.
func (c *Client) UnsubscribeHeartbeat(ctx context.Context) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelHeartbeat,
	})
}

// UnsubscribeExecutions unsubscribes from execution updates.
func (c *Client) UnsubscribeExecutions(ctx context.Context) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelExecutions,
		Token:   c.token,
	})
}

// UnsubscribeBalances unsubscribes from balance updates.
func (c *Client) UnsubscribeBalances(ctx context.Context) error {
	return c.Unsubscribe(ctx, UnsubscribeParams{
		Channel: ChannelBalances,
		Token:   c.token,
	})
}
