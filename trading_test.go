package kraken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

func TestTradingService_AddOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/AddOrder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", r.Header.Get("Content-Type"))
		}

		var body struct {
			Nonce     int64  `json:"nonce"`
			OrderType string `json:"ordertype"`
			Side      string `json:"type"`
			Volume    string `json:"volume"`
			Pair      string `json:"pair"`
			Price     string `json:"price,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Nonce == 0 {
			t.Errorf("nonce = %d, want non-zero", body.Nonce)
		}
		if body.OrderType != "limit" {
			t.Errorf("ordertype = %s, want limit", body.OrderType)
		}
		if body.Side != "buy" {
			t.Errorf("type = %s, want buy", body.Side)
		}
		if body.Volume != "1.0" {
			t.Errorf("volume = %s, want 1.0", body.Volume)
		}
		if body.Pair != "XBTUSD" {
			t.Errorf("pair = %s, want XBTUSD", body.Pair)
		}
		if body.Price != "50000" {
			t.Errorf("price = %s, want 50000", body.Price)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"descr": map[string]interface{}{
					"order": "buy 1.0 XBTUSD @ limit 50000",
				},
				"txid": []string{"OQCLML-BW3P3-BUCMWZ"},
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.AddOrder(context.Background(), &AddOrderRequest{
		OrderType: types.OrderTypeLimit,
		Side:      types.SideBuy,
		Volume:    "1.0",
		Pair:      "XBTUSD",
		Price:     "50000",
	})
	if err != nil {
		t.Fatalf("AddOrder failed: %v", err)
	}

	if result.Description.Order != "buy 1.0 XBTUSD @ limit 50000" {
		t.Errorf("Description.Order = %s, want buy 1.0 XBTUSD @ limit 50000", result.Description.Order)
	}
	if len(result.TransactionIDs) != 1 || result.TransactionIDs[0] != "OQCLML-BW3P3-BUCMWZ" {
		t.Errorf("TransactionIDs = %v, want [OQCLML-BW3P3-BUCMWZ]", result.TransactionIDs)
	}
}

func TestTradingService_AddOrderWithValidate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Validate bool `json:"validate"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if !body.Validate {
			t.Errorf("validate = %v, want true", body.Validate)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"descr": map[string]interface{}{
					"order": "buy 1.0 XBTUSD @ limit 50000",
				},
				"txid": []string{},
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.AddOrder(context.Background(), &AddOrderRequest{
		OrderType: types.OrderTypeLimit,
		Side:      types.SideBuy,
		Volume:    "1.0",
		Pair:      "XBTUSD",
		Price:     "50000",
		Validate:  true,
	})
	if err != nil {
		t.Fatalf("AddOrder failed: %v", err)
	}

	// Validate mode returns no transaction IDs
	if len(result.TransactionIDs) != 0 {
		t.Errorf("TransactionIDs = %v, want []", result.TransactionIDs)
	}
}

func TestTradingService_EditOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/EditOrder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("txid") != "OQCLML-BW3P3-BUCMWZ" {
			t.Errorf("txid = %s, want OQCLML-BW3P3-BUCMWZ", r.Form.Get("txid"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"descr": map[string]interface{}{
					"order": "buy 1.0 XBTUSD @ limit 51000",
				},
				"txid":             "NEW-TXID-123",
				"originaltxid":     "OQCLML-BW3P3-BUCMWZ",
				"volume":           "1.0",
				"price":            "51000",
				"orders_cancelled": 1,
				"status":           "ok",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.EditOrder(context.Background(), &EditOrderRequest{
		TxID:  "OQCLML-BW3P3-BUCMWZ",
		Pair:  "XBTUSD",
		Price: "51000",
	})
	if err != nil {
		t.Fatalf("EditOrder failed: %v", err)
	}

	if result.NewTransactionID != "NEW-TXID-123" {
		t.Errorf("NewTransactionID = %s, want NEW-TXID-123", result.NewTransactionID)
	}
	if result.OriginalTransactionID != "OQCLML-BW3P3-BUCMWZ" {
		t.Errorf("OriginalTransactionID = %s, want OQCLML-BW3P3-BUCMWZ", result.OriginalTransactionID)
	}
}

func TestTradingService_CancelOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/CancelOrder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("txid") != "OQCLML-BW3P3-BUCMWZ" {
			t.Errorf("txid = %s, want OQCLML-BW3P3-BUCMWZ", r.Form.Get("txid"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"count": 1,
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.CancelOrder(context.Background(), "OQCLML-BW3P3-BUCMWZ")
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
}

func TestTradingService_CancelAllOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/CancelAll" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"count": 5,
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.CancelAllOrders(context.Background())
	if err != nil {
		t.Fatalf("CancelAllOrders failed: %v", err)
	}

	if result.Count != 5 {
		t.Errorf("Count = %d, want 5", result.Count)
	}
}

func TestTradingService_CancelAllOrdersAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/CancelAllOrdersAfter" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("timeout") != "60" {
			t.Errorf("timeout = %s, want 60", r.Form.Get("timeout"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"currentTime": "2021-03-24T17:41:56Z",
				"triggerTime": "2021-03-24T17:42:56Z",
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.CancelAllOrdersAfter(context.Background(), 60)
	if err != nil {
		t.Fatalf("CancelAllOrdersAfter failed: %v", err)
	}

	if result.CurrentTime != "2021-03-24T17:41:56Z" {
		t.Errorf("CurrentTime = %s, want 2021-03-24T17:41:56Z", result.CurrentTime)
	}
	if result.TriggerTime != "2021-03-24T17:42:56Z" {
		t.Errorf("TriggerTime = %s, want 2021-03-24T17:42:56Z", result.TriggerTime)
	}
}

func TestTradingService_AddOrderBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/AddOrderBatch" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("pair") != "XBTUSD" {
			t.Errorf("pair = %s, want XBTUSD", r.Form.Get("pair"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"orders": []interface{}{
					map[string]interface{}{
						"descr": map[string]interface{}{
							"order": "buy 1.0 XBTUSD @ limit 50000",
						},
						"txid": "OQCLML-BW3P3-BUCMWZ",
					},
					map[string]interface{}{
						"descr": map[string]interface{}{
							"order": "sell 1.0 XBTUSD @ limit 55000",
						},
						"txid": "OQCLML-BW3P3-ABCDEF",
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Trading.AddOrderBatch(context.Background(), &AddOrderBatchRequest{
		Pair: "XBTUSD",
		Orders: []AddOrderRequest{
			{OrderType: types.OrderTypeLimit, Side: types.SideBuy, Volume: "1.0", Price: "50000"},
			{OrderType: types.OrderTypeLimit, Side: types.SideSell, Volume: "1.0", Price: "55000"},
		},
	})
	if err != nil {
		t.Fatalf("AddOrderBatch failed: %v", err)
	}

	if len(result.Orders) != 2 {
		t.Errorf("len(Orders) = %d, want 2", len(result.Orders))
	}
}
