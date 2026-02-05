package kraken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

// MockKrakenServer provides a comprehensive mock Kraken API server for integration testing.
type MockKrakenServer struct {
	server     *httptest.Server
	mu         sync.RWMutex
	orders     map[string]map[string]interface{}
	balances   map[string]string
	positions  map[string]map[string]interface{}
	ledgers    map[string]map[string]interface{}
	orderIDSeq int
}

// NewMockKrakenServer creates a new mock Kraken server.
func NewMockKrakenServer(t *testing.T) *MockKrakenServer {
	m := &MockKrakenServer{
		orders:     make(map[string]map[string]interface{}),
		balances:   map[string]string{"XXBT": "1.5", "ZUSD": "50000.00"},
		positions:  make(map[string]map[string]interface{}),
		ledgers:    make(map[string]map[string]interface{}),
		orderIDSeq: 1000,
	}

	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("/0/public/Time", m.handleTime)
	mux.HandleFunc("/0/public/SystemStatus", m.handleSystemStatus)
	mux.HandleFunc("/0/public/Assets", m.handleAssets)
	mux.HandleFunc("/0/public/AssetPairs", m.handleAssetPairs)
	mux.HandleFunc("/0/public/Ticker", m.handleTicker)
	mux.HandleFunc("/0/public/Depth", m.handleOrderBook)
	mux.HandleFunc("/0/public/OHLC", m.handleOHLC)
	mux.HandleFunc("/0/public/Trades", m.handleRecentTrades)

	// Private account endpoints
	mux.HandleFunc("/0/private/Balance", m.handleBalance)
	mux.HandleFunc("/0/private/TradeBalance", m.handleTradeBalance)
	mux.HandleFunc("/0/private/OpenOrders", m.handleOpenOrders)
	mux.HandleFunc("/0/private/ClosedOrders", m.handleClosedOrders)
	mux.HandleFunc("/0/private/QueryOrders", m.handleQueryOrders)
	mux.HandleFunc("/0/private/TradesHistory", m.handleTradesHistory)
	mux.HandleFunc("/0/private/Ledgers", m.handleLedgers)
	mux.HandleFunc("/0/private/TradeVolume", m.handleTradeVolume)
	mux.HandleFunc("/0/private/GetWebSocketsToken", m.handleWebSocketsToken)
	mux.HandleFunc("/0/private/OpenPositions", m.handleOpenPositions)

	// Private trading endpoints
	mux.HandleFunc("/0/private/AddOrder", m.handleAddOrder)
	mux.HandleFunc("/0/private/EditOrder", m.handleEditOrder)
	mux.HandleFunc("/0/private/CancelOrder", m.handleCancelOrder)
	mux.HandleFunc("/0/private/CancelAll", m.handleCancelAll)

	m.server = httptest.NewServer(mux)
	return m
}

// URL returns the mock server URL.
func (m *MockKrakenServer) URL() string {
	return m.server.URL
}

// Close shuts down the mock server.
func (m *MockKrakenServer) Close() {
	m.server.Close()
}

// SetBalance sets a balance for testing.
func (m *MockKrakenServer) SetBalance(asset, amount string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.balances[asset] = amount
}

// AddOrder adds an order for testing.
func (m *MockKrakenServer) AddOrder(id string, order map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders[id] = order
}

func (m *MockKrakenServer) respond(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  []string{},
		"result": result,
	})
}

func (m *MockKrakenServer) respondError(w http.ResponseWriter, errors ...string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  errors,
		"result": nil,
	})
}

func (m *MockKrakenServer) handleTime(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"unixtime": time.Now().Unix(),
		"rfc1123":  time.Now().UTC().Format(time.RFC1123),
	})
}

func (m *MockKrakenServer) handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"status":    "online",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (m *MockKrakenServer) handleAssets(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"XXBT": map[string]interface{}{
			"aclass":           "currency",
			"altname":          "XBT",
			"decimals":         10,
			"display_decimals": 5,
		},
		"ZUSD": map[string]interface{}{
			"aclass":           "currency",
			"altname":          "USD",
			"decimals":         4,
			"display_decimals": 2,
		},
		"XETH": map[string]interface{}{
			"aclass":           "currency",
			"altname":          "ETH",
			"decimals":         10,
			"display_decimals": 5,
		},
	})
}

