package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicService_GetServerTime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/Time" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  []string{},
			"result": map[string]interface{}{"unixtime": 1616663618, "rfc1123": "Thu, 25 Mar 21 09:53:38 +0000"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetServerTime(context.Background())
	if err != nil {
		t.Fatalf("GetServerTime failed: %v", err)
	}

	if result.UnixTime != 1616663618 {
		t.Errorf("UnixTime = %d, want 1616663618", result.UnixTime)
	}
}

func TestPublicService_GetSystemStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/SystemStatus" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  []string{},
			"result": map[string]interface{}{"status": "online", "timestamp": "2021-03-25T09:53:38Z"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus failed: %v", err)
	}

	if result.Status != "online" {
		t.Errorf("Status = %s, want online", result.Status)
	}
}

func TestPublicService_GetAssets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/Assets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Check query params
		if asset := r.URL.Query().Get("asset"); asset != "BTC,ETH" {
			t.Errorf("asset = %s, want BTC,ETH", asset)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBT": map[string]interface{}{"aclass": "currency", "altname": "XBT", "decimals": 10, "display_decimals": 5},
				"XETH": map[string]interface{}{"aclass": "currency", "altname": "ETH", "decimals": 10, "display_decimals": 5},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetAssets(context.Background(), &GetAssetsOptions{
		Assets: []string{"BTC", "ETH"},
	})
	if err != nil {
		t.Fatalf("GetAssets failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}

	if result["XXBT"].Altname != "XBT" {
		t.Errorf("XXBT.Altname = %s, want XBT", result["XXBT"].Altname)
	}
}

func TestPublicService_GetAssetPairs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/AssetPairs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
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
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetAssetPairs(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetAssetPairs failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}

	pair := result["XXBTZUSD"]
	if pair.Altname != "XBTUSD" {
		t.Errorf("Altname = %s, want XBTUSD", pair.Altname)
	}
	if pair.WSName != "XBT/USD" {
		t.Errorf("WSName = %s, want XBT/USD", pair.WSName)
	}
}

func TestPublicService_GetTicker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/Ticker" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if pair := r.URL.Query().Get("pair"); pair != "XBTUSD" {
			t.Errorf("pair = %s, want XBTUSD", pair)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBTZUSD": map[string]interface{}{
					"a": []string{"50000.00000", "1", "1.000"},
					"b": []string{"49999.00000", "1", "1.000"},
					"c": []string{"50000.00000", "0.10000000"},
					"v": []string{"1234.56789", "5678.90123"},
					"p": []string{"49500.00000", "49000.00000"},
					"t": []int{1000, 5000},
					"l": []string{"48000.00000", "47000.00000"},
					"h": []string{"51000.00000", "52000.00000"},
					"o": "48500.00000",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetTicker(context.Background(), []string{"XBTUSD"})
	if err != nil {
		t.Fatalf("GetTicker failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}

	ticker := result["XXBTZUSD"]
	if len(ticker.Ask) != 3 {
		t.Errorf("len(Ask) = %d, want 3", len(ticker.Ask))
	}
}

func TestPublicService_GetOHLC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/OHLC" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBTZUSD": [][]interface{}{
					{1616663580.0, "50000.0", "50100.0", "49900.0", "50050.0", "50025.0", "10.5", 100},
					{1616663640.0, "50050.0", "50150.0", "50000.0", "50100.0", "50075.0", "8.3", 80},
				},
				"last": 1616663640.0,
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetOHLC(context.Background(), "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetOHLC failed: %v", err)
	}

	if result.Last != 1616663640 {
		t.Errorf("Last = %d, want 1616663640", result.Last)
	}

	ohlc := result.Data["XXBTZUSD"]
	if len(ohlc) != 2 {
		t.Errorf("len(ohlc) = %d, want 2", len(ohlc))
	}

	if ohlc[0].Count != 100 {
		t.Errorf("Count = %d, want 100", ohlc[0].Count)
	}
}

