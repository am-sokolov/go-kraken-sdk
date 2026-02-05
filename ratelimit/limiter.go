// Package ratelimit provides rate limiting for Kraken API requests.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Tier represents the Kraken account tier for rate limiting.
// Current Kraken documentation (as of 2024) defines two tiers:
//   - Standard verification (formerly Intermediate): MaxCounter=125, Decay=2.34/sec
//   - Verified with higher limits (formerly Pro): MaxCounter=180, Decay=3.75/sec
//
// TierStarter is provided for backwards compatibility but may not reflect
// current Kraken rate limits. Use TierIntermediate or TierPro for production.
type Tier int

const (
	// TierStarter is a legacy tier with lowest limits.
	// Deprecated: This tier may not be supported by current Kraken API.
	// Consider using TierIntermediate instead.
	TierStarter Tier = iota
	// TierIntermediate is the standard verification tier (formerly "Intermediate").
	// MaxCounter: 125, DecayPerSecond: 2.34
	TierIntermediate
	// TierPro is the verified with higher limits tier (formerly "Pro").
	// MaxCounter: 180, DecayPerSecond: 3.75
	TierPro
)

// TierConfig contains rate limit configuration for a tier.
type TierConfig struct {
	// MaxCounter is the maximum counter value before requests are blocked.
	MaxCounter float64
	// DecayPerSecond is how much the counter decays per second.
	DecayPerSecond float64
	// MaxOrdersPerBatch is the maximum orders per batch request.
	MaxOrdersPerBatch int
}

// DefaultTierConfigs contains the default tier configurations based on Kraken docs.
// See: https://support.kraken.com/hc/en-us/articles/360045239571-Trading-rate-limits
var DefaultTierConfigs = map[Tier]TierConfig{
	// TierStarter is a legacy configuration - may not reflect current API limits.
	TierStarter: {
		MaxCounter:        60,
		DecayPerSecond:    1.0,
		MaxOrdersPerBatch: 60,
	},
	TierIntermediate: {
		MaxCounter:        125,
		DecayPerSecond:    2.34,
		MaxOrdersPerBatch: 80,
	},
	TierPro: {
		MaxCounter:        180,
		DecayPerSecond:    3.75,
		MaxOrdersPerBatch: 225,
	},
}

// EndpointCosts defines the counter cost for different endpoint types.
var EndpointCosts = map[string]float64{
	"public":  0, // Public endpoints don't count
	"account": 1, // Balance, TradeBalance, etc.
	"trade":   0, // AddOrder, CancelOrder (uses separate matching engine limits)
	"ledger":  2, // Ledgers, QueryLedgers
	"export":  2, // RequestExport, ExportStatus
	"other":   1, // Default cost
}

// MatchingEngineCosts defines costs for matching engine operations.
var MatchingEngineCosts = map[string]int{
	"add_order":    1,
	"cancel_order": 1,
	"edit_order":   2, // cancel + add
	"batch_add":    1, // per order
	"batch_cancel": 1, // per order
	"cancel_all":   1,
}

// Limiter provides rate limiting for API requests.
type Limiter interface {
	// Wait blocks until the request can proceed.
	Wait(ctx context.Context, endpoint string) error
	// Allow returns immediately whether the request can proceed.
	Allow(endpoint string) bool
	// Reserve reserves capacity for a request, returns wait time.
	Reserve(endpoint string) time.Duration
}

