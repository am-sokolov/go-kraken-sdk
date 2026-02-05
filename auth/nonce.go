// Package auth provides authentication utilities for the Kraken API.
package auth

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// NonceGenerator generates unique, always-increasing nonce values.
type NonceGenerator interface {
	// Next returns the next nonce value as a string.
	Next() string
}

// TimestampNonceGenerator generates nonces based on Unix timestamp in milliseconds.
// This is the recommended approach for production use.
type TimestampNonceGenerator struct {
	lastNonce int64
	mu        sync.Mutex
}

// NewTimestampNonceGenerator creates a new timestamp-based nonce generator.
func NewTimestampNonceGenerator() *TimestampNonceGenerator {
	return &TimestampNonceGenerator{}
}

// Next returns the next nonce value.
// The value is guaranteed to be always increasing, even if called multiple times
// within the same millisecond.
func (g *TimestampNonceGenerator) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	nonce := time.Now().UnixMilli()

	// Ensure nonce is always increasing
	if nonce <= g.lastNonce {
		nonce = g.lastNonce + 1
	}

	g.lastNonce = nonce
	return strconv.FormatInt(nonce, 10)
}

// MicroTimestampNonceGenerator generates nonces based on Unix timestamp in microseconds.
// This provides higher resolution for high-frequency trading scenarios.
type MicroTimestampNonceGenerator struct {
	lastNonce int64
	mu        sync.Mutex
}

// NewMicroTimestampNonceGenerator creates a new microsecond timestamp-based nonce generator.
func NewMicroTimestampNonceGenerator() *MicroTimestampNonceGenerator {
	return &MicroTimestampNonceGenerator{}
}

// Next returns the next nonce value.
func (g *MicroTimestampNonceGenerator) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	nonce := time.Now().UnixMicro()

	// Ensure nonce is always increasing
	if nonce <= g.lastNonce {
		nonce = g.lastNonce + 1
	}

	g.lastNonce = nonce
	return strconv.FormatInt(nonce, 10)
}

// CounterNonceGenerator uses a simple atomic counter.
// This is primarily useful for testing where deterministic nonces are needed.
type CounterNonceGenerator struct {
	counter int64
}

// NewCounterNonceGenerator creates a new counter-based nonce generator.
// The start value should be set to a timestamp or other value that ensures
// the generated nonces will be greater than any previously used nonces.
func NewCounterNonceGenerator(start int64) *CounterNonceGenerator {
	return &CounterNonceGenerator{
		counter: start,
	}
}

// Next returns the next nonce value.
func (g *CounterNonceGenerator) Next() string {
	nonce := atomic.AddInt64(&g.counter, 1)
	return strconv.FormatInt(nonce, 10)
}

// FixedNonceGenerator always returns the same nonce value.
// This is ONLY for testing purposes and should never be used in production.
type FixedNonceGenerator struct {
	value string
}

// NewFixedNonceGenerator creates a new fixed nonce generator.
// WARNING: This should only be used for testing.
func NewFixedNonceGenerator(value string) *FixedNonceGenerator {
	return &FixedNonceGenerator{value: value}
}

// Next returns the fixed nonce value.
func (g *FixedNonceGenerator) Next() string {
	return g.value
}
