package rest

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Response represents a Kraken API response.
type Response struct {
	// Error contains any error messages from the API.
	Error []string `json:"error"`

	// Result contains the response data.
	Result json.RawMessage `json:"result"`

	// HTTPStatus is the HTTP status code.
	HTTPStatus int `json:"-"`

	// RawData contains the raw response body for binary responses.
	RawData []byte `json:"-"`
}

// HasError returns true if the response contains errors.
func (r *Response) HasError() bool {
	return len(r.Error) > 0
}

// GetError returns the parsed error or nil.
func (r *Response) GetError() error {
	if !r.HasError() {
		return nil
	}

	return &APIError{
		Category:   extractCategory(r.Error),
		Messages:   r.Error,
		HTTPStatus: r.HTTPStatus,
	}
}

// Decode decodes the result into the given value.
func (r *Response) Decode(v interface{}) error {
	if r.Result == nil {
		return nil
	}
	return json.Unmarshal(r.Result, v)
}

// APIError represents an error returned by the Kraken API.
type APIError struct {
	// Category is the error category (e.g., "EGeneral", "EOrder").
	Category string

	// Messages contains the error messages from the API.
	Messages []string

	// HTTPStatus is the HTTP status code, if applicable.
	HTTPStatus int
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Messages) == 0 {
		return fmt.Sprintf("kraken API error: %s", e.Category)
	}
	return fmt.Sprintf("kraken API error: %s", strings.Join(e.Messages, "; "))
}

// extractCategory extracts the error category from error messages.
func extractCategory(messages []string) string {
	if len(messages) == 0 {
		return ""
	}

	// Extract category from first message (e.g., "EOrder:Rate limit exceeded")
	parts := strings.SplitN(messages[0], ":", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// IsRetryable returns true if the error is temporary and the request should be retried.
func (e *APIError) IsRetryable() bool {
	// Service errors are typically temporary
	if e.Category == "EService" {
		return true
	}

	// Check for specific retryable conditions
	for _, msg := range e.Messages {
		if strings.Contains(msg, "Rate limit exceeded") ||
			strings.Contains(msg, "Busy") ||
			strings.Contains(msg, "Timeout") {
			return true
		}
	}

	// HTTP 5xx errors
	if e.HTTPStatus >= 500 {
		return true
	}

	// HTTP 429 (Too Many Requests)
	if e.HTTPStatus == 429 {
		return true
	}

	return false
}

// IsAuthError returns true if the error is related to authentication.
func (e *APIError) IsAuthError() bool {
	if e.Category == "EAPI" || e.Category == "EAuth" {
		return true
	}

	for _, msg := range e.Messages {
		if strings.Contains(msg, "Invalid key") ||
			strings.Contains(msg, "Invalid signature") ||
			strings.Contains(msg, "Invalid nonce") ||
			strings.Contains(msg, "Permission denied") {
			return true
		}
	}

	return false
}
