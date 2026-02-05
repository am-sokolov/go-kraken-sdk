package kraken

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors for common error categories.
var (
	// ErrGeneral represents general errors (EGeneral prefix).
	ErrGeneral = errors.New("general error")

	// ErrService represents service errors (EService prefix).
	ErrService = errors.New("service error")

	// ErrAPI represents API errors (EAPI prefix).
	ErrAPI = errors.New("API error")

	// ErrOrder represents order errors (EOrder prefix).
	ErrOrder = errors.New("order error")

	// ErrAuth represents authentication errors (EAuth prefix).
	ErrAuth = errors.New("authentication error")

	// ErrTrade represents trade errors (ETrade prefix).
	ErrTrade = errors.New("trade error")

	// ErrFunding represents funding errors (EFunding prefix).
	ErrFunding = errors.New("funding error")

	// ErrRateLimit represents rate limit exceeded errors.
	ErrRateLimit = errors.New("rate limit exceeded")

	// ErrInsufficientFunds represents insufficient funds error.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidNonce represents invalid nonce error.
	ErrInvalidNonce = errors.New("invalid nonce")

	// ErrInvalidKey represents invalid API key error.
	ErrInvalidKey = errors.New("invalid API key")

	// ErrInvalidSignature represents invalid signature error.
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrPermissionDenied represents permission denied error.
	ErrPermissionDenied = errors.New("permission denied")

	// ErrOrderMinNotMet represents order minimum not met error.
	ErrOrderMinNotMet = errors.New("order minimum not met")

	// ErrUnknownOrder represents unknown order error.
	ErrUnknownOrder = errors.New("unknown order")
)

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

// Is implements errors.Is for comparing with sentinel errors.
func (e *APIError) Is(target error) bool {
	switch {
	case errors.Is(target, ErrGeneral):
		return e.Category == "EGeneral"
	case errors.Is(target, ErrService):
		return e.Category == "EService"
	case errors.Is(target, ErrAPI):
		return e.Category == "EAPI"
	case errors.Is(target, ErrOrder):
		return e.Category == "EOrder"
	case errors.Is(target, ErrAuth):
		return e.Category == "EAuth"
	case errors.Is(target, ErrTrade):
		return e.Category == "ETrade"
	case errors.Is(target, ErrFunding):
		return e.Category == "EFunding"
	case errors.Is(target, ErrRateLimit):
		return e.containsMessage("Rate limit exceeded")
	case errors.Is(target, ErrInsufficientFunds):
		return e.containsMessage("Insufficient funds")
	case errors.Is(target, ErrInvalidNonce):
		return e.containsMessage("Invalid nonce")
	case errors.Is(target, ErrInvalidKey):
		return e.containsMessage("Invalid key")
	case errors.Is(target, ErrInvalidSignature):
		return e.containsMessage("Invalid signature")
	case errors.Is(target, ErrPermissionDenied):
		return e.containsMessage("Permission denied")
	case errors.Is(target, ErrOrderMinNotMet):
		return e.containsMessage("Order minimum not met")
	case errors.Is(target, ErrUnknownOrder):
		return e.containsMessage("Unknown order")
	}
	return false
}

// Unwrap returns the base error category.
func (e *APIError) Unwrap() error {
	switch e.Category {
	case "EGeneral":
		return ErrGeneral
	case "EService":
		return ErrService
	case "EAPI":
		return ErrAPI
	case "EOrder":
		return ErrOrder
	case "EAuth":
		return ErrAuth
	case "ETrade":
		return ErrTrade
	case "EFunding":
		return ErrFunding
	default:
		return nil
	}
}

// containsMessage checks if any message contains the given substring.
func (e *APIError) containsMessage(substr string) bool {
	for _, msg := range e.Messages {
		if strings.Contains(msg, substr) {
			return true
		}
	}
	return false
}

// WebSocketError represents a WebSocket-specific error.
type WebSocketError struct {
	// Method is the WebSocket method that failed.
	Method string

	// Message is the error message.
	Message string

	// ReqID is the request ID, if available.
	ReqID int64
}

// Error implements the error interface.
func (e *WebSocketError) Error() string {
	if e.ReqID != 0 {
		return fmt.Sprintf("kraken websocket error [%s] (req_id=%d): %s", e.Method, e.ReqID, e.Message)
	}
	return fmt.Sprintf("kraken websocket error [%s]: %s", e.Method, e.Message)
}

// ParseAPIError parses Kraken error messages into a structured APIError.
func ParseAPIError(messages []string) error {
	if len(messages) == 0 {
		return nil
	}

	apiErr := &APIError{
		Messages: messages,
	}

	// Extract category from first message
	if len(messages) > 0 {
		parts := strings.SplitN(messages[0], ":", 2)
		if len(parts) > 0 {
			apiErr.Category = parts[0]
		}
	}

	return apiErr
}

// IsRetryable returns true if the error is temporary and the request should be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		// Service errors are typically temporary
		if apiErr.Category == "EService" {
			return true
		}

		// Rate limit errors should be retried after waiting
		for _, msg := range apiErr.Messages {
			if strings.Contains(msg, "Rate limit exceeded") {
				return true
			}
			if strings.Contains(msg, "Busy") {
				return true
			}
			if strings.Contains(msg, "Timeout") {
				return true
			}
		}

		// HTTP 5xx errors
		if apiErr.HTTPStatus >= 500 {
			return true
		}

		// HTTP 429 (Too Many Requests)
		if apiErr.HTTPStatus == 429 {
			return true
		}
	}

	return false
}

// IsAuthError returns true if the error is related to authentication.
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.Category == "EAPI" || apiErr.Category == "EAuth" {
			return true
		}

		for _, msg := range apiErr.Messages {
			if strings.Contains(msg, "Invalid key") ||
				strings.Contains(msg, "Invalid signature") ||
				strings.Contains(msg, "Invalid nonce") ||
				strings.Contains(msg, "Permission denied") {
				return true
			}
		}
	}

	return false
}

// IsOrderError returns true if the error is related to order operations.
func IsOrderError(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Category == "EOrder"
	}

	return false
}
