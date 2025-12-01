package dsa

import (
	"sync"
	"time"
)

// TokenBucket implements a token bucket rate limiter algorithm
// It allows for controlled bursts while maintaining an average rate
type TokenBucket struct {
	capacity     float64       // Maximum number of tokens in the bucket
	tokens       float64       // Current number of tokens
	refillRate   float64       // Tokens added per second
	lastRefillAt time.Time     // Last time tokens were refilled
	mu           sync.Mutex    // Mutex for thread-safe operations
}

// NewTokenBucket creates a new token bucket rate limiter
// capacity: maximum number of tokens (burst size)
// refillRate: number of tokens added per second
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity, // Start with full bucket
		refillRate:   refillRate,
		lastRefillAt: time.Now(),
	}
}

// Allow checks if a request can be processed (consumes 1 token)
// Returns true if the request is allowed, false otherwise
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN checks if n requests can be processed (consumes n tokens)
// Returns true if the requests are allowed, false otherwise
func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

// refill adds tokens based on elapsed time since last refill
// This method should be called with the mutex locked
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillAt).Seconds()
	
	// Calculate tokens to add based on elapsed time and refill rate
	tokensToAdd := elapsed * tb.refillRate
	
	// Add tokens but don't exceed capacity
	tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
	tb.lastRefillAt = now
}

// AvailableTokens returns the current number of available tokens
func (tb *TokenBucket) AvailableTokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}

// Reset resets the bucket to full capacity
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens = tb.capacity
	tb.lastRefillAt = time.Now()
}

// Release returns n tokens back to the bucket
// This is useful when a request completes and you want to free up capacity
func (tb *TokenBucket) Release(n float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Add tokens back but don't exceed capacity
	tb.tokens = min(tb.capacity, tb.tokens+n)
}

// TryAllowWithTimeout attempts to acquire n tokens with a timeout
// Returns true if tokens were acquired, false if timeout occurred
func (tb *TokenBucket) TryAllowWithTimeout(n float64, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	
	for {
		if tb.AllowN(n) {
			return true
		}
		
		// Check if timeout reached
		if time.Now().After(deadline) {
			return false
		}
		
		// Sleep for a short time before checking again
		time.Sleep(10 * time.Millisecond)
	}
}

// WaitForToken blocks until a token is available
// Returns when a token can be consumed
func (tb *TokenBucket) WaitForToken() {
	tb.WaitForTokens(1)
}

// WaitForTokens blocks until n tokens are available
// Returns when n tokens can be consumed
func (tb *TokenBucket) WaitForTokens(n float64) {
	for {
		if tb.AllowN(n) {
			return
		}
		// Sleep for a short time before checking again
		time.Sleep(10 * time.Millisecond)
	}
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