func TestPublicService_GetOrderBook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/Depth" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBTZUSD": map[string]interface{}{
					"asks": [][]interface{}{
						{"50000.00000", "1.000", 1616663618},
						{"50001.00000", "2.000", 1616663619},
					},
					"bids": [][]interface{}{
						{"49999.00000", "1.500", 1616663618},
						{"49998.00000", "3.000", 1616663617},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetOrderBook(context.Background(), "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetOrderBook failed: %v", err)
	}

	book := result["XXBTZUSD"]
	if len(book.Asks) != 2 {
		t.Errorf("len(Asks) = %d, want 2", len(book.Asks))
	}
	if len(book.Bids) != 2 {
		t.Errorf("len(Bids) = %d, want 2", len(book.Bids))
	}
}

func TestPublicService_GetRecentTrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/Trades" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBTZUSD": [][]interface{}{
					{"50000.00000", "0.10000000", 1616663618.1234, "b", "m", "", 12345678.0},
					{"50001.00000", "0.20000000", 1616663619.5678, "s", "l", "", 12345679.0},
				},
				"last": "1616663619567800000",
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetRecentTrades(context.Background(), "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetRecentTrades failed: %v", err)
	}

	if result.Last != "1616663619567800000" {
		t.Errorf("Last = %s, want 1616663619567800000", result.Last)
	}

	if len(result.Trades) != 2 {
		t.Errorf("len(Trades) = %d, want 2", len(result.Trades))
	}

	if result.Trades[0].Side != "b" {
		t.Errorf("Trades[0].Side = %s, want b", result.Trades[0].Side)
	}
}

func TestPublicService_GetRecentSpreads(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/Spread" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"XXBTZUSD": [][]interface{}{
					{1616663618.0, "49999.00000", "50000.00000"},
					{1616663619.0, "49998.00000", "50001.00000"},
				},
				"last": 1616663619.0,
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetRecentSpreads(context.Background(), "XBTUSD", nil)
	if err != nil {
		t.Fatalf("GetRecentSpreads failed: %v", err)
	}

	if result.Last != 1616663619 {
		t.Errorf("Last = %d, want 1616663619", result.Last)
	}

	spreads := result.Data["XXBTZUSD"]
	if len(spreads) != 2 {
		t.Errorf("len(spreads) = %d, want 2", len(spreads))
	}
}

func TestPublicService_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  []string{"EGeneral:Invalid arguments"},
			"result": nil,
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	_, err := service.GetServerTime(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Category != "EGeneral" {
		t.Errorf("Category = %s, want EGeneral", apiErr.Category)
	}
}

