package auth

import (
	"net/url"
	"testing"
)

// TestSigner_Sign tests the signature generation against Kraken's documented test vector.
// This test uses the exact example from:
// https://docs.kraken.com/api/docs/guides/spot-rest-auth
func TestSigner_Sign(t *testing.T) {
	// Test vector from Kraken documentation
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="
	apiKey := "test-api-key"

	signer, err := NewSigner(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	nonce := "1616492376594"
	path := "/0/private/AddOrder"

	// Create params in the exact order from documentation
	params := url.Values{}
	params.Set("nonce", nonce)
	params.Set("ordertype", "limit")
	params.Set("pair", "XBTUSD")
	params.Set("price", "37500")
	params.Set("type", "buy")
	params.Set("volume", "1.25")

	// Expected signature from Kraken documentation
	expected := "4/dpxb3iT4tp/ZCVEwSnEsLxx0bqyhLpdfOpc6fn7OR8+UClSV5n9E6aSS8MPtnRfp32bAb0nmbRn6H8ndwLUQ=="

	signature, err := signer.Sign(path, nonce, params)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if signature != expected {
		t.Errorf("signature mismatch:\ngot:  %s\nwant: %s", signature, expected)
	}
}

// TestSigner_SignJSON tests JSON body signature generation.
func TestSigner_SignJSON(t *testing.T) {
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="
	apiKey := "test-api-key"

	signer, err := NewSigner(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	nonce := "1616492376594"
	path := "/0/private/AddOrder"

	// JSON body (the signature will be different from form-encoded)
	jsonBody := `{"nonce":1616492376594,"ordertype":"limit","pair":"XBTUSD","price":"37500","type":"buy","volume":"1.25"}`

	signature, err := signer.SignJSON(path, nonce, jsonBody)
	if err != nil {
		t.Fatalf("SignJSON failed: %v", err)
	}

	// Just verify we get a non-empty signature
	// The exact value will differ from form-encoded signature
	if signature == "" {
		t.Error("expected non-empty signature")
	}
}

// TestSigner_NewSigner_EmptyKey tests that empty API key is rejected.
func TestSigner_NewSigner_EmptyKey(t *testing.T) {
	_, err := NewSigner("", "secret")
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

// TestSigner_NewSigner_EmptySecret tests that empty API secret is rejected.
func TestSigner_NewSigner_EmptySecret(t *testing.T) {
	_, err := NewSigner("key", "")
	if err == nil {
		t.Error("expected error for empty API secret")
	}
}

// TestSigner_NewSigner_InvalidBase64 tests that invalid base64 secret is rejected.
func TestSigner_NewSigner_InvalidBase64(t *testing.T) {
	_, err := NewSigner("key", "not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64 secret")
	}
}

// TestSigner_APIKey tests that the API key is correctly returned.
func TestSigner_APIKey(t *testing.T) {
	apiKey := "my-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="

	signer, err := NewSigner(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	if signer.APIKey() != apiKey {
		t.Errorf("APIKey() = %s, want %s", signer.APIKey(), apiKey)
	}
}

// TestEncodeParams tests URL parameter encoding.
func TestEncodeParams(t *testing.T) {
	tests := []struct {
		name     string
		params   url.Values
		expected string
	}{
		{
			name:     "empty params",
			params:   url.Values{},
			expected: "",
		},
		{
			name: "single param",
			params: url.Values{
				"nonce": {"123"},
			},
			expected: "nonce=123",
		},
		{
			name: "multiple params sorted",
			params: url.Values{
				"z": {"last"},
				"a": {"first"},
				"m": {"middle"},
			},
			expected: "a=first&m=middle&z=last",
		},
		{
			name: "special characters encoded",
			params: url.Values{
				"key": {"value with spaces"},
			},
			expected: "key=value+with+spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeParams(tt.params)
			if result != tt.expected {
				t.Errorf("encodeParams() = %q, want %q", result, tt.expected)
			}
		})
	}
}
