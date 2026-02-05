# Kraken Go SDK

A comprehensive Go SDK for the Kraken cryptocurrency exchange REST API and WebSocket v2 API.

## Features

- **REST API**: Complete coverage of ~50 REST endpoints
  - Public endpoints (no authentication): Server time, system status, assets, pairs, ticker, OHLC, order book, trades, spreads
  - Account endpoints: Balance, trade balance, open/closed orders, trades history, ledgers, positions
  - Trading endpoints: Add/edit/amend/cancel orders, batch operations, dead man's switch
  - Funding endpoints: Deposit/withdrawal methods, addresses, status, wallet transfers
  - Earn endpoints: Staking strategies, allocations

- **WebSocket v2 API**: Real-time market data and trading
  - Public channels: ticker, book, trade, ohlc, instrument, status
  - Private channels: executions, balances
  - Trading methods: add_order, edit_order, amend_order, cancel_order, batch operations

- **Authentication**: HMAC-SHA512 signature generation with nonce management
- **Rate Limiting**: Token bucket implementation respecting Kraken's tier-based limits
- **Type Safety**: Strongly typed requests/responses with decimal precision for financial calculations

## Installation

```bash
go get github.com/am-sokolov/go-kraken-sdk
```

## Quick Start

### REST API

```go
package main

import (
    "context"
    "fmt"
    "log"

    kraken "github.com/am-sokolov/go-kraken-sdk"
    "github.com/am-sokolov/go-kraken-sdk/types"
)

func main() {
    // Create client with API credentials
    client, err := kraken.New(
        kraken.WithAPIKey("your-api-key", "your-api-secret"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Public endpoints (no auth required)
    time, err := client.Public.GetServerTime(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Server time: %d\n", time.UnixTime)

    // Get ticker
    tickers, err := client.Public.GetTicker(ctx, []string{"XBTUSD"})
    if err != nil {
        log.Fatal(err)
    }
    for pair, ticker := range tickers {
        fmt.Printf("%s: Last=%v\n", pair, ticker.Close[0])
    }

    // Private endpoints (auth required)
    balance, err := client.Account.GetBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for asset, amount := range balance {
        fmt.Printf("%s: %s\n", asset, amount)
    }

    // Place an order (validate only)
    result, err := client.Trading.AddOrder(ctx, &kraken.AddOrderRequest{
        OrderType: types.OrderTypeLimit,
        Side:      types.SideBuy,
        Volume:    "0.001",
        Pair:      "XBTUSD",
        Price:     "50000",
        Validate:  true, // Validate only, don't actually place
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order: %s\n", result.Description.Order)
}
```

### WebSocket API

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    kraken "github.com/am-sokolov/go-kraken-sdk"
    "github.com/am-sokolov/go-kraken-sdk/websocket"
)

func main() {
    // Create client for WebSocket
    client, _ := kraken.New()

    // Create public WebSocket client
    ws := client.WebSocket()

    // Set up handlers
    ws.OnTicker(func(tickers []websocket.TickerData) {
        for _, t := range tickers {
            fmt.Printf("Ticker %s: Bid=%v Ask=%v Last=%v\n",
                t.Symbol, t.Bid, t.Ask, t.Last)
        }
    })

    ws.OnTrade(func(trades []websocket.TradeData) {
        for _, t := range trades {
            fmt.Printf("Trade %s: %s %v @ %v\n",
                t.Symbol, t.Side, t.Qty, t.Price)
        }
    })

    ws.OnError(func(err error) {
        log.Printf("WebSocket error: %v", err)
    })

    ws.OnConnect(func() {
        log.Println("Connected to WebSocket")
    })

    ws.OnDisconnect(func() {
        log.Println("Disconnected from WebSocket")
    })

    // Connect
    ctx := context.Background()
    if err := ws.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    // Subscribe to channels
    ws.SubscribeTicker(ctx, []string{"BTC/USD", "ETH/USD"})
    ws.SubscribeTrade(ctx, []string{"BTC/USD"})
    ws.SubscribeBook(ctx, []string{"BTC/USD"}, 10)

    // Run for a while
    time.Sleep(30 * time.Second)
}
```

### Authenticated WebSocket

```go
// Create client with API credentials
client, _ := kraken.New(
    kraken.WithAPIKey("your-api-key", "your-api-secret"),
)

