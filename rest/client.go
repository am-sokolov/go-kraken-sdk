// Package rest provides the REST API client for Kraken.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/am-sokolov/go-kraken-sdk/auth"
)

// Client handles REST API communication with Kraken.
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *auth.Authenticator
	userAgent  string

	// Rate limiting (optional)
	// limiter ratelimit.Limiter

	// Callbacks
	onError func(error)

	// Debug
	debugWriter io.Writer
}

// ClientOption configures the REST client.
type ClientOption func(*Client)

// NewClient creates a new REST client.
func NewClient(baseURL string, httpClient *http.Client, opts ...ClientOption) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	c := &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: httpClient,
		userAgent:  "go-kraken-sdk/1.0",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithAuth sets the authenticator for private endpoints.
func WithAuth(authenticator *auth.Authenticator) ClientOption {
	return func(c *Client) {
		c.auth = authenticator
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithOnError sets the error callback.
func WithOnError(fn func(error)) ClientOption {
	return func(c *Client) {
		c.onError = fn
	}
}

// WithDebugWriter sets a writer for debug output.
func WithDebugWriter(w io.Writer) ClientOption {
	return func(c *Client) {
		c.debugWriter = w
	}
}

// DoPublic executes a public (unauthenticated) GET request.
func (c *Client) DoPublic(ctx context.Context, path string, params url.Values) (*Response, error) {
	// Build URL with query parameters
	reqURL := c.baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.do(req)
}

// DoPrivate executes a private (authenticated) POST request.
func (c *Client) DoPrivate(ctx context.Context, path string, params url.Values) (*Response, error) {
	if c.auth == nil {
		return nil, fmt.Errorf("authentication required for private endpoint")
	}

	if params == nil {
		params = url.Values{}
	}

	// Build request
	reqURL := c.baseURL + path
	body := params.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add authentication
	if err := c.auth.Authenticate(req, path, params); err != nil {
		return nil, fmt.Errorf("failed to authenticate request: %w", err)
	}

	// Update body with authenticated params (includes nonce)
	body = params.Encode()
	req.Body = io.NopCloser(strings.NewReader(body))
	req.ContentLength = int64(len(body))

	return c.do(req)
}

// DoPrivateJSON executes a private POST request with JSON body.
func (c *Client) DoPrivateJSON(ctx context.Context, path string, body interface{}) (*Response, error) {
	if c.auth == nil {
		return nil, fmt.Errorf("authentication required for private endpoint")
	}

	// Marshal body to JSON.
	// Treat nil as an empty JSON object (private endpoints expect an object payload).
	jsonBody := []byte("{}")
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Parse JSON object so we can inject nonce/otp without losing numeric precision.
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(jsonBody, &obj); err != nil {
		return nil, fmt.Errorf("failed to parse request body: %w", err)
	}
	if obj == nil {
		obj = map[string]json.RawMessage{}
	}

	// Extract or generate nonce (required for signing).
	nonceStr := ""
	if rawNonce, ok := obj["nonce"]; ok && len(bytes.TrimSpace(rawNonce)) > 0 && string(bytes.TrimSpace(rawNonce)) != "null" {
		trim := bytes.TrimSpace(rawNonce)
		if len(trim) > 0 && trim[0] == '"' {
			var nonceValue string
			if err := json.Unmarshal(trim, &nonceValue); err != nil {
				return nil, fmt.Errorf("invalid nonce field in JSON body: %w", err)
			}
			nonceStr = strings.TrimSpace(nonceValue)
		} else {
			nonceStr = string(trim)
		}

		nonceInt, err := strconv.ParseInt(nonceStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid nonce field in JSON body: %w", err)
		}
		nonceStr = strconv.FormatInt(nonceInt, 10)
		// Normalize nonce to an integer JSON value.
		obj["nonce"] = json.RawMessage(nonceStr)
	} else {
		nonceStr = c.auth.GenerateNonce()
		obj["nonce"] = json.RawMessage(nonceStr)
	}

	// Add OTP if configured (2FA).
	if c.auth.HasOTP() {
		if _, ok := obj["otp"]; !ok {
			otp := strings.TrimSpace(c.auth.GenerateOTP())
			if otp != "" {
				obj["otp"] = json.RawMessage(strconv.Quote(otp))
			}
		}
	}

	// Re-marshal normalized body (must match the payload used for signing).
	var err error
	jsonBody, err = json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Build request
	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication
	if err := c.auth.AuthenticateJSON(req, path, nonceStr, string(jsonBody)); err != nil {
		return nil, fmt.Errorf("failed to authenticate request: %w", err)
	}

	return c.do(req)
}

// do executes an HTTP request and returns the parsed response.
func (c *Client) do(req *http.Request) (*Response, error) {
	// Set common headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	// Debug logging
	if c.debugWriter != nil {
		fmt.Fprintf(c.debugWriter, ">>> %s %s\n", req.Method, req.URL)
		for k, v := range req.Header {
			fmt.Fprintf(c.debugWriter, ">>> %s: %s\n", k, v)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.onError != nil {
			c.onError(err)
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug logging
	if c.debugWriter != nil {
		fmt.Fprintf(c.debugWriter, "<<< %d %s\n", resp.StatusCode, resp.Status)
		// Only log body for JSON responses to avoid binary data in logs
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") || strings.HasPrefix(contentType, "text/") {
			fmt.Fprintf(c.debugWriter, "<<< %s\n", string(bodyBytes))
		} else {
			fmt.Fprintf(c.debugWriter, "<<< [binary data: %d bytes]\n", len(bodyBytes))
		}
	}

	// Check if response is binary (e.g., ZIP file from RetrieveExport)
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/octet-stream") ||
		strings.Contains(contentType, "application/zip") {
		// Return raw binary data
		return &Response{
			HTTPStatus: resp.StatusCode,
			RawData:    bodyBytes,
		}, nil
	}

	// Parse response as JSON
	var apiResp Response
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(bodyBytes))
	}

	apiResp.HTTPStatus = resp.StatusCode
	apiResp.RawData = bodyBytes

	// Check for API errors
	if apiResp.HasError() {
		apiErr := apiResp.GetError()
		if c.onError != nil {
			c.onError(apiErr)
		}
		return &apiResp, apiErr
	}

	return &apiResp, nil
}

// IsAuthenticated returns true if the client has authentication configured.
func (c *Client) IsAuthenticated() bool {
	return c.auth != nil
}

// GenerateNonce generates a new nonce value using the authenticator.
func (c *Client) GenerateNonce() string {
	if c.auth != nil {
		return c.auth.GenerateNonce()
	}
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}
