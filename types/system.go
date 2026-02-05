package types

import "time"

// ServerTime represents the server time response.
type ServerTime struct {
	// UnixTime is the Unix timestamp.
	UnixTime int64 `json:"unixtime"`

	// RFC1123 is the RFC1123 formatted time.
	RFC1123 string `json:"rfc1123"`
}

// SystemStatus represents the system status response.
type SystemStatus struct {
	// Status is the current status (online, maintenance, etc.).
	Status string `json:"status"`

	// Timestamp is the status timestamp.
	Timestamp time.Time `json:"timestamp"`
}

// WSStatus represents WebSocket system status.
type WSStatus struct {
	// API is the API version.
	API string `json:"api_version"`

	// ConnectionID is the connection identifier.
	ConnectionID uint64 `json:"connection_id"`

	// System is the system status.
	System string `json:"system"`

	// Version is the WebSocket version.
	Version string `json:"version"`
}

// WebSocketToken represents the token for WebSocket authentication.
type WebSocketToken struct {
	// Token is the authentication token.
	Token string `json:"token"`

	// Expires is the expiration time in seconds.
	Expires int64 `json:"expires"`
}

// SubaccountResult represents the result of creating a subaccount.
type SubaccountResult struct {
	// Result is true if successful.
	Result bool `json:"result"`
}
