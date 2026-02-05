package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Signer generates API signatures for authenticated requests.
// The signature algorithm is: HMAC-SHA512(path + SHA256(nonce + postData), base64Decode(secret))
type Signer struct {
	apiKey    string
	apiSecret []byte // Base64-decoded secret
}

// NewSigner creates a new Signer with the given credentials.
// The apiSecret should be the base64-encoded secret from your Kraken API key-pair.
func NewSigner(apiKey, apiSecret string) (*Signer, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}
	if apiSecret == "" {
		return nil, fmt.Errorf("API secret cannot be empty")
	}

	// Decode the base64-encoded secret
	decoded, err := base64.StdEncoding.DecodeString(apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode API secret: %w", err)
	}

	return &Signer{
		apiKey:    apiKey,
		apiSecret: decoded,
	}, nil
}

// APIKey returns the public API key.
func (s *Signer) APIKey() string {
	return s.apiKey
}

// Sign generates the API-Sign header value for a request.
//
// The signature is computed as:
//  1. Concatenate nonce and POST data: nonce + urlEncode(params)
//  2. SHA256 hash of the concatenated string
//  3. Concatenate URI path and SHA256 hash
//  4. HMAC-SHA512 with the decoded API secret
//  5. Base64 encode the result
//
// Parameters:
//   - path: The URI path (e.g., "/0/private/AddOrder")
//   - nonce: The nonce value used in the request
//   - params: The POST parameters (must include the nonce)
//
// Returns the base64-encoded signature string.
func (s *Signer) Sign(path string, nonce string, params url.Values) (string, error) {
	// Ensure nonce is in params
	if params.Get("nonce") == "" {
		params.Set("nonce", nonce)
	}

	// Encode the POST data
	postData := encodeParams(params)

	// Step 1: Concatenate nonce and POST data
	noncePostData := nonce + postData

	// Step 2: SHA256 hash
	sha256Hash := sha256.Sum256([]byte(noncePostData))

	// Step 3: Concatenate path and SHA256 hash
	message := append([]byte(path), sha256Hash[:]...)

	// Step 4: HMAC-SHA512 with decoded secret
	mac := hmac.New(sha512.New, s.apiSecret)
	mac.Write(message)
	macSum := mac.Sum(nil)

	// Step 5: Base64 encode
	signature := base64.StdEncoding.EncodeToString(macSum)

	return signature, nil
}

// SignJSON generates the signature for JSON-encoded requests.
// This is used when sending JSON body instead of form-encoded data.
//
// Parameters:
//   - path: The URI path
//   - nonce: The nonce value
//   - jsonBody: The JSON-encoded request body as a string
//
// Note: The nonce must be present in the JSON body.
func (s *Signer) SignJSON(path string, nonce string, jsonBody string) (string, error) {
	// Step 1: Concatenate nonce and JSON body
	noncePostData := nonce + jsonBody

	// Step 2: SHA256 hash
	sha256Hash := sha256.Sum256([]byte(noncePostData))

	// Step 3: Concatenate path and SHA256 hash
	message := append([]byte(path), sha256Hash[:]...)

	// Step 4: HMAC-SHA512 with decoded secret
	mac := hmac.New(sha512.New, s.apiSecret)
	mac.Write(message)
	macSum := mac.Sum(nil)

	// Step 5: Base64 encode
	signature := base64.StdEncoding.EncodeToString(macSum)

	return signature, nil
}

// encodeParams encodes URL values in a consistent order.
// Kraken requires params to be in a specific order for signature verification.
func encodeParams(params url.Values) string {
	if len(params) == 0 {
		return ""
	}

	// Get sorted keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build encoded string
	var parts []string
	for _, key := range keys {
		values := params[key]
		for _, value := range values {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}

	return strings.Join(parts, "&")
}
