package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_AddOrder_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.AddOrder(ctx, AddOrderParams{})
	if err == nil {
		t.Error("expected error for no token")
	}
	if !strings.Contains(err.Error(), "authentication token required") {
		t.Errorf("expected token error, got: %v", err)
	}
}

func TestClient_EditOrder_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.EditOrder(ctx, EditOrderParams{})
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_AmendOrder_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.AmendOrder(ctx, AmendOrderParams{})
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_CancelOrder_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.CancelOrder(ctx, CancelOrderParams{})
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_CancelAll_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.CancelAll(ctx)
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_BatchAdd_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.BatchAdd(ctx, BatchAddParams{})
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_BatchCancel_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.BatchCancel(ctx, BatchCancelParams{})
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_CancelAllOrdersAfter_NoToken(t *testing.T) {
	client := New("wss://test.example.com")

	ctx := context.Background()
	_, err := client.CancelAllOrdersAfter(ctx, 60)
	if err == nil {
		t.Error("expected error for no token")
	}
}

func TestClient_AddOrder_Success(t *testing.T) {
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

			var req struct {
				Method string `json:"method"`
				ReqID  int64  `json:"req_id"`
				Params struct {
					OrderType  string  `json:"order_type"`
					Side       string  `json:"side"`
					Symbol     string  `json:"symbol"`
					OrderQty   float64 `json:"order_qty"`
					LimitPrice float64 `json:"limit_price"`
					Token      string  `json:"token"`
					Validate   bool    `json:"validate,omitempty"`
				} `json:"params"`
			}
			if err := json.Unmarshal(data, &req); err != nil {
				t.Errorf("failed to parse request: %v", err)
				return
			}

			if req.Method != "add_order" {
				t.Errorf("method = %s, want add_order", req.Method)
			}
			if req.Params.Symbol != "BTC/USD" {
				t.Errorf("symbol = %s, want BTC/USD", req.Params.Symbol)
			}
			if req.Params.Side != "buy" {
				t.Errorf("side = %s, want buy", req.Params.Side)
			}
			if req.Params.OrderType != "limit" {
				t.Errorf("order_type = %s, want limit", req.Params.OrderType)
			}
			if req.Params.Token != "test-token" {
				t.Errorf("token = %s, want test-token", req.Params.Token)
			}
			if req.Params.OrderQty < 0.000999 || req.Params.OrderQty > 0.001001 {
				t.Errorf("order_qty = %f, want ~0.001", req.Params.OrderQty)
			}
			if req.Params.LimitPrice != 50000 {
				t.Errorf("limit_price = %f, want 50000", req.Params.LimitPrice)
			}

			reqID := req.ReqID

			resp := map[string]interface{}{
				"method":  "add_order",
				"req_id":  reqID,
				"success": true,
				"result": map[string]interface{}{
					"order_id":      "OXXXXX-XXXXX-XXXXXX",
					"order_userref": 0,
				},
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := NewAuthenticated(wsURL, "test-token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	result, err := client.AddOrder(ctx, AddOrderParams{
		Symbol:     "BTC/USD",
		Side:       "buy",
		OrderType:  "limit",
		OrderQty:   "0.001",
		LimitPrice: "50000",
		Validate:   true,
	})

	if err != nil {
		t.Fatalf("AddOrder failed: %v", err)
	}

	if result.OrderID != "OXXXXX-XXXXX-XXXXXX" {
		t.Errorf("expected order ID OXXXXX-XXXXX-XXXXXX, got %s", result.OrderID)
	}
}