// TokenBucket implements a token bucket rate limiter.
type TokenBucket struct {
	tier       Tier
	config     TierConfig
	counter    float64
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket rate limiter for the given tier.
func NewTokenBucket(tier Tier) *TokenBucket {
	config, ok := DefaultTierConfigs[tier]
	if !ok {
		config = DefaultTierConfigs[TierStarter]
	}

	return &TokenBucket{
		tier:       tier,
		config:     config,
		counter:    0,
		lastUpdate: time.Now(),
	}
}

// Wait blocks until the request can proceed or context is cancelled.
func (tb *TokenBucket) Wait(ctx context.Context, endpoint string) error {
	waitDuration := tb.Reserve(endpoint)
	if waitDuration == 0 {
		return nil
	}

	timer := time.NewTimer(waitDuration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// Allow returns immediately whether the request can proceed.
func (tb *TokenBucket) Allow(endpoint string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.decay()

	cost := tb.getCost(endpoint)
	if tb.counter+cost > tb.config.MaxCounter {
		return false
	}

	tb.counter += cost
	return true
}

// Reserve reserves capacity for a request and returns the wait time.
func (tb *TokenBucket) Reserve(endpoint string) time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.decay()

	cost := tb.getCost(endpoint)
	if tb.counter+cost <= tb.config.MaxCounter {
		tb.counter += cost
		return 0
	}

	// Calculate wait time needed
	excess := (tb.counter + cost) - tb.config.MaxCounter
	waitSeconds := excess / tb.config.DecayPerSecond
	waitDuration := time.Duration(waitSeconds * float64(time.Second))

	// Reserve the capacity
	tb.counter += cost

	return waitDuration
}

// decay updates the counter based on time elapsed.
func (tb *TokenBucket) decay() {
	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.lastUpdate = now

	tb.counter -= elapsed * tb.config.DecayPerSecond
	if tb.counter < 0 {
		tb.counter = 0
	}
}

// getCost returns the cost for the given endpoint type.
func (tb *TokenBucket) getCost(endpoint string) float64 {
	if cost, ok := EndpointCosts[endpoint]; ok {
		return cost
	}
	return EndpointCosts["other"]
}

// Counter returns the current counter value.
func (tb *TokenBucket) Counter() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.decay()
	return tb.counter
}

// Reset resets the counter to zero.
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.counter = 0
	tb.lastUpdate = time.Now()
}

// PerPairLimiter provides per-pair rate limiting for matching engine operations.
type PerPairLimiter struct {
	tier     Tier
	config   TierConfig
	counters map[string]*pairCounter
	mu       sync.RWMutex
}

type pairCounter struct {
	count      float64
	lastUpdate time.Time
}

// NewPerPairLimiter creates a new per-pair rate limiter.
func NewPerPairLimiter(tier Tier) *PerPairLimiter {
	config, ok := DefaultTierConfigs[tier]
	if !ok {
		config = DefaultTierConfigs[TierStarter]
	}

	return &PerPairLimiter{
		tier:     tier,
		config:   config,
		counters: make(map[string]*pairCounter),
	}
}