func (m *MockKrakenServer) handleAssetPairs(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"XXBTZUSD": map[string]interface{}{
			"altname":       "XBTUSD",
			"wsname":        "XBT/USD",
			"aclass_base":   "currency",
			"base":          "XXBT",
			"aclass_quote":  "currency",
			"quote":         "ZUSD",
			"pair_decimals": 1,
			"lot_decimals":  8,
			"ordermin":      "0.0001",
		},
		"XETHZUSD": map[string]interface{}{
			"altname":       "ETHUSD",
			"wsname":        "ETH/USD",
			"aclass_base":   "currency",
			"base":          "XETH",
			"aclass_quote":  "currency",
			"quote":         "ZUSD",
			"pair_decimals": 2,
			"lot_decimals":  8,
			"ordermin":      "0.01",
		},
	})
}

func (m *MockKrakenServer) handleTicker(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"XXBTZUSD": map[string]interface{}{
			"a": []string{"50000.0", "1", "1.000"},
			"b": []string{"49999.0", "1", "1.000"},
			"c": []string{"50000.0", "0.1"},
			"v": []string{"1000.0", "5000.0"},
			"p": []string{"49500.0", "49000.0"},
			"t": []int{1000, 5000},
			"l": []string{"48000.0", "47000.0"},
			"h": []string{"51000.0", "52000.0"},
			"o": "49000.0",
		},
	})
}

func (m *MockKrakenServer) handleOrderBook(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"XXBTZUSD": map[string]interface{}{
			"asks": [][]interface{}{
				{"50000.0", "1.0", 1616663618},
				{"50001.0", "2.0", 1616663618},
			},
			"bids": [][]interface{}{
				{"49999.0", "1.5", 1616663618},
				{"49998.0", "3.0", 1616663618},
			},
		},
	})
}

func (m *MockKrakenServer) handleOHLC(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"XXBTZUSD": [][]interface{}{
			{1616663580.0, "50000.0", "50100.0", "49900.0", "50050.0", "50025.0", "10.5", 100},
		},
		"last": 1616663580,
	})
}

func (m *MockKrakenServer) handleRecentTrades(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"XXBTZUSD": [][]interface{}{
			{"50000.0", "0.1", 1616663618.1234, "b", "m", "", 12345678.0},
		},
		"last": "1616663618123456789",
	})
}

func (m *MockKrakenServer) handleBalance(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	balances := make(map[string]string)
	for k, v := range m.balances {
		balances[k] = v
	}
	m.mu.RUnlock()

	m.respond(w, balances)
}

func (m *MockKrakenServer) handleTradeBalance(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"eb": "50000.0",
		"tb": "50000.0",
		"m":  "0.0",
		"n":  "0.0",
		"c":  "0.0",
		"v":  "0.0",
		"e":  "50000.0",
		"mf": "50000.0",
	})
}

func (m *MockKrakenServer) handleOpenOrders(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	openOrders := make(map[string]interface{})
	for id, order := range m.orders {
		if status, ok := order["status"].(string); ok && status == "open" {
			openOrders[id] = order
		}
	}
	m.mu.RUnlock()

	m.respond(w, map[string]interface{}{
		"open": openOrders,
	})
}

func (m *MockKrakenServer) handleClosedOrders(w http.ResponseWriter, r *http.Request) {
	// Check for pagination params
	if err := r.ParseForm(); err == nil {
		// Validate that pagination params are properly formatted
		if start := r.Form.Get("start"); start != "" {
			// Should be a valid timestamp string
			if !isNumeric(start) {
				m.respondError(w, "EGeneral:Invalid arguments")
				return
			}
		}
	}

	m.mu.RLock()
	closedOrders := make(map[string]interface{})
	for id, order := range m.orders {
		if status, ok := order["status"].(string); ok && status == "closed" {
			closedOrders[id] = order
		}
	}
	m.mu.RUnlock()

	m.respond(w, map[string]interface{}{
		"closed": closedOrders,
		"count":  len(closedOrders),
	})
}

