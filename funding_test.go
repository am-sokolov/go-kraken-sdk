package kraken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFundingService_GetDepositMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/DepositMethods" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XBT" {
			t.Errorf("asset = %s, want XBT", r.Form.Get("asset"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []interface{}{
				map[string]interface{}{
					"method":      "Bitcoin",
					"limit":       false,
					"fee":         "0",
					"gen-address": true,
					"minimum":     "0.0001",
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

	result, err := client.Funding.GetDepositMethods(context.Background(), &GetDepositMethodsOptions{
		Asset: "XBT",
	})
	if err != nil {
		t.Fatalf("GetDepositMethods failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if result[0].Method != "Bitcoin" {
		t.Errorf("Method = %s, want Bitcoin", result[0].Method)
	}
}

func TestFundingService_GetDepositAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/DepositAddresses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XBT" {
			t.Errorf("asset = %s, want XBT", r.Form.Get("asset"))
		}
		if r.Form.Get("method") != "Bitcoin" {
			t.Errorf("method = %s, want Bitcoin", r.Form.Get("method"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []interface{}{
				map[string]interface{}{
					"address":  "bc1q...",
					"expiretm": "0",
					"new":      false,
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

	result, err := client.Funding.GetDepositAddresses(context.Background(), &GetDepositAddressesOptions{
		Asset:  "XBT",
		Method: "Bitcoin",
	})
	if err != nil {
		t.Fatalf("GetDepositAddresses failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if result[0].Address != "bc1q..." {
		t.Errorf("Address = %s, want bc1q...", result[0].Address)
	}
}

func TestFundingService_GetDepositStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/DepositStatus" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []interface{}{
				map[string]interface{}{
					"method": "Bitcoin",
					"aclass": "currency",
					"asset":  "XXBT",
					"refid":  "DEPOSIT-REF-123",
					"txid":   "abc123...",
					"info":   "bc1q...",
					"amount": "0.5",
					"fee":    "0",
					"time":   1616663618,
					"status": "Success",
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

	result, err := client.Funding.GetDepositStatus(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetDepositStatus failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if result[0].Status != "Success" {
		t.Errorf("Status = %s, want Success", result[0].Status)
	}
}

func TestFundingService_GetWithdrawalMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/WithdrawMethods" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []interface{}{
				map[string]interface{}{
					"asset":   "XXBT",
					"method":  "Bitcoin",
					"network": "Bitcoin",
					"minimum": "0.0005",
					"maximum": "100",
					"fee":     "0.00015",
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

	result, err := client.Funding.GetWithdrawalMethods(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetWithdrawalMethods failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if result[0].Method != "Bitcoin" {
		t.Errorf("Method = %s, want Bitcoin", result[0].Method)
	}
}

func TestFundingService_GetWithdrawalAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/WithdrawAddresses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []interface{}{
				map[string]interface{}{
					"address":  "bc1q...",
					"asset":    "XBT",
					"method":   "Bitcoin",
					"key":      "my-btc-wallet",
					"verified": true,
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

	result, err := client.Funding.GetWithdrawalAddresses(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetWithdrawalAddresses failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if result[0].Key != "my-btc-wallet" {
		t.Errorf("Key = %s, want my-btc-wallet", result[0].Key)
	}
}

func TestFundingService_GetWithdrawalInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/WithdrawInfo" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XBT" {
			t.Errorf("asset = %s, want XBT", r.Form.Get("asset"))
		}
		if r.Form.Get("key") != "my-btc-wallet" {
			t.Errorf("key = %s, want my-btc-wallet", r.Form.Get("key"))
		}
		if r.Form.Get("amount") != "0.5" {
			t.Errorf("amount = %s, want 0.5", r.Form.Get("amount"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"method": "Bitcoin",
				"limit":  "100.0",
				"amount": "0.5",
				"fee":    "0.00015",
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

	result, err := client.Funding.GetWithdrawalInfo(context.Background(), &GetWithdrawalInfoRequest{
		Asset:  "XBT",
		Key:    "my-btc-wallet",
		Amount: "0.5",
	})
	if err != nil {
		t.Fatalf("GetWithdrawalInfo failed: %v", err)
	}

	if result.Method != "Bitcoin" {
		t.Errorf("Method = %s, want Bitcoin", result.Method)
	}
}

func TestFundingService_Withdraw(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Withdraw" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XBT" {
			t.Errorf("asset = %s, want XBT", r.Form.Get("asset"))
		}
		if r.Form.Get("key") != "my-btc-wallet" {
			t.Errorf("key = %s, want my-btc-wallet", r.Form.Get("key"))
		}
		if r.Form.Get("amount") != "0.5" {
			t.Errorf("amount = %s, want 0.5", r.Form.Get("amount"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"refid": "WITHDRAW-REF-123",
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

	result, err := client.Funding.Withdraw(context.Background(), &WithdrawRequest{
		Asset:  "XBT",
		Key:    "my-btc-wallet",
		Amount: "0.5",
	})
	if err != nil {
		t.Fatalf("Withdraw failed: %v", err)
	}

	if result.RefID != "WITHDRAW-REF-123" {
		t.Errorf("RefID = %s, want WITHDRAW-REF-123", result.RefID)
	}
}

func TestFundingService_GetWithdrawalStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/WithdrawStatus" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": []interface{}{
				map[string]interface{}{
					"method": "Bitcoin",
					"aclass": "currency",
					"asset":  "XXBT",
					"refid":  "WITHDRAW-REF-123",
					"txid":   "abc123...",
					"info":   "bc1q...",
					"amount": "0.5",
					"fee":    "0.00015",
					"time":   1616663618,
					"status": "Success",
					"key":    "my-btc-wallet",
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

	result, err := client.Funding.GetWithdrawalStatus(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetWithdrawalStatus failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}
	if result[0].Status != "Success" {
		t.Errorf("Status = %s, want Success", result[0].Status)
	}
}

func TestFundingService_CancelWithdrawal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/WithdrawCancel" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XBT" {
			t.Errorf("asset = %s, want XBT", r.Form.Get("asset"))
		}
		if r.Form.Get("refid") != "WITHDRAW-REF-123" {
			t.Errorf("refid = %s, want WITHDRAW-REF-123", r.Form.Get("refid"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  []string{},
			"result": true,
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

	result, err := client.Funding.CancelWithdrawal(context.Background(), "XBT", "WITHDRAW-REF-123")
	if err != nil {
		t.Fatalf("CancelWithdrawal failed: %v", err)
	}

	if !result {
		t.Error("CancelWithdrawal returned false, want true")
	}
}

func TestFundingService_WalletTransfer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/WalletTransfer" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "XBT" {
			t.Errorf("asset = %s, want XBT", r.Form.Get("asset"))
		}
		if r.Form.Get("from") != "Spot Wallet" {
			t.Errorf("from = %s, want Spot Wallet", r.Form.Get("from"))
		}
		if r.Form.Get("to") != "Futures Wallet" {
			t.Errorf("to = %s, want Futures Wallet", r.Form.Get("to"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"refid": "TRANSFER-REF-123",
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

	result, err := client.Funding.WalletTransfer(context.Background(), &WalletTransferRequest{
		Asset:  "XBT",
		From:   "Spot Wallet",
		To:     "Futures Wallet",
		Amount: "0.5",
	})
	if err != nil {
		t.Fatalf("WalletTransfer failed: %v", err)
	}

	if result.RefID != "TRANSFER-REF-123" {
		t.Errorf("RefID = %s, want TRANSFER-REF-123", result.RefID)
	}
}
