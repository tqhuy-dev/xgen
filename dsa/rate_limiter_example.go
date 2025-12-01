package dsa

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/tqhuy-dev/xgen/utilities"
)

// basicRateLimiterExample demonstrates basic rate limiting
func basicRateLimiterExample() {
	// Create a rate limiter: 5 tokens capacity, refills at 1 token/second
	rateLimiter := NewTokenBucket(5, 1)

	// Try to process 8 requests
	for i := 1; i <= 8; i++ {
		if rateLimiter.Allow() {
			fmt.Printf("Request %d: ✓ Allowed\n", i)
		} else {
			fmt.Printf("Request %d: ✗ Rejected (rate limit exceeded)\n", i)
		}
	}
}

// burstHandlingExample shows how token bucket handles burst traffic
func burstHandlingExample() {
	// Create a rate limiter: 3 tokens capacity, refills at 2 tokens/second
	rateLimiter := NewTokenBucket(3, 2)

	fmt.Println("Initial burst of 3 requests:")
	for i := 1; i <= 3; i++ {
		if rateLimiter.Allow() {
			fmt.Printf("  Request %d: ✓ Allowed (tokens left: %.1f)\n", i, rateLimiter.AvailableTokens())
		}
	}

	fmt.Println("\nImmediate 4th request (should fail):")
	if rateLimiter.Allow() {
		fmt.Println("  Request 4: ✓ Allowed")
	} else {
		fmt.Printf("  Request 4: ✗ Rejected (tokens left: %.1f)\n", rateLimiter.AvailableTokens())
	}
}

// refillExample demonstrates token refilling over time
func refillExample() {
	// Create a rate limiter: 10 tokens capacity, refills at 5 tokens/second
	rateLimiter := NewTokenBucket(10, 5)

	// Consume all tokens
	for i := 0; i < 10; i++ {
		rateLimiter.Allow()
	}
	fmt.Printf("After consuming all tokens: %.1f tokens available\n", rateLimiter.AvailableTokens())

	// Wait 1 second for refill
	fmt.Println("Waiting 1 second for refill...")
	time.Sleep(1 * time.Second)
	fmt.Printf("After 1 second: %.1f tokens available\n", rateLimiter.AvailableTokens())

	// Wait another second
	fmt.Println("Waiting another second...")
	time.Sleep(1 * time.Second)
	fmt.Printf("After 2 seconds total: %.1f tokens available\n", rateLimiter.AvailableTokens())
}

// apiRateLimitExample simulates API rate limiting with UUID generation
func apiRateLimitExample() {
	// Simulate an API with rate limit: 3 requests per second
	rateLimiter := NewTokenBucket(3, 3)

	fmt.Println("Simulating API calls (3 requests/second limit):")

	// Try to make 10 API calls rapidly
	for i := 1; i <= 10; i++ {
		if rateLimiter.Allow() {
			uuid := utilities.GenerateUUIDV7()
			fmt.Printf("API Call %d: ✓ Success - Generated UUID: %s\n", i, uuid)
		} else {
			fmt.Printf("API Call %d: ✗ Rate limited - Please try again later\n", i)
		}

		// Small delay between requests
		time.Sleep(200 * time.Millisecond)
	}

	// Reset the rate limiter
	fmt.Println("\nResetting rate limiter...")
	rateLimiter.Reset()
	fmt.Printf("After reset: %.1f tokens available\n", rateLimiter.AvailableTokens())
}

// workerPoolExample demonstrates a worker pool with capacity 5
// Requests must wait if all 5 slots are busy, with a 2s timeout
func workerPoolExample() {
	// Create a rate limiter with capacity 5, no auto-refill (refill rate = 0)
	// This ensures tokens are only returned when workers complete
	workerPool := NewTokenBucket(5, 0)

	totalRequests := 12
	timeout := 2 * time.Second

	fmt.Printf("Worker Pool: Capacity = 5, Total Requests = %d, Timeout = %v\n\n", totalRequests, timeout)

	// Track statistics
	var successCount, timeoutCount int
	var mu sync.Mutex

	// Launch all requests concurrently
	var wg sync.WaitGroup
	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)
		requestID := i

		go func() {
			defer wg.Done()

			startTime := time.Now()
			fmt.Printf("[Request %2d] Attempting to acquire worker slot... (Available: %.0f/5)\n",
				requestID, workerPool.AvailableTokens())

			// Try to acquire a worker slot with timeout
			if !workerPool.TryAllowWithTimeout(1, timeout) {
				elapsed := time.Since(startTime)
				fmt.Printf("[Request %2d] ✗ TIMEOUT after %.2fs - All workers busy\n", requestID, elapsed.Seconds())
				mu.Lock()
				timeoutCount++
				mu.Unlock()
				return
			}

			acquireTime := time.Since(startTime)
			if acquireTime > 10*time.Millisecond {
				fmt.Printf("[Request %2d] ✓ Acquired worker slot (waited %.2fs)\n", requestID, acquireTime.Seconds())
			} else {
				fmt.Printf("[Request %2d] ✓ Acquired worker slot immediately\n", requestID)
			}

			// Simulate random processing time (100ms to 800ms)
			processingTime := time.Duration(100+rand.Intn(700)) * time.Millisecond

			// Process the request
			time.Sleep(processingTime)
			uuid := utilities.GenerateUUIDV7()

			fmt.Printf("[Request %2d] ✓ Completed in %.2fs - UUID: %s\n",
				requestID, processingTime.Seconds(), uuid)

			// Release the worker slot
			workerPool.Release(1)
			fmt.Printf("[Request %2d] Released worker slot (Available: %.0f/5)\n",
				requestID, workerPool.AvailableTokens())

			mu.Lock()
			successCount++
			mu.Unlock()
		}()

		// Stagger the request launches slightly to show the queuing behavior
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for all requests to complete
	wg.Wait()

	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("Summary:\n")
	fmt.Printf("  Total Requests:     %d\n", totalRequests)
	fmt.Printf("  Successful:         %d (%.1f%%)\n", successCount, float64(successCount)/float64(totalRequests)*100)
	fmt.Printf("  Timed Out:          %d (%.1f%%)\n", timeoutCount, float64(timeoutCount)/float64(totalRequests)*100)
	fmt.Printf("  Available Workers:  %.0f/5\n", workerPool.AvailableTokens())
	fmt.Printf(strings.Repeat("=", 60) + "\n")
}