ctx := context.Background()

// Create authenticated WebSocket client (fetches token automatically)
ws, err := client.WebSocketAuth(ctx)
if err != nil {
    log.Fatal(err)
}

// Set up execution handler
ws.OnExecution(func(executions []websocket.ExecutionData) {
    for _, e := range executions {
        fmt.Printf("Execution: OrderID=%s Status=%s\n",
            e.OrderID, e.OrderStatus)
    }
})

// Connect and subscribe
if err := ws.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer ws.Close()

ws.SubscribeExecutions(ctx, true) // true for snapshot

// Place order via WebSocket
result, err := ws.AddOrder(ctx, websocket.AddOrderParams{
    Symbol:    "BTC/USD",
    Side:      "buy",
    OrderType: "limit",
    OrderQty:  "0.001",
    LimitPrice: "50000",
    Validate:  true,
})
```

## Configuration Options

```go
client, err := kraken.New(
    kraken.WithAPIKey("key", "secret"),
    kraken.WithBaseURL("https://api.kraken.com"), // Custom base URL
    kraken.WithTimeout(30*time.Second),
    kraken.WithRateLimitTier(kraken.TierIntermediate),
    kraken.WithUserAgent("my-app/1.0"),
    kraken.With2FA(func() string { return getOTP() }),
    kraken.WithLogger(myLogger),
)
```

## Rate Limiting

The SDK includes built-in rate limiting that respects Kraken's tier-based limits:

| Tier | Max Counter | Decay/sec | Max Orders/Batch |
|------|-------------|-----------|------------------|
| Starter | 60 | 1.0 | 60 |
| Intermediate | 125 | 2.34 | 80 |
| Pro | 180 | 3.75 | 225 |

```go
import "github.com/am-sokolov/go-kraken-sdk/ratelimit"

// Create rate limiter
limiter := ratelimit.NewCombinedLimiter(ratelimit.TierStarter)

// Check if request allowed
if limiter.AllowAPI("account") {
    // Make request
}

// Or wait for capacity
ctx := context.Background()
if err := limiter.WaitAPI(ctx, "account"); err != nil {
    // Context cancelled
}
```

## Error Handling

```go
result, err := client.Trading.AddOrder(ctx, req)
if err != nil {
    if apiErr, ok := err.(*kraken.APIError); ok {
        fmt.Printf("API Error: Category=%s, Messages=%v\n",
            apiErr.Category, apiErr.Messages)

        // Check error type
        if kraken.IsAuthError(err) {
            // Handle authentication error
        }
        if kraken.IsRetryable(err) {
            // Retry the request
        }
    }
}
```

## Project Structure

```
kraken-sdk/
├── kraken.go          # Main client
├── options.go         # Configuration options
├── errors.go          # Error types
├── constants.go       # API URLs, constants
├── account.go         # Account service
├── trading.go         # Trading service
├── funding.go         # Funding service
├── earn.go            # Earn service
├── subaccount.go      # Subaccount service
├── auth/              # Authentication
│   ├── auth.go
│   ├── signer.go      # HMAC-SHA512 signatures
│   └── nonce.go       # Nonce generation
├── rest/              # REST client
│   ├── client.go
│   ├── public.go      # Public endpoints
│   └── response.go
├── websocket/         # WebSocket client
│   ├── client.go
│   ├── subscription.go
│   ├── trading.go
│   └── message.go
├── ratelimit/         # Rate limiting
│   └── limiter.go
└── types/             # Data types
    ├── order.go
    ├── asset.go
    ├── ticker.go
    └── ...
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./auth/...
go test -v ./websocket/...
```

## Dependencies

- `github.com/gorilla/websocket` - WebSocket client
- `github.com/shopspring/decimal` - Precise decimal arithmetic

## License

MIT License

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Disclaimer

This SDK is not officially affiliated with Kraken. Use at your own risk. Always test with small amounts first and use the `validate=true` parameter when testing order placement.
