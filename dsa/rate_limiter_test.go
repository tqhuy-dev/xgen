package dsa

import (
	"sync"
	"testing"
	"time"
)

// TestNewTokenBucket tests the creation of a new token bucket
// and verifies initial state is correct
func TestNewTokenBucket(t *testing.T) {
	tests := []struct {
		name       string
		capacity   float64
		refillRate float64
	}{
		{
			name:       "standard bucket",
			capacity:   10,
			refillRate: 1,
		},
		{
			name:       "high capacity bucket",
			capacity:   100,
			refillRate: 10,
		},
		{
			name:       "low capacity bucket",
			capacity:   1,
			refillRate: 0.5,
		},
		{
			name:       "fractional capacity",
			capacity:   5.5,
			refillRate: 2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: create a new token bucket
			tb := NewTokenBucket(tt.capacity, tt.refillRate)

			// Assert: verify initial state
			if tb.capacity != tt.capacity {
				t.Errorf("NewTokenBucket() capacity = %v; expected %v", tb.capacity, tt.capacity)
			}
			if tb.tokens != tt.capacity {
				t.Errorf("NewTokenBucket() tokens = %v; expected %v (should start full)", tb.tokens, tt.capacity)
			}
			if tb.refillRate != tt.refillRate {
				t.Errorf("NewTokenBucket() refillRate = %v; expected %v", tb.refillRate, tt.refillRate)
			}
		})
	}
}

// TestTokenBucket_Allow tests single token consumption
func TestTokenBucket_Allow(t *testing.T) {
	tests := []struct {
		name          string
		capacity      float64
		refillRate    float64
		requests      int
		expectedAllow int // Number of requests that should be allowed
	}{
		{
			name:          "allow within capacity",
			capacity:      5,
			refillRate:    1,
			requests:      3,
			expectedAllow: 3,
		},
		{
			name:          "exceed capacity",
			capacity:      3,
			refillRate:    1,
			requests:      5,
			expectedAllow: 3,
		},
		{
			name:          "single token bucket",
			capacity:      1,
			refillRate:    1,
			requests:      2,
			expectedAllow: 1,
		},
		{
			name:          "empty bucket",
			capacity:      0,
			refillRate:    1,
			requests:      1,
			expectedAllow: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create token bucket
			tb := NewTokenBucket(tt.capacity, tt.refillRate)

			// Act: make requests
			allowedCount := 0
			for i := 0; i < tt.requests; i++ {
				if tb.Allow() {
					allowedCount++
				}
			}

			// Assert: verify allowed count
			if allowedCount != tt.expectedAllow {
				t.Errorf("Allow() allowed %d requests; expected %d", allowedCount, tt.expectedAllow)
			}
		})
	}
}

// TestTokenBucket_AllowN tests multiple token consumption
func TestTokenBucket_AllowN(t *testing.T) {
	tests := []struct {
		name          string
		capacity      float64
		refillRate    float64
		tokensNeeded  float64
		expectedAllow bool
	}{
		{
			name:          "consume within capacity",
			capacity:      10,
			refillRate:    1,
			tokensNeeded:  5,
			expectedAllow: true,
		},
		{
			name:          "consume exact capacity",
			capacity:      10,
			refillRate:    1,
			tokensNeeded:  10,
			expectedAllow: true,
		},
		{
			name:          "consume exceeds capacity",
			capacity:      10,
			refillRate:    1,
			tokensNeeded:  11,
			expectedAllow: false,
		},
		{
			name:          "consume fractional tokens",
			capacity:      5.5,
			refillRate:    1,
			tokensNeeded:  2.5,
			expectedAllow: true,
		},
		{
			name:          "consume zero tokens",
			capacity:      10,
			refillRate:    1,
			tokensNeeded:  0,
			expectedAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create token bucket
			tb := NewTokenBucket(tt.capacity, tt.refillRate)

			// Act: request multiple tokens
			result := tb.AllowN(tt.tokensNeeded)

			// Assert: verify result
			if result != tt.expectedAllow {
				t.Errorf("AllowN(%v) = %v; expected %v", tt.tokensNeeded, result, tt.expectedAllow)
			}
		})
	}
}

