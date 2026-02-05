package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestClient_Connect(t *testing.T) {
	// Create a test WebSocket server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		// Send a status message
		msg := map[string]interface{}{
			"channel": "status",
			"type":    "update",
			"data": map[string]interface{}{
				"system":        "online",
				"version":       "2.0.0",
				"api_version":   "v2",
				"connection_id": 12345,
			},
		}
		conn.WriteJSON(msg)

		// Keep connection alive
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !client.IsConnected() {
		t.Error("client should be connected")
	}

	err = client.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if client.IsConnected() {
		t.Error("client should be disconnected after close")
	}
}

func TestClient_Subscribe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var req SubscribeRequest
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			// Send subscription response
			resp := map[string]interface{}{
				"method":  "subscribe",
				"req_id":  req.ReqID,
				"success": true,
				"result": map[string]interface{}{
					"channel": req.Params.Channel,
					"symbol":  req.Params.Symbol,
				},
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	err := client.SubscribeTicker(ctx, []string{"BTC/USD"})
	if err != nil {
		t.Fatalf("SubscribeTicker failed: %v", err)
	}
}

func TestClient_TickerHandler(t *testing.T) {
	receivedTicker := make(chan []TickerData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Send ticker data
		msg := map[string]interface{}{
			"channel": "ticker",
			"type":    "update",
			"data": []map[string]interface{}{
				{
					"symbol":     "BTC/USD",
					"bid":        "50000.00",
					"bid_qty":    "1.0",
					"ask":        "50001.00",
					"ask_qty":    "1.0",
					"last":       "50000.50",
					"volume":     "1000.0",
					"vwap":       "49500.00",
					"low":        "48000.00",
					"high":       "51000.00",
					"change":     "1500.00",
					"change_pct": "3.0",
				},
			},
		}
		conn.WriteJSON(msg)

		// Keep connection alive
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnTicker(func(tickers []TickerData) {
		receivedTicker <- tickers
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case tickers := <-receivedTicker:
		if len(tickers) != 1 {
			t.Errorf("expected 1 ticker, got %d", len(tickers))
		}
		if tickers[0].Symbol != "BTC/USD" {
			t.Errorf("expected symbol BTC/USD, got %s", tickers[0].Symbol)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for ticker data")
	}
}

func TestClient_TradeHandler(t *testing.T) {
	receivedTrades := make(chan []TradeData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Send trade data
		msg := map[string]interface{}{
			"channel": "trade",
			"type":    "update",
			"data": []map[string]interface{}{
				{
					"symbol":    "BTC/USD",
					"side":      "buy",
					"price":     "50000.00",
					"qty":       "0.5",
					"ord_type":  "market",
					"trade_id":  12345,
					"timestamp": "2024-01-01T00:00:00Z",
				},
			},
		}
		conn.WriteJSON(msg)

		// Keep connection alive
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnTrade(func(trades []TradeData) {
		receivedTrades <- trades
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case trades := <-receivedTrades:
		if len(trades) != 1 {
			t.Errorf("expected 1 trade, got %d", len(trades))
		}
		if trades[0].Symbol != "BTC/USD" {
			t.Errorf("expected symbol BTC/USD, got %s", trades[0].Symbol)
		}
		if trades[0].Side != "buy" {
			t.Errorf("expected side buy, got %s", trades[0].Side)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for trade data")
	}
}

func TestClient_NewAuthenticated(t *testing.T) {
	client := NewAuthenticated("wss://test.example.com", "test-token")

	if client.url != "wss://test.example.com" {
		t.Errorf("expected URL wss://test.example.com, got %s", client.url)
	}

	if client.token != "test-token" {
		t.Errorf("expected token test-token, got %s", client.token)
	}
}

func TestClient_Options(t *testing.T) {
	client := New(
		"wss://test.example.com",
		WithReconnectInterval(10*time.Second),
		WithPingInterval(20*time.Second),
		WithReadTimeout(30*time.Second),
		WithWriteTimeout(5*time.Second),
	)

	if client.reconnectInterval != 10*time.Second {
		t.Errorf("expected reconnect interval 10s, got %s", client.reconnectInterval)
	}

	if client.pingInterval != 20*time.Second {
		t.Errorf("expected ping interval 20s, got %s", client.pingInterval)
	}

	if client.readTimeout != 30*time.Second {
		t.Errorf("expected read timeout 30s, got %s", client.readTimeout)
	}

	if client.writeTimeout != 5*time.Second {
		t.Errorf("expected write timeout 5s, got %s", client.writeTimeout)
	}
}

func TestSubscriptionKey(t *testing.T) {
	key := subscriptionKey(ChannelTicker, []string{"BTC/USD", "ETH/USD"}, 0, 0)
	expected := "ticker:BTC/USD,ETH/USD:0:0"
	if key != expected {
		t.Errorf("expected key %s, got %s", expected, key)
	}

	key = subscriptionKey(ChannelBook, []string{"BTC/USD"}, 10, 0)
	expected = "book:BTC/USD:10:0"
	if key != expected {
		t.Errorf("expected key %s, got %s", expected, key)
	}

	key = subscriptionKey(ChannelOHLC, []string{"BTC/USD"}, 0, 60)
	expected = "ohlc:BTC/USD:0:60"
	if key != expected {
		t.Errorf("expected key %s, got %s", expected, key)
	}
}

func TestClient_BookHandler(t *testing.T) {
	receivedBook := make(chan []BookData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "book",
			"type":    "snapshot",
			"data": []map[string]interface{}{
				{
					"symbol": "BTC/USD",
					"bids": []map[string]interface{}{
						{"price": "50000.00", "qty": "1.0"},
					},
					"asks": []map[string]interface{}{
						{"price": "50001.00", "qty": "0.5"},
					},
				},
			},
		}
		conn.WriteJSON(msg)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnBook(func(books []BookData) {
		receivedBook <- books
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case books := <-receivedBook:
		if len(books) != 1 {
			t.Errorf("expected 1 book, got %d", len(books))
		}
		if books[0].Symbol != "BTC/USD" {
			t.Errorf("expected symbol BTC/USD, got %s", books[0].Symbol)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for book data")
	}
}

func TestClient_OnMessage_Parsed(t *testing.T) {
	received := make(chan MessageEvent, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "ticker",
			"type":    "update",
			"data": []map[string]interface{}{
				{
					"symbol":     "BTC/USD",
					"bid":        "50000.00",
					"bid_qty":    "1.0",
					"ask":        "50001.00",
					"ask_qty":    "1.0",
					"last":       "50000.50",
					"volume":     "1000.0",
					"vwap":       "49500.00",
					"low":        "48000.00",
					"high":       "51000.00",
					"change":     "1500.00",
					"change_pct": "3.0",
				},
			},
		}
		conn.WriteJSON(msg)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)
	client.OnMessage(func(ev MessageEvent) {
		received <- ev
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case ev := <-received:
		if !ev.Parsed {
			t.Fatalf("expected Parsed=true, got false (err=%q)", ev.ParseError)
		}
		if len(ev.Raw) == 0 {
			t.Error("expected raw payload")
		}
		if ev.ReceivedAt.IsZero() {
			t.Error("expected ReceivedAt to be set")
		}
		if ev.ReceivedMonoNs < 0 {
			t.Errorf("expected non-negative ReceivedMonoNs, got %d", ev.ReceivedMonoNs)
		}
		if ev.Message.Channel != ChannelTicker {
			t.Errorf("expected channel %q, got %q", ChannelTicker, ev.Message.Channel)
		}
		if ev.Message.Type != "update" {
			t.Errorf("expected type update, got %q", ev.Message.Type)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for OnMessage event")
	}
}

func TestClient_OnMessage_ParseError(t *testing.T) {
	received := make(chan MessageEvent, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		conn.WriteMessage(websocket.TextMessage, []byte("not json"))

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)
	client.OnMessage(func(ev MessageEvent) {
		received <- ev
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case ev := <-received:
		if ev.Parsed {
			t.Fatalf("expected Parsed=false, got true: %+v", ev.Message)
		}
		if ev.ParseError == "" {
			t.Fatal("expected ParseError")
		}
		if len(ev.Raw) == 0 {
			t.Error("expected raw payload")
		}
		if ev.ReceivedAt.IsZero() {
			t.Error("expected ReceivedAt to be set")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for OnMessage parse error event")
	}
}

func TestClient_OHLCHandler(t *testing.T) {
	receivedOHLC := make(chan []OHLCData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "ohlc",
			"type":    "update",
			"data": []map[string]interface{}{
				{
					"symbol":   "BTC/USD",
					"open":     "49000.00",
					"high":     "51000.00",
					"low":      "48500.00",
					"close":    "50500.00",
					"volume":   "100.5",
					"vwap":     "50000.00",
					"trades":   150,
					"interval": 60,
				},
			},
		}
		conn.WriteJSON(msg)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnOHLC(func(ohlc []OHLCData) {
		receivedOHLC <- ohlc
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case ohlc := <-receivedOHLC:
		if len(ohlc) != 1 {
			t.Errorf("expected 1 ohlc, got %d", len(ohlc))
		}
		if ohlc[0].Symbol != "BTC/USD" {
			t.Errorf("expected symbol BTC/USD, got %s", ohlc[0].Symbol)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for ohlc data")
	}
}

func TestClient_StatusHandler(t *testing.T) {
	receivedStatus := make(chan *StatusData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "status",
			"type":    "update",
			"data": map[string]interface{}{
				"system":        "online",
				"version":       "2.0.0",
				"api_version":   "v2",
				"connection_id": 12345,
			},
		}
		conn.WriteJSON(msg)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnStatus(func(status *StatusData) {
		receivedStatus <- status
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case status := <-receivedStatus:
		if status.System != "online" {
			t.Errorf("expected system online, got %s", status.System)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for status data")
	}
}

func TestClient_ErrorHandler(t *testing.T) {
	receivedError := make(chan error, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Close immediately to trigger error
		conn.Close()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnError(func(err error) {
		receivedError <- err
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case err := <-receivedError:
		if err == nil {
			t.Error("expected error, got nil")
		}
	case <-time.After(3 * time.Second):
		// May not receive error if connection closed cleanly
	}
}

func TestClient_ConnectDisconnectCallbacks(t *testing.T) {
	connectCalled := make(chan struct{}, 1)
	disconnectCalled := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	client.OnConnect(func() {
		connectCalled <- struct{}{}
	})

	client.OnDisconnect(func() {
		disconnectCalled <- struct{}{}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	select {
	case <-connectCalled:
		// Good
	case <-time.After(1 * time.Second):
		t.Error("OnConnect not called")
	}

	client.Close()

	select {
	case <-disconnectCalled:
		// Good
	case <-time.After(1 * time.Second):
		t.Error("OnDisconnect not called")
	}
}

func TestClient_Unsubscribe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}

			method, _ := msg["method"].(string)
			reqID, _ := msg["req_id"].(float64)

			resp := map[string]interface{}{
				"method":  method,
				"req_id":  int64(reqID),
				"success": true,
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	// Subscribe first
	err := client.SubscribeTicker(ctx, []string{"BTC/USD"})
	if err != nil {
		t.Fatalf("SubscribeTicker failed: %v", err)
	}

	// Then unsubscribe
	err = client.UnsubscribeTicker(ctx, []string{"BTC/USD"})
	if err != nil {
		t.Fatalf("UnsubscribeTicker failed: %v", err)
	}
}

func TestClient_SubscribeBook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var req SubscribeRequest
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			resp := map[string]interface{}{
				"method":  "subscribe",
				"req_id":  req.ReqID,
				"success": true,
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	err := client.SubscribeBook(ctx, []string{"BTC/USD"}, 10)
	if err != nil {
		t.Fatalf("SubscribeBook failed: %v", err)
	}
}

func TestClient_SubscribeOHLC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var req SubscribeRequest
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			resp := map[string]interface{}{
				"method":  "subscribe",
				"req_id":  req.ReqID,
				"success": true,
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	err := client.SubscribeOHLC(ctx, []string{"BTC/USD"}, 60)
	if err != nil {
		t.Fatalf("SubscribeOHLC failed: %v", err)
	}
}

func TestClient_SubscribeTrade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var req SubscribeRequest
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			resp := map[string]interface{}{
				"method":  "subscribe",
				"req_id":  req.ReqID,
				"success": true,
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	err := client.SubscribeTrade(ctx, []string{"BTC/USD"})
	if err != nil {
		t.Fatalf("SubscribeTrade failed: %v", err)
	}
}

func TestClient_SetToken(t *testing.T) {
	client := New("wss://test.example.com")
	client.SetToken("new-token")

	if client.token != "new-token" {
		t.Errorf("expected token new-token, got %s", client.token)
	}
}

func TestClient_NotConnected(t *testing.T) {
	client := New("wss://test.example.com")

	// Try to send message without connection
	err := client.sendMessage(map[string]string{"test": "data"})
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestClient_DoubleClose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// First close
	err := client.Close()
	if err != nil {
		t.Fatalf("First close failed: %v", err)
	}

	// Second close should be safe
	err = client.Close()
	if err != nil {
		t.Fatalf("Second close failed: %v", err)
	}
}

func TestClient_AlreadyConnected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := New(wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First connect
	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("First connect failed: %v", err)
	}
	defer client.Close()

	// Second connect should return nil (already connected)
	err = client.Connect(ctx)
	if err != nil {
		t.Errorf("Second connect should succeed (already connected), got: %v", err)
	}
}

func TestValidateBookDepth(t *testing.T) {
	validDepths := []int{10, 25, 100, 500, 1000}
	invalidDepths := []int{0, 5, 50, 200, 2000}

	for _, depth := range validDepths {
		if err := ValidateBookDepth(depth); err != nil {
			t.Errorf("depth %d should be valid, got error: %v", depth, err)
		}
	}

	for _, depth := range invalidDepths {
		if err := ValidateBookDepth(depth); err == nil {
			t.Errorf("depth %d should be invalid", depth)
		}
	}
}

func TestValidateOHLCInterval(t *testing.T) {
	validIntervals := []int{1, 5, 15, 30, 60, 240, 1440, 10080, 21600}
	invalidIntervals := []int{0, 2, 10, 100, 500}

	for _, interval := range validIntervals {
		if err := ValidateOHLCInterval(interval); err != nil {
			t.Errorf("interval %d should be valid, got error: %v", interval, err)
		}
	}

	for _, interval := range invalidIntervals {
		if err := ValidateOHLCInterval(interval); err == nil {
			t.Errorf("interval %d should be invalid", interval)
		}
	}
}

func TestValidateSymbols(t *testing.T) {
	// Valid cases
	if err := ValidateSymbols([]string{"BTC/USD"}); err != nil {
		t.Errorf("single symbol should be valid: %v", err)
	}
	if err := ValidateSymbols([]string{"BTC/USD", "ETH/USD"}); err != nil {
		t.Errorf("multiple symbols should be valid: %v", err)
	}

	// Invalid cases
	if err := ValidateSymbols(nil); err == nil {
		t.Error("nil symbols should be invalid")
	}
	if err := ValidateSymbols([]string{}); err == nil {
		t.Error("empty symbols should be invalid")
	}
	if err := ValidateSymbols([]string{""}); err == nil {
		t.Error("empty string symbol should be invalid")
	}
	if err := ValidateSymbols([]string{"BTC/USD", ""}); err == nil {
		t.Error("symbols with empty string should be invalid")
	}
}

func TestSubscribeTickerValidation(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()

	// Empty symbols should fail validation
	err := client.SubscribeTicker(ctx, []string{})
	if err == nil {
		t.Error("expected error for empty symbols")
	}
}

func TestSubscribeBookValidation(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()

	// Invalid depth should fail validation
	err := client.SubscribeBook(ctx, []string{"BTC/USD"}, 50)
	if err == nil {
		t.Error("expected error for invalid depth")
	}

	// Empty symbols should fail validation
	err = client.SubscribeBook(ctx, []string{}, 10)
	if err == nil {
		t.Error("expected error for empty symbols")
	}
}

func TestSubscribeOHLCValidation(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()

	// Invalid interval should fail validation
	err := client.SubscribeOHLC(ctx, []string{"BTC/USD"}, 10)
	if err == nil {
		t.Error("expected error for invalid interval")
	}

	// Empty symbols should fail validation
	err = client.SubscribeOHLC(ctx, []string{}, 60)
	if err == nil {
		t.Error("expected error for empty symbols")
	}
}

func TestSubscribeTradeValidation(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()

	// Empty symbols should fail validation
	err := client.SubscribeTrade(ctx, []string{})
	if err == nil {
		t.Error("expected error for empty symbols")
	}
}