// WaitForPair blocks until a matching engine operation can proceed for the pair.
func (ppl *PerPairLimiter) WaitForPair(ctx context.Context, pair string, cost int) error {
	waitDuration := ppl.ReserveForPair(pair, cost)
	if waitDuration == 0 {
		return nil
	}

	timer := time.NewTimer(waitDuration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// AllowForPair returns whether a matching engine operation can proceed for the pair.
func (ppl *PerPairLimiter) AllowForPair(pair string, cost int) bool {
	ppl.mu.Lock()
	defer ppl.mu.Unlock()

	pc := ppl.getOrCreateCounter(pair)
	ppl.decayCounter(pc)

	if pc.count+float64(cost) > float64(ppl.config.MaxOrdersPerBatch) {
		return false
	}

	pc.count += float64(cost)
	return true
}

// ReserveForPair reserves capacity for a matching engine operation.
func (ppl *PerPairLimiter) ReserveForPair(pair string, cost int) time.Duration {
	ppl.mu.Lock()
	defer ppl.mu.Unlock()

	pc := ppl.getOrCreateCounter(pair)
	ppl.decayCounter(pc)

	costFloat := float64(cost)
	maxFloat := float64(ppl.config.MaxOrdersPerBatch)

	if pc.count+costFloat <= maxFloat {
		pc.count += costFloat
		return 0
	}

	// Calculate wait time needed
	excess := (pc.count + costFloat) - maxFloat
	waitSeconds := excess / ppl.config.DecayPerSecond
	waitDuration := time.Duration(waitSeconds * float64(time.Second))

	// Reserve the capacity
	pc.count += costFloat

	return waitDuration
}

func (ppl *PerPairLimiter) getOrCreateCounter(pair string) *pairCounter {
	pc, ok := ppl.counters[pair]
	if !ok {
		pc = &pairCounter{
			count:      0,
			lastUpdate: time.Now(),
		}
		ppl.counters[pair] = pc
	}
	return pc
}

func (ppl *PerPairLimiter) decayCounter(pc *pairCounter) {
	now := time.Now()
	elapsed := now.Sub(pc.lastUpdate).Seconds()
	pc.lastUpdate = now

	pc.count -= elapsed * ppl.config.DecayPerSecond
	if pc.count < 0 {
		pc.count = 0
	}
}

// CleanupStaleCounters removes counters that have been idle for longer than maxIdleTime.
// This prevents unbounded memory growth when trading many different pairs.
func (ppl *PerPairLimiter) CleanupStaleCounters(maxIdleTime time.Duration) int {
	ppl.mu.Lock()
	defer ppl.mu.Unlock()

	now := time.Now()
	removed := 0

	for pair, pc := range ppl.counters {
		if now.Sub(pc.lastUpdate) > maxIdleTime {
			delete(ppl.counters, pair)
			removed++
		}
	}

	return removed
}

// CounterCount returns the number of pair counters currently tracked.
func (ppl *PerPairLimiter) CounterCount() int {
	ppl.mu.RLock()
	defer ppl.mu.RUnlock()
	return len(ppl.counters)
}

// CombinedLimiter combines API rate limiting and matching engine limits.
type CombinedLimiter struct {
	api         *TokenBucket
	matchEngine *PerPairLimiter
}

// NewCombinedLimiter creates a new combined rate limiter.
func NewCombinedLimiter(tier Tier) *CombinedLimiter {
	return &CombinedLimiter{
		api:         NewTokenBucket(tier),
		matchEngine: NewPerPairLimiter(tier),
	}
}

// WaitAPI waits for API rate limit.
func (cl *CombinedLimiter) WaitAPI(ctx context.Context, endpoint string) error {
	return cl.api.Wait(ctx, endpoint)
}

// WaitMatchEngine waits for matching engine rate limit.
func (cl *CombinedLimiter) WaitMatchEngine(ctx context.Context, pair string, cost int) error {
	return cl.matchEngine.WaitForPair(ctx, pair, cost)
}

// AllowAPI returns whether an API request can proceed.
func (cl *CombinedLimiter) AllowAPI(endpoint string) bool {
	return cl.api.Allow(endpoint)
}

// AllowMatchEngine returns whether a matching engine operation can proceed.
func (cl *CombinedLimiter) AllowMatchEngine(pair string, cost int) bool {
	return cl.matchEngine.AllowForPair(pair, cost)
}

// Reset resets all rate limiters.
func (cl *CombinedLimiter) Reset() {
	cl.api.Reset()
	// Clear per-pair counters
	cl.matchEngine.mu.Lock()
	cl.matchEngine.counters = make(map[string]*pairCounter)
	cl.matchEngine.mu.Unlock()
}

// NoopLimiter is a rate limiter that allows all requests.
type NoopLimiter struct{}

// NewNoopLimiter creates a new noop limiter.
func NewNoopLimiter() *NoopLimiter {
	return &NoopLimiter{}
}

// Wait always returns immediately.
func (nl *NoopLimiter) Wait(ctx context.Context, endpoint string) error {
	return nil
}

// Allow always returns true.
func (nl *NoopLimiter) Allow(endpoint string) bool {
	return true
}

// Reserve always returns zero wait time.
func (nl *NoopLimiter) Reserve(endpoint string) time.Duration {
	return 0
}
