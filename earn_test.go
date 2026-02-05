package kraken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEarnService_GetStrategies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Earn/Strategies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("asset") != "DOT" {
			t.Errorf("asset = %s, want DOT", r.Form.Get("asset"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"id":    "DOT-staking",
						"asset": "DOT",
						"lock_type": map[string]interface{}{
							"type":             "bonded",
							"payout_frequency": 7,
							"bonding_period":   2419200,
							"unbonding_period": 2419200,
						},
						"apr_estimate": map[string]interface{}{
							"low":  "0.08",
							"high": "0.12",
						},
						"user_min_allocation": "1.0",
						"auto_compound": map[string]interface{}{
							"type":    "optional",
							"default": true,
						},
						"can_allocate":   true,
						"can_deallocate": true,
					},
				},
				"next_cursor": "",
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

	result, err := client.Earn.GetStrategies(context.Background(), &GetStrategiesOptions{
		Asset: "DOT",
	})
	if err != nil {
		t.Fatalf("GetStrategies failed: %v", err)
	}

	if len(result.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(result.Items))
	}
	if result.Items[0].ID != "DOT-staking" {
		t.Errorf("ID = %s, want DOT-staking", result.Items[0].ID)
	}
	if result.Items[0].Asset != "DOT" {
		t.Errorf("Asset = %s, want DOT", result.Items[0].Asset)
	}
}

func TestEarnService_GetAllocations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Earn/Allocations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"converted_asset": "USD",
				"total_allocated": "1000.00",
				"total_rewarded":  "50.00",
				"items": []interface{}{
					map[string]interface{}{
						"strategy_id":  "DOT-staking",
						"native_asset": "DOT",
						"amount_allocated": map[string]interface{}{
							"total": map[string]interface{}{
								"native":    "100.0",
								"converted": "50.0",
							},
						},
						"total_rewarded": map[string]interface{}{
							"total": map[string]interface{}{
								"native":    "5.0",
								"converted": "2.5",
							},
						},
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

	result, err := client.Earn.GetAllocations(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetAllocations failed: %v", err)
	}

	if result.ConvertedAsset != "USD" {
		t.Errorf("ConvertedAsset = %s, want USD", result.ConvertedAsset)
	}
	if len(result.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(result.Items))
	}
	if result.Items[0].StrategyID != "DOT-staking" {
		t.Errorf("StrategyID = %s, want DOT-staking", result.Items[0].StrategyID)
	}
}

func TestEarnService_Allocate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Earn/Allocate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("strategy_id") != "DOT-staking" {
			t.Errorf("strategy_id = %s, want DOT-staking", r.Form.Get("strategy_id"))
		}
		if r.Form.Get("amount") != "10.0" {
			t.Errorf("amount = %s, want 10.0", r.Form.Get("amount"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"pending": true,
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

	result, err := client.Earn.Allocate(context.Background(), &AllocateRequest{
		StrategyID: "DOT-staking",
		Amount:     "10.0",
	})
	if err != nil {
		t.Fatalf("Allocate failed: %v", err)
	}

	if !result.Pending {
		t.Error("Pending = false, want true")
	}
}

func TestEarnService_Deallocate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Earn/Deallocate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("strategy_id") != "DOT-staking" {
			t.Errorf("strategy_id = %s, want DOT-staking", r.Form.Get("strategy_id"))
		}
		if r.Form.Get("amount") != "5.0" {
			t.Errorf("amount = %s, want 5.0", r.Form.Get("amount"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"pending": true,
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

	result, err := client.Earn.Deallocate(context.Background(), &DeallocateRequest{
		StrategyID: "DOT-staking",
		Amount:     "5.0",
	})
	if err != nil {
		t.Fatalf("Deallocate failed: %v", err)
	}

	if !result.Pending {
		t.Error("Pending = false, want true")
	}
}

func TestEarnService_GetAllocationStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Earn/AllocateStatus" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("strategy_id") != "DOT-staking" {
			t.Errorf("strategy_id = %s, want DOT-staking", r.Form.Get("strategy_id"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"pending": false,
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

	result, err := client.Earn.GetAllocationStatus(context.Background(), "DOT-staking")
	if err != nil {
		t.Fatalf("GetAllocationStatus failed: %v", err)
	}

	if result.Pending {
		t.Error("Pending = true, want false")
	}
}

func TestEarnService_GetDeallocationStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/0/private/Earn/DeallocateStatus" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("strategy_id") != "DOT-staking" {
			t.Errorf("strategy_id = %s, want DOT-staking", r.Form.Get("strategy_id"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": []string{},
			"result": map[string]interface{}{
				"pending": true,
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

	result, err := client.Earn.GetDeallocationStatus(context.Background(), "DOT-staking")
	if err != nil {
		t.Fatalf("GetDeallocationStatus failed: %v", err)
	}

	if !result.Pending {
		t.Error("Pending = false, want true")
	}
}