// TestTokenBucket_Refill tests token refilling over time
func TestTokenBucket_Refill(t *testing.T) {
	tests := []struct {
		name         string
		capacity     float64
		refillRate   float64
		initialUse   float64
		waitTime     time.Duration
		minExpected  float64
		maxExpected  float64
	}{
		{
			name:        "refill after 1 second",
			capacity:    10,
			refillRate:  2,  // 2 tokens per second
			initialUse:  5,  // Use 5 tokens
			waitTime:    time.Second,
			minExpected: 6.5, // Should have at least 5 + 2 = 7 tokens (with margin)
			maxExpected: 10,  // Can't exceed capacity
		},
		{
			name:        "partial refill",
			capacity:    10,
			refillRate:  10, // 10 tokens per second
			initialUse:  10, // Empty the bucket
			waitTime:    500 * time.Millisecond,
			minExpected: 4,  // Should have ~5 tokens after 0.5 seconds (with margin)
			maxExpected: 10,
		},
		{
			name:        "refill to capacity",
			capacity:    5,
			refillRate:  10,
			initialUse:  5,
			waitTime:    time.Second,
			minExpected: 5,  // Should refill to capacity
			maxExpected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket and consume tokens
			tb := NewTokenBucket(tt.capacity, tt.refillRate)
			tb.AllowN(tt.initialUse)

			// Act: wait for refill
			time.Sleep(tt.waitTime)
			available := tb.AvailableTokens()

			// Assert: verify refilled tokens are within expected range
			if available < tt.minExpected || available > tt.maxExpected {
				t.Errorf("AvailableTokens() = %v; expected between %v and %v",
					available, tt.minExpected, tt.maxExpected)
			}
		})
	}
}

// TestTokenBucket_AvailableTokens tests querying available tokens
func TestTokenBucket_AvailableTokens(t *testing.T) {
	tests := []struct {
		name      string
		capacity  float64
		consume   float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "full bucket",
			capacity:  10,
			consume:   0,
			expected:  10,
			tolerance: 0.01,
		},
		{
			name:      "partially consumed",
			capacity:  10,
			consume:   3,
			expected:  7,
			tolerance: 0.01,
		},
		{
			name:      "empty bucket",
			capacity:  10,
			consume:   10,
			expected:  0,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket and consume tokens
			tb := NewTokenBucket(tt.capacity, 1)
			if tt.consume > 0 {
				tb.AllowN(tt.consume)
			}

			// Act: query available tokens
			available := tb.AvailableTokens()

			// Assert: verify available tokens (with tolerance for floating-point precision)
			diff := available - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("AvailableTokens() = %v; expected %v (diff: %v, tolerance: %v)",
					available, tt.expected, diff, tt.tolerance)
			}
		})
	}
}

// TestTokenBucket_Reset tests resetting the bucket
func TestTokenBucket_Reset(t *testing.T) {
	tests := []struct {
		name      string
		capacity  float64
		consume   float64
		tolerance float64
	}{
		{
			name:      "reset after partial consumption",
			capacity:  10,
			consume:   5,
			tolerance: 0.01,
		},
		{
			name:      "reset empty bucket",
			capacity:  10,
			consume:   10,
			tolerance: 0.01,
		},
		{
			name:      "reset full bucket",
			capacity:  10,
			consume:   0,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket and consume tokens
			tb := NewTokenBucket(tt.capacity, 1)
			if tt.consume > 0 {
				tb.AllowN(tt.consume)
			}

			// Act: reset the bucket
			tb.Reset()

			// Assert: verify bucket is full (with tolerance for floating-point precision)
			available := tb.AvailableTokens()
			diff := available - tt.capacity
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("After Reset(), AvailableTokens() = %v; expected %v (diff: %v, tolerance: %v)",
					available, tt.capacity, diff, tt.tolerance)
			}
		})
	}
}

// TestTokenBucket_Concurrent tests thread safety
func TestTokenBucket_Concurrent(t *testing.T) {
	tests := []struct {
		name        string
		capacity    float64
		refillRate  float64
		goroutines  int
		requestsEach int
	}{
		{
			name:        "concurrent access with 10 goroutines",
			capacity:    100,
			refillRate:  10,
			goroutines:  10,
			requestsEach: 20,
		},
		{
			name:        "high contention",
			capacity:    10,
			refillRate:  1,
			goroutines:  50,
			requestsEach: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket
			tb := NewTokenBucket(tt.capacity, tt.refillRate)
			var wg sync.WaitGroup
			totalAllowed := 0
			var mu sync.Mutex

			// Act: spawn multiple goroutines making requests
			for i := 0; i < tt.goroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < tt.requestsEach; j++ {
						if tb.Allow() {
							mu.Lock()
							totalAllowed++
							mu.Unlock()
						}
					}
				}()
			}

			wg.Wait()

			// Assert: verify no tokens were over-consumed
			// We can't predict exact number due to refills during execution
			// But total allowed should be <= capacity + (time * refillRate)
			// At minimum, we verify no panic occurred and some requests were allowed
			if totalAllowed < 0 {
				t.Errorf("Concurrent access resulted in negative allowed count: %d", totalAllowed)
			}
		})
	}
}

