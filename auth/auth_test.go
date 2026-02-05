package auth

import (
	"net/http"
	"net/url"
	"testing"
)

func TestAuthenticator_Authenticate(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="

	auth, err := NewAuthenticator(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	req, _ := http.NewRequest("POST", "https://api.kraken.com/0/private/Balance", nil)
	params := url.Values{}

	err = auth.Authenticate(req, "/0/private/Balance", params)
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}

	// Check headers are set
	if req.Header.Get("API-Key") != apiKey {
		t.Errorf("API-Key header = %s, want %s", req.Header.Get("API-Key"), apiKey)
	}

	if req.Header.Get("API-Sign") == "" {
		t.Error("API-Sign header should not be empty")
	}

	// Check nonce was added to params
	if params.Get("nonce") == "" {
		t.Error("nonce should be added to params")
	}
}

func TestAuthenticator_WithFixedNonce(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="
	fixedNonce := "1616492376594"

	auth, err := NewAuthenticator(apiKey, apiSecret,
		WithNonceGenerator(NewFixedNonceGenerator(fixedNonce)),
	)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	req, _ := http.NewRequest("POST", "https://api.kraken.com/0/private/AddOrder", nil)
	params := url.Values{}
	params.Set("ordertype", "limit")
	params.Set("pair", "XBTUSD")
	params.Set("price", "37500")
	params.Set("type", "buy")
	params.Set("volume", "1.25")

	err = auth.Authenticate(req, "/0/private/AddOrder", params)
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}

	// With fixed nonce, we should get the documented signature
	expected := "4/dpxb3iT4tp/ZCVEwSnEsLxx0bqyhLpdfOpc6fn7OR8+UClSV5n9E6aSS8MPtnRfp32bAb0nmbRn6H8ndwLUQ=="
	if req.Header.Get("API-Sign") != expected {
		t.Errorf("API-Sign = %s, want %s", req.Header.Get("API-Sign"), expected)
	}
}

func TestAuthenticator_WithOTP(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="
	otpValue := "123456"

	auth, err := NewAuthenticator(apiKey, apiSecret,
		WithOTP(func() string { return otpValue }),
	)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	req, _ := http.NewRequest("POST", "https://api.kraken.com/0/private/Balance", nil)
	params := url.Values{}

	err = auth.Authenticate(req, "/0/private/Balance", params)
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}

	// Check OTP was added
	if params.Get("otp") != otpValue {
		t.Errorf("otp = %s, want %s", params.Get("otp"), otpValue)
	}

	// Check HasOTP returns true
	if !auth.HasOTP() {
		t.Error("HasOTP should return true")
	}
}

func TestAuthenticator_NoOTP(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="

	auth, err := NewAuthenticator(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	req, _ := http.NewRequest("POST", "https://api.kraken.com/0/private/Balance", nil)
	params := url.Values{}

	err = auth.Authenticate(req, "/0/private/Balance", params)
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}

	// Check OTP was NOT added
	if params.Get("otp") != "" {
		t.Errorf("otp should not be set, got %s", params.Get("otp"))
	}

	// Check HasOTP returns false
	if auth.HasOTP() {
		t.Error("HasOTP should return false")
	}
}

func TestAuthenticator_GenerateNonce(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="

	auth, err := NewAuthenticator(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	nonce := auth.GenerateNonce()
	if nonce == "" {
		t.Error("GenerateNonce should return non-empty string")
	}
}

func TestAuthenticator_APIKey(t *testing.T) {
	apiKey := "my-test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="

	auth, err := NewAuthenticator(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	if auth.APIKey() != apiKey {
		t.Errorf("APIKey() = %s, want %s", auth.APIKey(), apiKey)
	}
}

func TestAuthenticator_AuthenticateJSON(t *testing.T) {
	apiKey := "test-api-key"
	apiSecret := "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg=="

	auth, err := NewAuthenticator(apiKey, apiSecret)
	if err != nil {
		t.Fatalf("NewAuthenticator failed: %v", err)
	}

	req, _ := http.NewRequest("POST", "https://api.kraken.com/0/private/AddOrder", nil)
	nonce := "1616492376594"
	jsonBody := `{"nonce":1616492376594,"ordertype":"limit"}`

	err = auth.AuthenticateJSON(req, "/0/private/AddOrder", nonce, jsonBody)
	if err != nil {
		t.Fatalf("AuthenticateJSON failed: %v", err)
	}

	if req.Header.Get("API-Key") != apiKey {
		t.Errorf("API-Key header = %s, want %s", req.Header.Get("API-Key"), apiKey)
	}

	if req.Header.Get("API-Sign") == "" {
		t.Error("API-Sign header should not be empty")
	}
}
