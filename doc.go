// Package kraken provides a Go SDK for the Kraken cryptocurrency exchange API.
//
// The SDK supports both the REST API and WebSocket v2 API with full type safety,
// proper authentication, rate limiting, and comprehensive error handling.
//
// # Quick Start
//
// Create a client and make requests:
//
//	client, err := kraken.New()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Get server time (public endpoint)
//	time, err := client.Public.GetServerTime(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Server time: %v\n", time)
//
// # Authentication
//
// For private endpoints, provide API credentials:
//
//	client, err := kraken.New(
//	    kraken.WithAPIKey("your-api-key", "your-api-secret"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get account balance (private endpoint)
//	balance, err := client.Account.GetBalance(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Balance: %+v\n", balance)
//
// # WebSocket
//
// For real-time data, use the WebSocket client:
//
//	ws, err := client.WebSocket(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer ws.Close()
//
//	// Subscribe to ticker updates
//	sub, err := ws.SubscribeTicker(ctx, []string{"BTC/USD", "ETH/USD"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for msg := range sub.Messages() {
//	    fmt.Printf("Ticker update: %+v\n", msg)
//	}
//
// # Error Handling
//
// The SDK returns structured errors that can be inspected:
//
//	_, err := client.Trading.AddOrder(ctx, &AddOrderRequest{...})
//	var apiErr *kraken.APIError
//	if errors.As(err, &apiErr) {
//	    if errors.Is(apiErr, kraken.ErrInsufficientFunds) {
//	        // Handle insufficient funds
//	    }
//	}
//
// # Rate Limiting
//
// The SDK includes built-in rate limiting that respects Kraken's limits:
//
//	client, err := kraken.New(
//	    kraken.WithRateLimitTier(kraken.TierPro),
//	)
//
// For more information, see:
//   - REST API: https://docs.kraken.com/api/docs/rest-api
//   - WebSocket API: https://docs.kraken.com/api/docs/websocket-v2
package kraken