func TestClient_CancelOrder_Success(t *testing.T) {
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

			var req map[string]interface{}
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			reqID := int64(req["req_id"].(float64))

			resp := map[string]interface{}{
				"method":  "cancel_order",
				"req_id":  reqID,
				"success": true,
				"result": []map[string]interface{}{
					{
						"order_id": "OXXXXX-XXXXX-XXXXXX",
					},
				},
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := NewAuthenticated(wsURL, "test-token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	results, err := client.CancelOrder(ctx, CancelOrderParams{
		OrderID: []string{"OXXXXX-XXXXX-XXXXXX"},
	})

	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestClient_CancelAll_Success(t *testing.T) {
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

			var req map[string]interface{}
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			reqID := int64(req["req_id"].(float64))

			resp := map[string]interface{}{
				"method":  "cancel_all",
				"req_id":  reqID,
				"success": true,
				"result": map[string]interface{}{
					"count": 5,
				},
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := NewAuthenticated(wsURL, "test-token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	result, err := client.CancelAll(ctx)
	if err != nil {
		t.Fatalf("CancelAll failed: %v", err)
	}

	if result.Count != 5 {
		t.Errorf("expected count 5, got %d", result.Count)
	}
}

func TestClient_AddOrder_Error(t *testing.T) {
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

			var req map[string]interface{}
			if err := json.Unmarshal(data, &req); err != nil {
				continue
			}

			reqID := int64(req["req_id"].(float64))

			resp := map[string]interface{}{
				"method":  "add_order",
				"req_id":  reqID,
				"success": false,
				"error":   "EOrder:Invalid order",
			}
			conn.WriteJSON(resp)
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client := NewAuthenticated(wsURL, "test-token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	_, err := client.AddOrder(ctx, AddOrderParams{
		Symbol:    "INVALID",
		Side:      "buy",
		OrderType: "limit",
	})

	if err == nil {
		t.Error("expected error for invalid order")
	}
}

func TestClient_InstrumentHandler(t *testing.T) {
	receivedInstrument := make(chan []InstrumentData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "instrument",
			"type":    "snapshot",
			"data": []map[string]interface{}{
				{
					"symbol":          "BTC/USD",
					"base":            "BTC",
					"quote":           "USD",
					"status":          "online",
					"price_precision": 1,
					"qty_precision":   8,
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

	client.OnInstrument(func(instruments []InstrumentData) {
		receivedInstrument <- instruments
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case instruments := <-receivedInstrument:
		if len(instruments) != 1 {
			t.Errorf("expected 1 instrument, got %d", len(instruments))
		}
		if instruments[0].Symbol != "BTC/USD" {
			t.Errorf("expected symbol BTC/USD, got %s", instruments[0].Symbol)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for instrument data")
	}
}

func TestClient_HeartbeatHandler(t *testing.T) {
	receivedHeartbeat := make(chan *HeartbeatData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "heartbeat",
			"type":    "update",
			"data":    map[string]interface{}{},
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

	client.OnHeartbeat(func(heartbeat *HeartbeatData) {
		receivedHeartbeat <- heartbeat
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case <-receivedHeartbeat:
		// Success
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for heartbeat data")
	}
}

func TestClient_ExecutionHandler(t *testing.T) {
	receivedExecution := make(chan []ExecutionData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "executions",
			"type":    "update",
			"data": []map[string]interface{}{
				{
					"order_id":     "OXXXXX-XXXXX-XXXXXX",
					"order_status": "filled",
					"exec_type":    "filled",
					"symbol":       "BTC/USD",
					"side":         "buy",
					"order_qty":    "0.001",
					"order_type":   "limit",
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

	client.OnExecution(func(executions []ExecutionData) {
		receivedExecution <- executions
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case executions := <-receivedExecution:
		if len(executions) != 1 {
			t.Errorf("expected 1 execution, got %d", len(executions))
		}
		if executions[0].OrderID != "OXXXXX-XXXXX-XXXXXX" {
			t.Errorf("expected order ID OXXXXX-XXXXX-XXXXXX, got %s", executions[0].OrderID)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for execution data")
	}
}

func TestClient_BalanceHandler(t *testing.T) {
	receivedBalance := make(chan []BalanceData, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		msg := map[string]interface{}{
			"channel": "balances",
			"type":    "snapshot",
			"data": []map[string]interface{}{
				{
					"asset":   "BTC",
					"balance": "1.5",
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

	client.OnBalance(func(balances []BalanceData) {
		receivedBalance <- balances
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	select {
	case balances := <-receivedBalance:
		if len(balances) != 1 {
			t.Errorf("expected 1 balance, got %d", len(balances))
		}
		if balances[0].Asset != "BTC" {
			t.Errorf("expected asset BTC, got %s", balances[0].Asset)
		}
	case <-time.After(3 * time.Second):
		t.Error("timeout waiting for balance data")
	}
}