func (m *MockKrakenServer) handleQueryOrders(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		m.respondError(w, "EGeneral:Invalid arguments")
		return
	}

	txids := r.Form["txid"]
	m.mu.RLock()
	result := make(map[string]interface{})
	for _, id := range txids {
		if order, ok := m.orders[id]; ok {
			result[id] = order
		}
	}
	m.mu.RUnlock()

	m.respond(w, result)
}

func (m *MockKrakenServer) handleTradesHistory(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"trades": map[string]interface{}{
			"TXID-123": map[string]interface{}{
				"pair":  "XBTUSD",
				"type":  "buy",
				"price": "50000.0",
				"vol":   "0.1",
				"time":  1616663618.1234,
			},
		},
		"count": 1,
	})
}

func (m *MockKrakenServer) handleLedgers(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"ledger": map[string]interface{}{
			"L123": map[string]interface{}{
				"refid":   "REFID-123",
				"time":    1616663618.1234,
				"type":    "deposit",
				"aclass":  "currency",
				"asset":   "XXBT",
				"amount":  "1.0",
				"fee":     "0.0",
				"balance": "1.5",
			},
		},
		"count": 1,
	})
}

func (m *MockKrakenServer) handleTradeVolume(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"currency": "ZUSD",
		"volume":   "10000.0",
	})
}

func (m *MockKrakenServer) handleWebSocketsToken(w http.ResponseWriter, r *http.Request) {
	m.respond(w, map[string]interface{}{
		"token":   "mock-ws-token-12345",
		"expires": 900,
	})
}

func (m *MockKrakenServer) handleOpenPositions(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	positions := make(map[string]interface{})
	for id, pos := range m.positions {
		positions[id] = pos
	}
	m.mu.RUnlock()

	m.respond(w, positions)
}

