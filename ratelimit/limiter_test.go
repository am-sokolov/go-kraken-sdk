package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestTokenBucket_Allow(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Should allow initial requests
	for i := 0; i < 50; i++ {
		if !tb.Allow("account") {
			t.Errorf("request %d should be allowed", i)
		}
	}
}

func TestTokenBucket_BlockWhenFull(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Fill the bucket
	for i := 0; i < 60; i++ {
		tb.Allow("account")
	}

	// Next request should be blocked
	if tb.Allow("account") {
		t.Error("request should be blocked when bucket is full")
	}
}

func TestTokenBucket_Decay(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Fill the bucket
	for i := 0; i < 60; i++ {
		tb.Allow("account")
	}

	// Wait for some decay
	time.Sleep(2 * time.Second)

	// Should now have room
	if !tb.Allow("account") {
		t.Error("request should be allowed after decay")
	}
}

func TestTokenBucket_Wait(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Fill the bucket
	for i := 0; i < 60; i++ {
		tb.Allow("account")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	err := tb.Wait(ctx, "account")
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Wait returned error: %v", err)
	}

	// Should have waited about 1 second (1 unit decay per second)
	if elapsed < 500*time.Millisecond {
		t.Errorf("expected to wait, but elapsed only %v", elapsed)
	}
}

func TestTokenBucket_WaitContextCancelled(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Fill the bucket
	for i := 0; i < 60; i++ {
		tb.Allow("account")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := tb.Wait(ctx, "account")
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestTokenBucket_PublicEndpointsFree(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Public endpoints should not increment counter
	for i := 0; i < 1000; i++ {
		if !tb.Allow("public") {
			t.Errorf("public request %d should always be allowed", i)
		}
	}

	if tb.Counter() != 0 {
		t.Errorf("counter should be 0 for public endpoints, got %f", tb.Counter())
	}
}

func TestTokenBucket_LedgerEndpointsCostMore(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Ledger endpoints cost 2
	tb.Allow("ledger")

	counter := tb.Counter()
	// Use approximate comparison for floating point
	if counter < 1.99 || counter > 2.01 {
		t.Errorf("expected counter ~2 for ledger endpoint, got %f", counter)
	}
}

func TestTokenBucket_Reset(t *testing.T) {
	tb := NewTokenBucket(TierStarter)

	// Fill the bucket
	for i := 0; i < 30; i++ {
		tb.Allow("account")
	}

	// Reset
	tb.Reset()

	if tb.Counter() != 0 {
		t.Errorf("counter should be 0 after reset, got %f", tb.Counter())
	}
}

func TestPerPairLimiter_Allow(t *testing.T) {
	ppl := NewPerPairLimiter(TierStarter)

	// Should allow initial requests
	for i := 0; i < 50; i++ {
		if !ppl.AllowForPair("XBTUSD", 1) {
			t.Errorf("request %d should be allowed", i)
		}
	}
}

func TestPerPairLimiter_SeparatePairs(t *testing.T) {
	ppl := NewPerPairLimiter(TierStarter)

	// Fill one pair
	for i := 0; i < 60; i++ {
		ppl.AllowForPair("XBTUSD", 1)
	}

	// Other pair should still be allowed
	if !ppl.AllowForPair("ETHUSD", 1) {
		t.Error("ETHUSD should be allowed when XBTUSD is full")
	}
}

func TestPerPairLimiter_BlockWhenFull(t *testing.T) {
	ppl := NewPerPairLimiter(TierStarter)

	// Fill the pair's limit
	for i := 0; i < 60; i++ {
		ppl.AllowForPair("XBTUSD", 1)
	}

	// Next request should be blocked
	if ppl.AllowForPair("XBTUSD", 1) {
		t.Error("request should be blocked when pair limit is full")
	}
}

func TestCombinedLimiter(t *testing.T) {
	cl := NewCombinedLimiter(TierStarter)

	// API requests should work
	if !cl.AllowAPI("account") {
		t.Error("API request should be allowed")
	}

	// Match engine requests should work
	if !cl.AllowMatchEngine("XBTUSD", 1) {
		t.Error("match engine request should be allowed")
	}
}

func TestNoopLimiter(t *testing.T) {
	nl := NewNoopLimiter()

	// Should always allow
	for i := 0; i < 1000; i++ {
		if !nl.Allow("account") {
			t.Errorf("noop limiter should always allow")
		}
	}

	// Should never wait
	ctx := context.Background()
	if err := nl.Wait(ctx, "account"); err != nil {
		t.Errorf("noop limiter Wait should not error: %v", err)
	}

	// Reserve should return 0
	if d := nl.Reserve("account"); d != 0 {
		t.Errorf("noop limiter Reserve should return 0, got %v", d)
	}
}

func TestTierConfigurations(t *testing.T) {
	tests := []struct {
		tier               Tier
		expectedMax        float64
		expectedDecay      float64
		expectedBatchLimit int
	}{
		{TierStarter, 60, 1.0, 60},
		{TierIntermediate, 125, 2.34, 80},
		{TierPro, 180, 3.75, 225},
	}

	for _, tt := range tests {
		config := DefaultTierConfigs[tt.tier]

		if config.MaxCounter != tt.expectedMax {
			t.Errorf("Tier %d: expected MaxCounter %f, got %f", tt.tier, tt.expectedMax, config.MaxCounter)
		}

		if config.DecayPerSecond != tt.expectedDecay {
			t.Errorf("Tier %d: expected DecayPerSecond %f, got %f", tt.tier, tt.expectedDecay, config.DecayPerSecond)
		}

		if config.MaxOrdersPerBatch != tt.expectedBatchLimit {
			t.Errorf("Tier %d: expected MaxOrdersPerBatch %d, got %d", tt.tier, tt.expectedBatchLimit, config.MaxOrdersPerBatch)
		}
	}
}
