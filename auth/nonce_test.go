package auth

import (
	"strconv"
	"sync"
	"testing"
)

func TestTimestampNonceGenerator_Next(t *testing.T) {
	gen := NewTimestampNonceGenerator()

	// Get first nonce
	nonce1 := gen.Next()
	if nonce1 == "" {
		t.Error("expected non-empty nonce")
	}

	// Parse as int64
	_, err := strconv.ParseInt(nonce1, 10, 64)
	if err != nil {
		t.Errorf("nonce should be valid int64: %v", err)
	}

	// Get second nonce - should be greater
	nonce2 := gen.Next()
	n1, _ := strconv.ParseInt(nonce1, 10, 64)
	n2, _ := strconv.ParseInt(nonce2, 10, 64)

	if n2 <= n1 {
		t.Errorf("second nonce (%d) should be greater than first (%d)", n2, n1)
	}
}

func TestTimestampNonceGenerator_Concurrent(t *testing.T) {
	gen := NewTimestampNonceGenerator()
	const numGoroutines = 100
	const numNonces = 100

	nonces := make(chan string, numGoroutines*numNonces)
	var wg sync.WaitGroup

	// Generate nonces concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numNonces; j++ {
				nonces <- gen.Next()
			}
		}()
	}

	wg.Wait()
	close(nonces)

	// Collect all nonces
	seen := make(map[string]bool)
	for nonce := range nonces {
		if seen[nonce] {
			t.Errorf("duplicate nonce: %s", nonce)
		}
		seen[nonce] = true
	}

	// Should have all unique nonces
	if len(seen) != numGoroutines*numNonces {
		t.Errorf("expected %d unique nonces, got %d", numGoroutines*numNonces, len(seen))
	}
}

func TestMicroTimestampNonceGenerator_Next(t *testing.T) {
	gen := NewMicroTimestampNonceGenerator()

	nonce1 := gen.Next()
	nonce2 := gen.Next()

	n1, _ := strconv.ParseInt(nonce1, 10, 64)
	n2, _ := strconv.ParseInt(nonce2, 10, 64)

	if n2 <= n1 {
		t.Errorf("second nonce (%d) should be greater than first (%d)", n2, n1)
	}
}

func TestCounterNonceGenerator_Next(t *testing.T) {
	start := int64(1000)
	gen := NewCounterNonceGenerator(start)

	nonce1 := gen.Next()
	nonce2 := gen.Next()
	nonce3 := gen.Next()

	if nonce1 != "1001" {
		t.Errorf("first nonce = %s, want 1001", nonce1)
	}
	if nonce2 != "1002" {
		t.Errorf("second nonce = %s, want 1002", nonce2)
	}
	if nonce3 != "1003" {
		t.Errorf("third nonce = %s, want 1003", nonce3)
	}
}

func TestCounterNonceGenerator_Concurrent(t *testing.T) {
	gen := NewCounterNonceGenerator(0)
	const numGoroutines = 100
	const numNonces = 100

	nonces := make(chan string, numGoroutines*numNonces)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numNonces; j++ {
				nonces <- gen.Next()
			}
		}()
	}

	wg.Wait()
	close(nonces)

	// Check all unique
	seen := make(map[string]bool)
	for nonce := range nonces {
		if seen[nonce] {
			t.Errorf("duplicate nonce: %s", nonce)
		}
		seen[nonce] = true
	}
}

func TestFixedNonceGenerator_Next(t *testing.T) {
	value := "fixed-nonce-value"
	gen := NewFixedNonceGenerator(value)

	for i := 0; i < 10; i++ {
		nonce := gen.Next()
		if nonce != value {
			t.Errorf("Next() = %s, want %s", nonce, value)
		}
	}
}

func BenchmarkTimestampNonceGenerator_Next(b *testing.B) {
	gen := NewTimestampNonceGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gen.Next()
	}
}

func BenchmarkMicroTimestampNonceGenerator_Next(b *testing.B) {
	gen := NewMicroTimestampNonceGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gen.Next()
	}
}

func BenchmarkCounterNonceGenerator_Next(b *testing.B) {
	gen := NewCounterNonceGenerator(0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gen.Next()
	}
}