func (m *MockKrakenServer) handleAddOrder(w http.ResponseWriter, r *http.Request) {
	var (
		pair      string
		orderType string
		orderKind string
		volume    string
		price     string
		validate  bool
	)

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var body struct {
			Pair      string `json:"pair"`
			Type      string `json:"type"`
			OrderType string `json:"ordertype"`
			Volume    string `json:"volume"`
			Price     string `json:"price"`
			Validate  bool   `json:"validate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			m.respondError(w, "EGeneral:Invalid arguments")
			return
		}
		pair = body.Pair
		orderType = body.Type
		orderKind = body.OrderType
		volume = body.Volume
		price = body.Price
		validate = body.Validate
	} else {
		if err := r.ParseForm(); err != nil {
			m.respondError(w, "EGeneral:Invalid arguments")
			return
		}

		pair = r.Form.Get("pair")
		orderType = r.Form.Get("type")
		orderKind = r.Form.Get("ordertype")
		volume = r.Form.Get("volume")
		price = r.Form.Get("price")
		validate = r.Form.Get("validate") == "true"
	}

	if pair == "" || orderType == "" || volume == "" {
		m.respondError(w, "EGeneral:Invalid arguments")
		return
	}

	// Simulate validation mode
	if validate {
		m.respond(w, map[string]interface{}{
			"descr": map[string]interface{}{
				"order": orderType + " " + volume + " " + pair + " @ " + orderKind + " " + price,
			},
		})
		return
	}

	m.mu.Lock()
	m.orderIDSeq++
	orderID := "O" + strings.Repeat("X", 5) + "-" + strings.Repeat("X", 5) + "-" + string(rune('A'+m.orderIDSeq%26)) + strings.Repeat("X", 5)

	m.orders[orderID] = map[string]interface{}{
		"status": "open",
		"descr": map[string]interface{}{
			"pair":      pair,
			"type":      orderType,
			"ordertype": orderKind,
			"price":     price,
		},
		"vol":      volume,
		"vol_exec": "0.0",
	}
	m.mu.Unlock()

	m.respond(w, map[string]interface{}{
		"descr": map[string]interface{}{
			"order": orderType + " " + volume + " " + pair + " @ " + orderKind + " " + price,
		},
		"txid": []string{orderID},
	})
}

func (m *MockKrakenServer) handleEditOrder(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		m.respondError(w, "EGeneral:Invalid arguments")
		return
	}

	txid := r.Form.Get("txid")
	if txid == "" {
		m.respondError(w, "EGeneral:Invalid arguments")
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.orders[txid]; !exists {
		m.respondError(w, "EOrder:Unknown order")
		return
	}

	// Create new order ID for the edited order
	m.orderIDSeq++
	newOrderID := "O" + strings.Repeat("X", 5) + "-NEW-" + string(rune('A'+m.orderIDSeq%26)) + strings.Repeat("X", 5)

	// Mark old order as canceled
	m.orders[txid]["status"] = "canceled"

	// Create new order
	m.orders[newOrderID] = map[string]interface{}{
		"status": "open",
	}

	m.respond(w, map[string]interface{}{
		"txid":         newOrderID,
		"originaltxid": txid,
		"status":       "ok",
	})
}

func (m *MockKrakenServer) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		m.respondError(w, "EGeneral:Invalid arguments")
		return
	}

	txid := r.Form.Get("txid")
	if txid == "" {
		m.respondError(w, "EGeneral:Invalid arguments")
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if order, exists := m.orders[txid]; exists {
		order["status"] = "canceled"
		m.respond(w, map[string]interface{}{
			"count": 1,
		})
	} else {
		m.respondError(w, "EOrder:Unknown order")
	}
}

func (m *MockKrakenServer) handleCancelAll(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	count := 0
	for _, order := range m.orders {
		if status, ok := order["status"].(string); ok && status == "open" {
			order["status"] = "canceled"
			count++
		}
	}
	m.mu.Unlock()

	m.respond(w, map[string]interface{}{
		"count": count,
	})
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// Integration Tests

func TestIntegration_TradingWorkflow(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Check system status
	status, err := client.Public.GetSystemStatus(ctx)
	if err != nil {
		t.Fatalf("GetSystemStatus failed: %v", err)
	}
	if status.Status != "online" {
		t.Errorf("expected online status, got %s", status.Status)
	}

	// Step 2: Get balance
	balances, err := client.Account.GetBalance(ctx)
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}
	if _, ok := balances["XXBT"]; !ok {
		t.Error("expected XXBT balance")
	}

	// Step 3: Get ticker
	ticker, err := client.Public.GetTicker(ctx, []string{"XBTUSD"})
	if err != nil {
		t.Fatalf("GetTicker failed: %v", err)
	}
	if _, ok := ticker["XXBTZUSD"]; !ok {
		t.Error("expected XXBTZUSD ticker")
	}

	// Step 4: Place order
	orderResult, err := client.Trading.AddOrder(ctx, &AddOrderRequest{
		Pair:      "XBTUSD",
		Side:      types.SideBuy,
		OrderType: types.OrderTypeLimit,
		Volume:    "0.001",
		Price:     "45000.0",
	})
	if err != nil {
		t.Fatalf("AddOrder failed: %v", err)
	}
	if len(orderResult.TransactionIDs) == 0 {
		t.Error("expected order ID")
	}
	orderID := orderResult.TransactionIDs[0]

	// Step 5: Check open orders
	openOrders, err := client.Account.GetOpenOrders(ctx, nil)
	if err != nil {
		t.Fatalf("GetOpenOrders failed: %v", err)
	}
	if _, ok := openOrders[orderID]; !ok {
		t.Error("expected order in open orders")
	}

	// Step 6: Cancel order
	cancelResult, err := client.Trading.CancelOrder(ctx, orderID)
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}
	if cancelResult.Count != 1 {
		t.Errorf("expected 1 canceled order, got %d", cancelResult.Count)
	}
}

func TestIntegration_MarketDataWorkflow(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Get server time
	serverTime, err := client.Public.GetServerTime(ctx)
	if err != nil {
		t.Fatalf("GetServerTime failed: %v", err)
	}
	if serverTime.UnixTime == 0 {
		t.Error("expected non-zero unix time")
	}

	// Step 2: Get assets
	assets, err := client.Public.GetAssets(ctx, nil)
	if err != nil {
		t.Fatalf("GetAssets failed: %v", err)
	}
	if len(assets) == 0 {
		t.Error("expected assets")
	}

	// Step 3: Get asset pairs
	pairs, err := client.Public.GetAssetPairs(ctx, nil)
	if err != nil {
		t.Fatalf("GetAssetPairs failed: %v", err)
	}
	if len(pairs) == 0 {
		t.Error("expected asset pairs")
	}

	// Step 4: Get order book
	orderBook, err := client.Public.GetOrderBook(ctx, "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetOrderBook failed: %v", err)
	}
	if book, ok := orderBook["XXBTZUSD"]; !ok {
		t.Error("expected XXBTZUSD order book")
	} else {
		if len(book.Asks) == 0 {
			t.Error("expected asks in order book")
		}
		if len(book.Bids) == 0 {
			t.Error("expected bids in order book")
		}
	}

	// Step 5: Get OHLC data
	ohlc, err := client.Public.GetOHLC(ctx, "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetOHLC failed: %v", err)
	}
	if ohlc.Last == 0 {
		t.Error("expected OHLC data")
	}

	// Step 6: Get recent trades
	trades, err := client.Public.GetRecentTrades(ctx, "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetRecentTrades failed: %v", err)
	}
	if trades.Last == "" {
		t.Error("expected trades")
	}
}

func TestIntegration_AccountDataWorkflow(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Get balance
	balances, err := client.Account.GetBalance(ctx)
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}
	if balances["ZUSD"] == "" {
		t.Error("expected USD balance")
	}

	// Step 2: Get trade balance
	tradeBalance, err := client.Account.GetTradeBalance(ctx, &GetTradeBalanceOptions{
		Asset: "ZUSD",
	})
	if err != nil {
		t.Fatalf("GetTradeBalance failed: %v", err)
	}
	if tradeBalance.EquivalentBalance.IsZero() {
		t.Error("expected non-zero equivalent balance")
	}

	// Step 3: Get trades history
	tradesHistory, err := client.Account.GetTradesHistory(ctx, &GetTradesHistoryOptions{
		Start:  1600000000,
		End:    1700000000,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("GetTradesHistory failed: %v", err)
	}
	if tradesHistory == nil {
		t.Error("expected trades history")
	}

	// Step 4: Get ledgers
	ledgers, err := client.Account.GetLedgers(ctx, &GetLedgersOptions{
		Asset: "XXBT",
	})
	if err != nil {
		t.Fatalf("GetLedgers failed: %v", err)
	}
	if ledgers == nil {
		t.Error("expected ledgers")
	}

	// Step 5: Get trade volume
	tradeVolume, err := client.Account.GetTradeVolume(ctx, &GetTradeVolumeOptions{
		Pair: "XBTUSD",
	})
	if err != nil {
		t.Fatalf("GetTradeVolume failed: %v", err)
	}
	if tradeVolume.Currency == "" {
		t.Error("expected trade volume currency")
	}

	// Step 6: Get WebSockets token
	wsToken, err := client.Account.GetWebSocketsToken(ctx)
	if err != nil {
		t.Fatalf("GetWebSocketsToken failed: %v", err)
	}
	if wsToken.Token == "" {
		t.Error("expected WebSocket token")
	}
}

func TestIntegration_PaginationParams(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test GetClosedOrders with pagination
	_, err = client.Account.GetClosedOrders(ctx, &GetClosedOrdersOptions{
		Start:  1600000000,
		End:    1700000000,
		Offset: 10,
	})
	if err != nil {
		t.Errorf("GetClosedOrders with pagination failed: %v", err)
	}

	// Test GetTradesHistory with pagination
	_, err = client.Account.GetTradesHistory(ctx, &GetTradesHistoryOptions{
		Start:  1600000000,
		End:    1700000000,
		Offset: 10,
	})
	if err != nil {
		t.Errorf("GetTradesHistory with pagination failed: %v", err)
	}

	// Test GetLedgers with pagination
	_, err = client.Account.GetLedgers(ctx, &GetLedgersOptions{
		Start:  1600000000,
		End:    1700000000,
		Offset: 10,
	})
	if err != nil {
		t.Errorf("GetLedgers with pagination failed: %v", err)
	}
}

func TestIntegration_RateLimiting(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
		WithRateLimitTier(TierIntermediate),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Make several rapid requests - should not fail with rate limiter
	for i := 0; i < 5; i++ {
		_, err := client.Public.GetServerTime(ctx)
		if err != nil {
			t.Errorf("request %d failed: %v", i, err)
		}
	}
}

func TestIntegration_OrderValidation(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test order validation (validate=true)
	result, err := client.Trading.AddOrder(ctx, &AddOrderRequest{
		Pair:      "XBTUSD",
		Side:      types.SideBuy,
		OrderType: types.OrderTypeLimit,
		Volume:    "0.001",
		Price:     "45000.0",
		Validate:  true,
	})
	if err != nil {
		t.Fatalf("AddOrder validation failed: %v", err)
	}

	// Validation should not return an order ID
	if len(result.TransactionIDs) > 0 {
		t.Error("validation mode should not return order ID")
	}

	// Should return order description
	if result.Description.Order == "" {
		t.Error("expected order description")
	}
}

func TestIntegration_CancelAllOrders(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	// Add some orders to the mock server
	mock.AddOrder("ORDER-1", map[string]interface{}{"status": "open"})
	mock.AddOrder("ORDER-2", map[string]interface{}{"status": "open"})
	mock.AddOrder("ORDER-3", map[string]interface{}{"status": "closed"})

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Cancel all open orders
	result, err := client.Trading.CancelAllOrders(ctx)
	if err != nil {
		t.Fatalf("CancelAllOrders failed: %v", err)
	}

	// Should cancel 2 orders (the open ones)
	if result.Count != 2 {
		t.Errorf("expected 2 canceled orders, got %d", result.Count)
	}
}

func TestIntegration_EditOrder(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	// Add an order to edit
	mock.AddOrder("ORIGINAL-ORDER", map[string]interface{}{"status": "open"})

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Edit the order
	result, err := client.Trading.EditOrder(ctx, &EditOrderRequest{
		TxID:   "ORIGINAL-ORDER",
		Pair:   "XBTUSD",
		Volume: "0.002",
		Price:  "46000.0",
	})
	if err != nil {
		t.Fatalf("EditOrder failed: %v", err)
	}

	if result.OriginalTransactionID != "ORIGINAL-ORDER" {
		t.Errorf("expected original txid ORIGINAL-ORDER, got %s", result.OriginalTransactionID)
	}
	if result.NewTransactionID == "" {
		t.Error("expected new txid")
	}
}

func TestIntegration_OptionalParams(t *testing.T) {
	mock := NewMockKrakenServer(t)
	defer mock.Close()

	client, err := New(
		WithBaseURL(mock.URL()),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test GetOpenOrders with optional params
	_, err = client.Account.GetOpenOrders(ctx, &GetOpenOrdersOptions{
		Trades:           true,
		UserRef:          12345,
		ClientOrderID:    "client-order-123",
		RebaseMultiplier: "rebased",
	})
	if err != nil {
		t.Errorf("GetOpenOrders with optional params failed: %v", err)
	}

	// Test GetClosedOrders with optional params
	consolidate := true
	_, err = client.Account.GetClosedOrders(ctx, &GetClosedOrdersOptions{
		Trades:           true,
		UserRef:          12345,
		ClientOrderID:    "client-order-123",
		ConsolidateTaker: &consolidate,
		WithoutCount:     true,
		RebaseMultiplier: "rebased",
	})
	if err != nil {
		t.Errorf("GetClosedOrders with optional params failed: %v", err)
	}

	// Test GetTradesHistory with optional params
	_, err = client.Account.GetTradesHistory(ctx, &GetTradesHistoryOptions{
		Type:             "all",
		Trades:           true,
		WithoutCount:     true,
		ConsolidateTaker: &consolidate,
		Ledgers:          true,
		RebaseMultiplier: "rebased",
	})
	if err != nil {
		t.Errorf("GetTradesHistory with optional params failed: %v", err)
	}
}
