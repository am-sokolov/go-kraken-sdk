package kraken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/am-sokolov/go-kraken-sdk/types"
)

func TestAccountService_GetBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Balance" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]string{
				"XXBT": "1.5000000000",
				"ZUSD": "25000.0000",
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

	result, err := client.Account.GetBalance(context.Background())
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}

	if result["XXBT"] != "1.5000000000" {
		t.Errorf("XXBT = %s, want 1.5000000000", result["XXBT"])
	}
	if result["ZUSD"] != "25000.0000" {
		t.Errorf("ZUSD = %s, want 25000.0000", result["ZUSD"])
	}
}

func TestAccountService_GetTradeBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/TradeBalance" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "ZUSD" {
			t.Errorf("asset = %s, want ZUSD", r.Form.Get("asset"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"eb": "25000.0000",
				"tb": "25000.0000",
				"m":  "0.0000",
				"n":  "0.0000",
				"c":  "0.0000",
				"v":  "0.0000",
				"e":  "25000.0000",
				"mf": "25000.0000",
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

	result, err := client.Account.GetTradeBalance(context.Background(), &GetTradeBalanceOptions{
		Asset: "ZUSD",
	})
	if err != nil {
		t.Fatalf("GetTradeBalance failed: %v", err)
	}

	if result.EquivalentBalance.String() != "25000" {
		t.Errorf("EquivalentBalance = %s, want 25000", result.EquivalentBalance.String())
	}
}

func TestAccountService_GetOpenOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/OpenOrders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"open": map[string]interface{}{
					"OQCLML-BW3P3-BUCMWZ": map[string]interface{}{
						"order_status": "open",
						"descr": map[string]interface{}{
							"pair":      "XBTUSD",
							"type":      "buy",
							"ordertype": "limit",
							"price":     "50000.0",
							"order":     "buy 1.0 XBTUSD @ limit 50000.0",
						},
						"vol":      "1.0",
						"vol_exec": "0.0",
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

	result, err := client.Account.GetOpenOrders(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetOpenOrders failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}

	order, ok := result["OQCLML-BW3P3-BUCMWZ"]
	if !ok {
		t.Fatal("order OQCLML-BW3P3-BUCMWZ not found")
	}
	if order.Status != "open" {
		t.Errorf("Status = %s, want open", order.Status)
	}
}

func TestAccountService_GetClosedOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/ClosedOrders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"closed": map[string]interface{}{
					"OLD-ORDER-123": map[string]interface{}{
						"order_status": "closed",
						"vol":          "1.0",
						"vol_exec":     "1.0",
					},
				},
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

	result, err := client.Account.GetClosedOrders(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetClosedOrders failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}

	order, ok := result["OLD-ORDER-123"]
	if !ok {
		t.Fatal("order OLD-ORDER-123 not found")
	}
	if order.Status != "closed" {
		t.Errorf("Status = %s, want closed", order.Status)
	}
}

func TestAccountService_QueryOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/QueryOrders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		txids := r.Form["txid"]
		if len(txids) != 2 {
			t.Errorf("expected 2 txids, got %d", len(txids))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"OQCLML-BW3P3-BUCMWZ": map[string]interface{}{
					"status": "open",
					"vol":    "1.0",
				},
				"OQCLML-BW3P3-ABCDEF": map[string]interface{}{
					"status": "closed",
					"vol":    "2.0",
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

	result, err := client.Account.QueryOrders(context.Background(), []string{"OQCLML-BW3P3-BUCMWZ", "OQCLML-BW3P3-ABCDEF"}, false)
	if err != nil {
		t.Fatalf("QueryOrders failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}
}

func TestAccountService_GetTradesHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/TradesHistory" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"trades": map[string]interface{}{
					"TXID-123": map[string]interface{}{
						"pair":  "XBTUSD",
						"type":  "buy",
						"price": "50000.0",
						"vol":   "1.0",
						"time":  1616663618.1234,
					},
				},
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

	result, err := client.Account.GetTradesHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetTradesHistory failed: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
}

func TestAccountService_GetLedgers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Ledgers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XXBT" {
			t.Errorf("asset = %s, want XXBT", r.Form.Get("asset"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
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

	result, err := client.Account.GetLedgers(context.Background(), &GetLedgersOptions{
		Asset: "XXBT",
	})
	if err != nil {
		t.Fatalf("GetLedgers failed: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
}

func TestAccountService_GetTradeVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/TradeVolume" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"currency": "ZUSD",
				"volume":   "10000.0000",
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

	result, err := client.Account.GetTradeVolume(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetTradeVolume failed: %v", err)
	}

	if result.Currency != "ZUSD" {
		t.Errorf("Currency = %s, want ZUSD", result.Currency)
	}
	if result.Volume.String() != "10000" {
		t.Errorf("Volume = %s, want 10000", result.Volume.String())
	}
}

func TestAccountService_GetWebSocketsToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/GetWebSocketsToken" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"token":   "1234567890abcdef",
				"expires": 900,
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

	result, err := client.Account.GetWebSocketsToken(context.Background())
	if err != nil {
		t.Fatalf("GetWebSocketsToken failed: %v", err)
	}

	if result.Token != "1234567890abcdef" {
		t.Errorf("Token = %s, want 1234567890abcdef", result.Token)
	}
	if result.Expires != 900 {
		t.Errorf("Expires = %d, want 900", result.Expires)
	}
}

func TestAccountService_GetL3OrderBook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Level3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"pair": "YFI/EUR",
				"bids": []map[string]interface{}{
					{"price": "3062.00000", "qty": "0.29665800", "order_id": "O5KJU4-IEQTM-NDMS6W", "timestamp": 1765622008594292000},
					{"price": "3062.00000", "qty": "0.13917400", "order_id": "OERRY6-MXYER-6EQKNY", "timestamp": 1765622011396903000},
				},
				"asks": []map[string]interface{}{
					{"price": "3066.00000", "qty": "0.00278335", "order_id": "ORAWGV-N5L4J-LBA3WH", "timestamp": 1765622008499456000},
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

	result, err := client.Account.GetL3OrderBook(context.Background(), "YFI/EUR", &GetL3OrderBookOptions{Depth: 10})
	if err != nil {
		t.Fatalf("GetL3OrderBook failed: %v", err)
	}

	if result.Pair != "YFI/EUR" {
		t.Errorf("Pair = %s, want YFI/EUR", result.Pair)
	}
	if len(result.Bids) != 2 {
		t.Errorf("len(Bids) = %d, want 2", len(result.Bids))
	}
	if len(result.Asks) != 1 {
		t.Errorf("len(Asks) = %d, want 1", len(result.Asks))
	}
	if result.Bids[0].OrderID != "O5KJU4-IEQTM-NDMS6W" {
		t.Errorf("Bids[0].OrderID = %s, want O5KJU4-IEQTM-NDMS6W", result.Bids[0].OrderID)
	}
}

func TestAccountService_GetExtendedBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/BalanceEx" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBT": map[string]interface{}{
					"balance":     "3.46840030",
					"credit":      "1.26844502",
					"credit_used": "0.10002300",
					"hold_trade":  "2.14560458",
				},
				"ZUSD": map[string]interface{}{
					"balance":    "50000.00",
					"hold_trade": "1000.00",
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

	result, err := client.Account.GetExtendedBalance(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetExtendedBalance failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}
	if result["XXBT"].Balance != "3.46840030" {
		t.Errorf("XXBT.Balance = %s, want 3.46840030", result["XXBT"].Balance)
	}
	if result["XXBT"].HoldTrade != "2.14560458" {
		t.Errorf("XXBT.HoldTrade = %s, want 2.14560458", result["XXBT"].HoldTrade)
	}
}

func TestAccountService_GetCreditLines(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/CreditLines" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"asset_details": map[string]interface{}{
					"ZUSD": map[string]interface{}{
						"balance":          "1000.5000",
						"credit_limit":     "50000.0000",
						"credit_used":      "12500.0000",
						"available_credit": "37500.0000",
					},
				},
				"limits_monitor": map[string]interface{}{
					"total_credit_usd":           "100000.0000",
					"total_credit_used_usd":      "25000.0000",
					"total_collateral_value_usd": "150000.0000",
					"equity_usd":                 "125000.0000",
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

	result, err := client.Account.GetCreditLines(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetCreditLines failed: %v", err)
	}

	if result.AssetDetails["ZUSD"].CreditLimit != "50000.0000" {
		t.Errorf("ZUSD.CreditLimit = %s, want 50000.0000", result.AssetDetails["ZUSD"].CreditLimit)
	}
	if result.LimitsMonitor.TotalCreditUSD != "100000.0000" {
		t.Errorf("TotalCreditUSD = %s, want 100000.0000", result.LimitsMonitor.TotalCreditUSD)
	}
}

func TestAccountService_RequestExportReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/AddExport" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("report") != "trades" {
			t.Errorf("report = %s, want trades", r.Form.Get("report"))
		}
		if r.Form.Get("description") != "yearly report" {
			t.Errorf("description = %s, want yearly report", r.Form.Get("description"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"id": "TCJA",
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

	result, err := client.Account.RequestExportReport(context.Background(), &AddExportRequest{
		Report:      types.ExportReportTrades,
		Description: "yearly report",
	})
	if err != nil {
		t.Fatalf("RequestExportReport failed: %v", err)
	}

	if result.ID != "TCJA" {
		t.Errorf("ID = %s, want TCJA", result.ID)
	}
}

func TestAccountService_GetExportReportStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/ExportStatus" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []map[string]interface{}{
				{
					"id":          "TCJA",
					"descr":       "yearly report",
					"format":      "CSV",
					"report":      "trades",
					"status":      "Processed",
					"createdtm":   "1695728276",
					"completedtm": "1695728280",
					"datastarttm": "1695728276",
					"dataendtm":   "1695828276",
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

	result, err := client.Account.GetExportReportStatus(context.Background(), types.ExportReportTrades)
	if err != nil {
		t.Fatalf("GetExportReportStatus failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("len(result) = %d, want 1", len(result))
	}
	if result[0].ID != "TCJA" {
		t.Errorf("ID = %s, want TCJA", result[0].ID)
	}
	if result[0].Status != types.ExportStatusProcessed {
		t.Errorf("Status = %s, want Processed", result[0].Status)
	}
}

func TestAccountService_RetrieveExportReport(t *testing.T) {
	expectedData := []byte("PK\x03\x04\x14\x00test zip data")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/RetrieveExport" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(expectedData)
	}))
	defer server.Close()

	client, err := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key", "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, err := client.Account.RetrieveExportReport(context.Background(), "TCJA")
	if err != nil {
		t.Fatalf("RetrieveExportReport failed: %v", err)
	}

	if string(result) != string(expectedData) {
		t.Errorf("result = %v, want %v", result, expectedData)
	}
}

func TestAccountService_RemoveExportReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/RemoveExport" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("id") != "TCJA" {
			t.Errorf("id = %s, want TCJA", r.Form.Get("id"))
		}
		if r.Form.Get("type") != "delete" {
			t.Errorf("type = %s, want delete", r.Form.Get("type"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"delete": true,
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

	result, err := client.Account.RemoveExportReport(context.Background(), "TCJA", RemoveExportDelete)
	if err != nil {
		t.Fatalf("RemoveExportReport failed: %v", err)
	}

	if !result.Delete {
		t.Errorf("Delete = %v, want true", result.Delete)
	}
}