func TestPublicService_GetGroupedOrderBook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/GroupedBook" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Check query params
		if pair := r.URL.Query().Get("pair"); pair != "BTC/USD" {
			t.Errorf("pair = %s, want BTC/USD", pair)
		}
		if depth := r.URL.Query().Get("depth"); depth != "10" {
			t.Errorf("depth = %s, want 10", depth)
		}
		if grouping := r.URL.Query().Get("grouping"); grouping != "1000" {
			t.Errorf("grouping = %s, want 1000", grouping)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"pair":     "BTC/USD",
				"grouping": 1000,
				"bids": []map[string]interface{}{
					{"price": "90400.00000", "qty": "19.83057746"},
					{"price": "90300.00000", "qty": "45.35073006"},
				},
				"asks": []map[string]interface{}{
					{"price": "90500.00000", "qty": "38.96185061"},
					{"price": "90600.00000", "qty": "55.96402032"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetGroupedOrderBook(context.Background(), "BTC/USD", &GetGroupedOrderBookOptions{
		Depth:    10,
		Grouping: 1000,
	})
	if err != nil {
		t.Fatalf("GetGroupedOrderBook failed: %v", err)
	}

	if result.Pair != "BTC/USD" {
		t.Errorf("Pair = %s, want BTC/USD", result.Pair)
	}
	if result.Grouping != 1000 {
		t.Errorf("Grouping = %d, want 1000", result.Grouping)
	}
	if len(result.Bids) != 2 {
		t.Errorf("len(Bids) = %d, want 2", len(result.Bids))
	}
	if len(result.Asks) != 2 {
		t.Errorf("len(Asks) = %d, want 2", len(result.Asks))
	}
}

func TestPublicService_GetPreTradeData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/PreTrade" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if symbol := r.URL.Query().Get("symbol"); symbol != "BTC/USD" {
			t.Errorf("symbol = %s, want BTC/USD", symbol)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"symbol":         "BTC/USD",
				"description":    "Bitcoin / US Dollars",
				"base_asset":     "BTC",
				"base_notation":  "NOML",
				"quote_asset":    "USD",
				"quote_notation": "MONE",
				"venue":          "PGSL",
				"system":         "CLOB",
				"bids": []map[string]interface{}{
					{"side": "BUY", "price": "102002.1", "qty": "1.5", "count": 10, "publication_ts": "2024-05-30T12:34:56.123456Z"},
				},
				"asks": []map[string]interface{}{
					{"side": "SELL", "price": "102003.1", "qty": "2.0", "count": 5, "publication_ts": "2024-05-30T12:34:56.123456Z"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetPreTradeData(context.Background(), "BTC/USD")
	if err != nil {
		t.Fatalf("GetPreTradeData failed: %v", err)
	}

	if result.Symbol != "BTC/USD" {
		t.Errorf("Symbol = %s, want BTC/USD", result.Symbol)
	}
	if result.BaseAsset != "BTC" {
		t.Errorf("BaseAsset = %s, want BTC", result.BaseAsset)
	}
	if len(result.Bids) != 1 {
		t.Errorf("len(Bids) = %d, want 1", len(result.Bids))
	}
	if len(result.Asks) != 1 {
		t.Errorf("len(Asks) = %d, want 1", len(result.Asks))
	}
}

func TestPublicService_GetPostTradeData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/public/PostTrade" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if symbol := r.URL.Query().Get("symbol"); symbol != "BTC/USD" {
			t.Errorf("symbol = %s, want BTC/USD", symbol)
		}
		if count := r.URL.Query().Get("count"); count != "100" {
			t.Errorf("count = %s, want 100", count)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"last_ts": "2024-05-30T12:34:56.123456789Z",
				"count":   2,
				"trades": []map[string]interface{}{
					{
						"trade_id":          "TGBB7L-HT5LX-J3BZ4A",
						"price":             "102002.1",
						"quantity":          "1.24",
						"symbol":            "BTC/USD",
						"description":       "Bitcoin / US Dollars",
						"base_asset":        "BTC",
						"base_notation":     "UNIT",
						"quote_asset":       "USD",
						"quote_notation":    "MONE",
						"trade_venue":       "PGSL",
						"trade_ts":          "2024-05-30T12:34:56.123456789Z",
						"publication_venue": "PGSL",
						"publication_ts":    "2024-05-30T12:34:56.123456789Z",
					},
					{
						"trade_id":          "TGBB7L-HT5LX-J3BZ4B",
						"price":             "102003.5",
						"quantity":          "0.5",
						"symbol":            "BTC/USD",
						"description":       "Bitcoin / US Dollars",
						"base_asset":        "BTC",
						"base_notation":     "UNIT",
						"quote_asset":       "USD",
						"quote_notation":    "MONE",
						"trade_venue":       "PGSL",
						"trade_ts":          "2024-05-30T12:34:57.123456789Z",
						"publication_venue": "PGSL",
						"publication_ts":    "2024-05-30T12:34:57.123456789Z",
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	service := NewPublicService(client)

	result, err := service.GetPostTradeData(context.Background(), &GetPostTradeOptions{
		Symbol: "BTC/USD",
		Count:  100,
	})
	if err != nil {
		t.Fatalf("GetPostTradeData failed: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("Count = %d, want 2", result.Count)
	}
	if len(result.Trades) != 2 {
		t.Errorf("len(Trades) = %d, want 2", len(result.Trades))
	}
	if result.Trades[0].TradeID != "TGBB7L-HT5LX-J3BZ4A" {
		t.Errorf("Trades[0].TradeID = %s, want TGBB7L-HT5LX-J3BZ4A", result.Trades[0].TradeID)
	}
}