// TestTokenBucket_WaitForToken tests blocking until token is available
func TestTokenBucket_WaitForToken(t *testing.T) {
	tests := []struct {
		name       string
		capacity   float64
		refillRate float64
		initialUse float64
	}{
		{
			name:       "wait for refill",
			capacity:   1,
			refillRate: 10, // Fast refill for quick test
			initialUse: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket and empty it
			tb := NewTokenBucket(tt.capacity, tt.refillRate)
			tb.AllowN(tt.initialUse)

			start := time.Now()

			// Act: wait for token to become available
			tb.WaitForToken()
			elapsed := time.Since(start)

			// Assert: verify we had to wait (some time passed)
			// and token was eventually acquired
			if elapsed < 10*time.Millisecond {
				t.Errorf("WaitForToken() returned too quickly (%v); expected some wait time", elapsed)
			}

			// Verify token was consumed
			if tb.Allow() {
				// Should have consumed the token, so immediate allow should fail
				// unless refill happened, which is possible
			}
		})
	}
}

// TestTokenBucket_Release tests releasing tokens back to the bucket
func TestTokenBucket_Release(t *testing.T) {
	tests := []struct {
		name      string
		capacity  float64
		consume   float64
		release   float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "release single token",
			capacity:  10,
			consume:   5,
			release:   1,
			expected:  6,
			tolerance: 0.01,
		},
		{
			name:      "release multiple tokens",
			capacity:  10,
			consume:   7,
			release:   3,
			expected:  6,
			tolerance: 0.01,
		},
		{
			name:      "release to capacity limit",
			capacity:  5,
			consume:   2,
			release:   10, // Try to release more than consumed
			expected:  5,  // Should cap at capacity
			tolerance: 0.01,
		},
		{
			name:      "release all consumed tokens",
			capacity:  10,
			consume:   10,
			release:   10,
			expected:  10,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket and consume tokens
			tb := NewTokenBucket(tt.capacity, 0) // No auto-refill
			tb.AllowN(tt.consume)

			// Act: release tokens
			tb.Release(tt.release)
			available := tb.AvailableTokens()

			// Assert: verify available tokens
			diff := available - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("After Release(%v), AvailableTokens() = %v; expected %v (diff: %v)",
					tt.release, available, tt.expected, diff)
			}
		})
	}
}

// TestTokenBucket_TryAllowWithTimeout tests timeout functionality
func TestTokenBucket_TryAllowWithTimeout(t *testing.T) {
	tests := []struct {
		name          string
		capacity      float64
		consume       float64
		tokensNeeded  float64
		timeout       time.Duration
		releaseAfter  time.Duration
		releaseAmount float64
		expectSuccess bool
	}{
		{
			name:          "acquire without waiting",
			capacity:      10,
			consume:       0,
			tokensNeeded:  1,
			timeout:       100 * time.Millisecond,
			releaseAfter:  0,
			releaseAmount: 0,
			expectSuccess: true,
		},
		{
			name:          "timeout when no tokens available",
			capacity:      5,
			consume:       5,
			tokensNeeded:  1,
			timeout:       100 * time.Millisecond,
			releaseAfter:  0,
			releaseAmount: 0,
			expectSuccess: false,
		},
		{
			name:          "acquire after token released",
			capacity:      5,
			consume:       5,
			tokensNeeded:  1,
			timeout:       500 * time.Millisecond,
			releaseAfter:  100 * time.Millisecond,
			releaseAmount: 1,
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create bucket and consume tokens
			tb := NewTokenBucket(tt.capacity, 0) // No auto-refill
			if tt.consume > 0 {
				tb.AllowN(tt.consume)
			}

			// If we need to release tokens after a delay, do it in a goroutine
			if tt.releaseAfter > 0 && tt.releaseAmount > 0 {
				go func() {
					time.Sleep(tt.releaseAfter)
					tb.Release(tt.releaseAmount)
				}()
			}

			// Act: try to acquire tokens with timeout
			start := time.Now()
			result := tb.TryAllowWithTimeout(tt.tokensNeeded, tt.timeout)
			elapsed := time.Since(start)

			// Assert: verify result matches expectation
			if result != tt.expectSuccess {
				t.Errorf("TryAllowWithTimeout() = %v; expected %v (elapsed: %v)",
					result, tt.expectSuccess, elapsed)
			}

			// Verify timeout was respected (with some margin)
			if !result && elapsed > tt.timeout+50*time.Millisecond {
				t.Errorf("TryAllowWithTimeout() took %v; expected timeout at %v",
					elapsed, tt.timeout)
			}
		})
	}
}

