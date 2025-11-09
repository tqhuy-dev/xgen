package dsa

import (
	"fmt"
	"sync"
	"testing"
)

// TestNewRoundRobinBasic tests the NewRoundRobinBasic constructor
func TestNewRoundRobinBasic(t *testing.T) {
	tests := []struct {
		name         string
		nodes        []int
		expectedLen  int
		expectNonNil bool
	}{
		{
			name:         "create with single node",
			nodes:        []int{1},
			expectedLen:  1,
			expectNonNil: true,
		},
		{
			name:         "create with multiple nodes",
			nodes:        []int{1, 2, 3},
			expectedLen:  3,
			expectNonNil: true,
		},
		{
			name:         "create with empty nodes",
			nodes:        []int{},
			expectedLen:  0,
			expectNonNil: true,
		},
		{
			name:         "create with nil nodes",
			nodes:        nil,
			expectedLen:  0,
			expectNonNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrb := NewRoundRobinBasic(tt.nodes)
			if tt.expectNonNil && rrb == nil {
				t.Errorf("NewRoundRobinBasic() returned nil")
			}
			if rrb != nil && len(rrb.nodes) != tt.expectedLen {
				t.Errorf("NewRoundRobinBasic() nodes length = %d; expected %d", len(rrb.nodes), tt.expectedLen)
			}
			if rrb != nil && rrb.index != 0 {
				t.Errorf("NewRoundRobinBasic() initial index = %d; expected 0", rrb.index)
			}
		})
	}
}

// TestRoundRobinBasicNext tests the Next method
func TestRoundRobinBasicNext(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []int
		callCount       int
		expectedPattern []int
	}{
		{
			name:            "single node returns same value",
			nodes:           []int{1},
			callCount:       5,
			expectedPattern: []int{1, 1, 1, 1, 1},
		},
		{
			name:            "two nodes alternate",
			nodes:           []int{1, 2},
			callCount:       4,
			expectedPattern: []int{2, 1, 2, 1},
		},
		{
			name:            "three nodes rotate",
			nodes:           []int{1, 2, 3},
			callCount:       6,
			expectedPattern: []int{2, 3, 1, 2, 3, 1},
		},
		{
			name:            "empty nodes return zero value",
			nodes:           []int{},
			callCount:       3,
			expectedPattern: []int{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrb := NewRoundRobinBasic(tt.nodes)

			for i := 0; i < tt.callCount; i++ {
				result := rrb.Next()
				if result != tt.expectedPattern[i] {
					t.Errorf("Next() call %d = %d; expected %d", i+1, result, tt.expectedPattern[i])
				}
			}
		})
	}
}

// TestNewRoundRobinSmoothWeighted tests the NewRoundRobinSmoothWeighted constructor
func TestNewRoundRobinSmoothWeighted(t *testing.T) {
	tests := []struct {
		name                string
		nodes               []WeightedNode[string]
		expectedLen         int
		expectedTotalWeight int
	}{
		{
			name: "create with single weighted node",
			nodes: []WeightedNode[string]{
				{Value: "server1", Weight: 5},
			},
			expectedLen:         1,
			expectedTotalWeight: 5,
		},
		{
			name: "create with multiple weighted nodes",
			nodes: []WeightedNode[string]{
				{Value: "server1", Weight: 5},
				{Value: "server2", Weight: 3},
				{Value: "server3", Weight: 2},
			},
			expectedLen:         3,
			expectedTotalWeight: 10,
		},
		{
			name:                "create with empty nodes",
			nodes:               []WeightedNode[string]{},
			expectedLen:         0,
			expectedTotalWeight: 0,
		},
		{
			name: "create with zero weight nodes filters them out",
			nodes: []WeightedNode[string]{
				{Value: "server1", Weight: 5},
				{Value: "server2", Weight: 0},
				{Value: "server3", Weight: 3},
			},
			expectedLen:         2,
			expectedTotalWeight: 8,
		},
		{
			name: "create with negative weight nodes filters them out",
			nodes: []WeightedNode[string]{
				{Value: "server1", Weight: 5},
				{Value: "server2", Weight: -1},
				{Value: "server3", Weight: 3},
			},
			expectedLen:         2,
			expectedTotalWeight: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)
			if rrs == nil {
				t.Errorf("NewRoundRobinSmoothWeighted() returned nil")
			}
			if len(rrs.nodes) != tt.expectedLen {
				t.Errorf("NewRoundRobinSmoothWeighted() nodes length = %d; expected %d", len(rrs.nodes), tt.expectedLen)
			}
			if rrs.totalWeight != tt.expectedTotalWeight {
				t.Errorf("NewRoundRobinSmoothWeighted() totalWeight = %d; expected %d", rrs.totalWeight, tt.expectedTotalWeight)
			}
		})
	}
}

// TestRoundRobinSmoothWeightedNext tests the Next method
func TestRoundRobinSmoothWeightedNext(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []WeightedNode[string]
		callCount       int
		expectedPattern []string
	}{
		{
			name: "single weighted node",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 5},
			},
			callCount:       3,
			expectedPattern: []string{"A", "A", "A"},
		},
		{
			name: "equal weights alternate",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 1},
				{Value: "B", Weight: 1},
			},
			callCount:       4,
			expectedPattern: []string{"A", "B", "A", "B"},
		},
		{
			name: "2:1 weight ratio",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 2},
				{Value: "B", Weight: 1},
			},
			callCount:       6,
			expectedPattern: []string{"A", "B", "A", "A", "B", "A"},
		},
		{
			name: "5:1:1 weight ratio",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 5},
				{Value: "B", Weight: 1},
				{Value: "C", Weight: 1},
			},
			callCount:       7,
			expectedPattern: []string{"A", "A", "B", "A", "C", "A", "A"},
		},
		{
			name: "classic 5:3:2 weight ratio",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 5},
				{Value: "B", Weight: 3},
				{Value: "C", Weight: 2},
			},
			callCount:       10,
			expectedPattern: []string{"A", "B", "C", "A", "A", "B", "A", "C", "B", "A"},
		},
		{
			name:            "empty nodes return zero value",
			nodes:           []WeightedNode[string]{},
			callCount:       2,
			expectedPattern: []string{"", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)

			for i := 0; i < tt.callCount; i++ {
				result := rrs.Next()
				if result != tt.expectedPattern[i] {
					t.Errorf("Next() call %d = %s; expected %s", i+1, result, tt.expectedPattern[i])
				}
			}
		})
	}
}

// TestRoundRobinSmoothWeightedDistribution tests distribution of weighted nodes
func TestRoundRobinSmoothWeightedDistribution(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []WeightedNode[string]
		totalCalls    int
		expectedCount map[string]int
	}{
		{
			name: "5:3:2 weight distribution",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 5},
				{Value: "B", Weight: 3},
				{Value: "C", Weight: 2},
			},
			totalCalls: 100,
			expectedCount: map[string]int{
				"A": 50, // 5/10 * 100
				"B": 30, // 3/10 * 100
				"C": 20, // 2/10 * 100
			},
		},
		{
			name: "equal weight distribution",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 1},
				{Value: "B", Weight: 1},
				{Value: "C", Weight: 1},
			},
			totalCalls: 90,
			expectedCount: map[string]int{
				"A": 30,
				"B": 30,
				"C": 30,
			},
		},
		{
			name: "10:5:3:2 weight distribution",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 10},
				{Value: "B", Weight: 5},
				{Value: "C", Weight: 3},
				{Value: "D", Weight: 2},
			},
			totalCalls: 100,
			expectedCount: map[string]int{
				"A": 50, // 10/20 * 100
				"B": 25, // 5/20 * 100
				"C": 15, // 3/20 * 100
				"D": 10, // 2/20 * 100
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)
			counts := make(map[string]int)

			// Make calls and count occurrences
			for i := 0; i < tt.totalCalls; i++ {
				result := rrs.Next()
				counts[result]++
			}

			// Verify distribution matches expected
			for node, expectedCount := range tt.expectedCount {
				actualCount := counts[node]
				if actualCount != expectedCount {
					t.Errorf("Distribution: node %s appeared %d times; expected %d", node, actualCount, expectedCount)
				}
			}
		})
	}
}

// TestRoundRobinSmoothWeightedReset tests the Reset method
func TestRoundRobinSmoothWeightedReset(t *testing.T) {
	tests := []struct {
		name             string
		nodes            []WeightedNode[string]
		callsBeforeReset int
		callsAfterReset  int
		expectedPattern  []string
	}{
		{
			name: "reset after some calls",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 2},
				{Value: "B", Weight: 1},
			},
			callsBeforeReset: 3,
			callsAfterReset:  3,
			expectedPattern:  []string{"A", "B", "A"}, // Same pattern after reset
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)

			// Make some calls
			for i := 0; i < tt.callsBeforeReset; i++ {
				rrs.Next()
			}

			// Reset
			rrs.Reset()

			// Verify pattern restarts correctly
			for i := 0; i < tt.callsAfterReset; i++ {
				result := rrs.Next()
				if result != tt.expectedPattern[i] {
					t.Errorf("After reset, call %d = %s; expected %s", i+1, result, tt.expectedPattern[i])
				}
			}
		})
	}
}

// TestRoundRobinSmoothWeightedConcurrent tests concurrent access
func TestRoundRobinSmoothWeightedConcurrent(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []WeightedNode[string]
		goroutines    int
		callsPerGo    int
		expectedTotal map[string]int
	}{
		{
			name: "concurrent access with 2:1 weight",
			nodes: []WeightedNode[string]{
				{Value: "A", Weight: 2},
				{Value: "B", Weight: 1},
			},
			goroutines: 10,
			callsPerGo: 30,
			expectedTotal: map[string]int{
				"A": 200, // 2/3 * 300
				"B": 100, // 1/3 * 300
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)
			var wg sync.WaitGroup
			results := make(chan string, tt.goroutines*tt.callsPerGo)

			// Launch goroutines
			for i := 0; i < tt.goroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < tt.callsPerGo; j++ {
						result := rrs.Next()
						results <- result
					}
				}()
			}

			// Wait for all goroutines to complete
			wg.Wait()
			close(results)

			// Collect and verify results
			counts := make(map[string]int)
			for result := range results {
				counts[result]++
			}

			// Verify counts match expected
			for node, expectedCount := range tt.expectedTotal {
				actualCount := counts[node]
				if actualCount != expectedCount {
					t.Errorf("Concurrent test: node %s appeared %d times; expected %d", node, actualCount, expectedCount)
				}
			}
		})
	}
}

// TestRoundRobinSmoothWeightedWithIntegers tests with integer values
func TestRoundRobinSmoothWeightedWithIntegers(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []WeightedNode[int]
		callCount       int
		expectedPattern []int
	}{
		{
			name: "integer nodes with 3:2:1 weight",
			nodes: []WeightedNode[int]{
				{Value: 1, Weight: 3},
				{Value: 2, Weight: 2},
				{Value: 3, Weight: 1},
			},
			callCount:       6,
			expectedPattern: []int{1, 2, 1, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)

			for i := 0; i < tt.callCount; i++ {
				result := rrs.Next()
				if result != tt.expectedPattern[i] {
					t.Errorf("Next() call %d = %d; expected %d", i+1, result, tt.expectedPattern[i])
				}
			}
		})
	}
}

// TestRoundRobinSmoothWeightedWithStructs tests with struct values
func TestRoundRobinSmoothWeightedWithStructs(t *testing.T) {
	type Server struct {
		Name string
		Port int
	}

	tests := []struct {
		name            string
		nodes           []WeightedNode[Server]
		callCount       int
		expectedPattern []Server
	}{
		{
			name: "struct nodes with different weights",
			nodes: []WeightedNode[Server]{
				{Value: Server{Name: "server1", Port: 8080}, Weight: 2},
				{Value: Server{Name: "server2", Port: 8081}, Weight: 1},
			},
			callCount: 3,
			expectedPattern: []Server{
				{Name: "server1", Port: 8080},
				{Name: "server2", Port: 8081},
				{Name: "server1", Port: 8080},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrs := NewRoundRobinSmoothWeighted(tt.nodes)

			for i := 0; i < tt.callCount; i++ {
				result := rrs.Next()
				if result != tt.expectedPattern[i] {
					t.Errorf("Next() call %d = %+v; expected %+v", i+1, result, tt.expectedPattern[i])
				}
			}
		})
	}
}

// TestRoundRobinSmoothWeightedSmoothness tests the smoothness of distribution
func TestRoundRobinSmoothWeightedSmoothness(t *testing.T) {
	// This test verifies that smooth weighted round-robin distributes
	// requests more evenly than basic weighted round-robin
	t.Run("smoothness test", func(t *testing.T) {
		nodes := []WeightedNode[string]{
			{Value: "A", Weight: 5},
			{Value: "B", Weight: 1},
		}
		rrs := NewRoundRobinSmoothWeighted(nodes)

		results := make([]string, 12)
		for i := 0; i < 12; i++ {
			results[i] = rrs.Next()
		}

		// With smooth weighted RR, we should see B distributed more evenly
		// Rather than seeing all A's first then B's
		// Pattern should be something like: A A A A A B A A A A A B
		// Not: A A A A A A A A A A B B

		// Count consecutive A's at the start
		consecutiveA := 0
		for i := 0; i < len(results) && results[i] == "A"; i++ {
			consecutiveA++
		}

		// In smooth weighted RR, we shouldn't see more than 5 consecutive A's at start
		if consecutiveA > 5 {
			t.Errorf("Smoothness test failed: saw %d consecutive A's at start; expected <= 5 (pattern: %v)", consecutiveA, results)
		}

		// Verify we have correct total counts
		countA, countB := 0, 0
		for _, result := range results {
			if result == "A" {
				countA++
			} else if result == "B" {
				countB++
			}
		}

		if countA != 10 || countB != 2 {
			t.Errorf("Smoothness test: incorrect distribution A=%d, B=%d; expected A=10, B=2", countA, countB)
		}
	})
}

// TestNewRoundRobinLimitDDL tests the NewRoundRobinLimitDDL constructor
func TestNewRoundRobinLimitDDL(t *testing.T) {
	tests := []struct {
		name        string
		nodes       []*RoundRobinLimitDDLNode[string]
		expectedLen int
	}{
		{
			name: "create with single node",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "server1", Limit: 2},
			},
			expectedLen: 1,
		},
		{
			name: "create with multiple nodes",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "server1", Limit: 2},
				{Value: "server2", Limit: 3},
				{Value: "server3", Limit: 1},
			},
			expectedLen: 3,
		},
		{
			name:        "create with empty nodes",
			nodes:       []*RoundRobinLimitDDLNode[string]{},
			expectedLen: 0,
		},
		{
			name: "filter out zero limit nodes",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "server1", Limit: 2},
				{Value: "server2", Limit: 0},
				{Value: "server3", Limit: 1},
			},
			expectedLen: 2,
		},
		{
			name: "filter out negative limit nodes",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "server1", Limit: 2},
				{Value: "server2", Limit: -1},
			},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrl := NewRoundRobinLimitDDL(tt.nodes)
			if rrl == nil {
				t.Errorf("NewRoundRobinLimitDDL() returned nil")
			}
			if rrl.Len() != tt.expectedLen {
				t.Errorf("NewRoundRobinLimitDDL() length = %d; expected %d", rrl.Len(), tt.expectedLen)
			}
		})
	}
}

// TestRoundRobinLimitDDLNext tests the Next method
func TestRoundRobinLimitDDLNext(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []*RoundRobinLimitDDLNode[string]
		callCount       int
		expectedPattern []string
	}{
		{
			name: "single node with limit 2",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "A", Limit: 2},
			},
			callCount:       3,
			expectedPattern: []string{"A", "A", ""},
		},
		{
			name: "two nodes with equal limits",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "A", Limit: 2},
				{Value: "B", Limit: 2},
			},
			callCount:       4,
			expectedPattern: []string{"A", "B", "A", "B"},
		},
		{
			name: "two nodes with different limits",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "A", Limit: 1},
				{Value: "B", Limit: 2},
			},
			callCount:       4,
			expectedPattern: []string{"A", "B", "B", ""},
		},
		{
			name: "three nodes rotate until exhausted",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "A", Limit: 1},
				{Value: "B", Limit: 1},
				{Value: "C", Limit: 1},
			},
			callCount:       4,
			expectedPattern: []string{"A", "B", "C", ""},
		},
		{
			name:            "empty nodes return zero value",
			nodes:           []*RoundRobinLimitDDLNode[string]{},
			callCount:       2,
			expectedPattern: []string{"", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrl := NewRoundRobinLimitDDL(tt.nodes)

			for i := 0; i < tt.callCount; i++ {
				result := rrl.Next()
				if result != tt.expectedPattern[i] {
					t.Errorf("Next() call %d = %s; expected %s", i+1, result, tt.expectedPattern[i])
				}
			}
		})
	}
}

// TestRoundRobinLimitDDLRelease tests the Release method
func TestRoundRobinLimitDDLRelease(t *testing.T) {
	tests := []struct {
		name         string
		setupCalls   int
		releasedNode int // index of node to release (-1 for none)
		afterRelease []string
	}{
		{
			name:         "release brings node back",
			setupCalls:   3,                  // A twice, now at limit
			releasedNode: 0,                  // release A
			afterRelease: []string{"B", "A"}, // B then A (now available again)
		},
		{
			name:         "release node not at limit has no effect",
			setupCalls:   1,                       // A once, under limit
			releasedNode: 0,                       // release A
			afterRelease: []string{"B", "A", "B"}, // Normal rotation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeA := &RoundRobinLimitDDLNode[string]{Value: "A", Limit: 2}
			nodeB := &RoundRobinLimitDDLNode[string]{Value: "B", Limit: 2}
			nodes := []*RoundRobinLimitDDLNode[string]{nodeA, nodeB}
			rrl := NewRoundRobinLimitDDL(nodes)

			// Setup: call Next() specified times
			for i := 0; i < tt.setupCalls; i++ {
				rrl.Next()
			}

			// Release node if specified
			if tt.releasedNode >= 0 {
				rrl.Release(nodes[tt.releasedNode])
			}

			// Verify pattern after release
			for i, expected := range tt.afterRelease {
				result := rrl.Next()
				if result != expected {
					t.Errorf("After release, call %d = %s; expected %s", i+1, result, expected)
				}
			}
		})
	}
}

// TestRoundRobinLimitDDLReleaseNil tests Release with nil node
func TestRoundRobinLimitDDLReleaseNil(t *testing.T) {
	t.Run("release nil node does not panic", func(t *testing.T) {
		nodes := []*RoundRobinLimitDDLNode[string]{
			{Value: "A", Limit: 1},
		}
		rrl := NewRoundRobinLimitDDL(nodes)

		// Should not panic
		rrl.Release(nil)

		// Should still work normally
		result := rrl.Next()
		if result != "A" {
			t.Errorf("Next() after Release(nil) = %s; expected A", result)
		}
	})
}

// TestRoundRobinLimitDDLConcurrent tests concurrent access
func TestRoundRobinLimitDDLConcurrent(t *testing.T) {
	tests := []struct {
		name       string
		nodes      []*RoundRobinLimitDDLNode[int]
		goroutines int
		operations int
	}{
		{
			name: "concurrent next calls",
			nodes: []*RoundRobinLimitDDLNode[int]{
				{Value: 1, Limit: 50},
				{Value: 2, Limit: 50},
				{Value: 3, Limit: 50},
			},
			goroutines: 10,
			operations: 10,
		},
		{
			name: "concurrent next and release",
			nodes: []*RoundRobinLimitDDLNode[int]{
				{Value: 1, Limit: 20},
				{Value: 2, Limit: 20},
			},
			goroutines: 5,
			operations: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrl := NewRoundRobinLimitDDL(tt.nodes)
			var wg sync.WaitGroup

			// Concurrent Next() calls
			for i := 0; i < tt.goroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < tt.operations; j++ {
						rrl.Next()
					}
				}()
			}

			wg.Wait()

			// Verify no race conditions (test will fail with -race flag if issues)
			// All nodes should be exhausted or have consistent state
		})
	}
}

// TestRoundRobinLimitDDLLen tests the Len method
func TestRoundRobinLimitDDLLen(t *testing.T) {
	tests := []struct {
		name         string
		nodes        []*RoundRobinLimitDDLNode[string]
		operations   []string // "next" or "release:index"
		expectedLens []int
	}{
		{
			name: "length decreases as nodes exhaust",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "A", Limit: 1},
				{Value: "B", Limit: 1},
			},
			operations:   []string{"next", "next"},
			expectedLens: []int{1, 0}, // After each Next()
		},
		{
			name: "length increases when node released",
			nodes: []*RoundRobinLimitDDLNode[string]{
				{Value: "A", Limit: 1},
			},
			operations:   []string{"next", "release:0"},
			expectedLens: []int{0, 1}, // Exhausted, then restored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rrl := NewRoundRobinLimitDDL(tt.nodes)

			for i, op := range tt.operations {
				if op == "next" {
					rrl.Next()
				} else if len(op) > 8 && op[:7] == "release" {
					var idx int
					fmt.Sscanf(op, "release:%d", &idx)
					rrl.Release(tt.nodes[idx])
				}

				length := rrl.Len()
				if length != tt.expectedLens[i] {
					t.Errorf("After operation %d (%s): Len() = %d; expected %d", i+1, op, length, tt.expectedLens[i])
				}
			}
		})
	}
}

// TestRoundRobinLimitDDLReset tests the Reset method
func TestRoundRobinLimitDDLReset(t *testing.T) {
	t.Run("reset restores all nodes", func(t *testing.T) {
		nodeA := &RoundRobinLimitDDLNode[string]{Value: "A", Limit: 1}
		nodeB := &RoundRobinLimitDDLNode[string]{Value: "B", Limit: 1}
		nodes := []*RoundRobinLimitDDLNode[string]{nodeA, nodeB}
		rrl := NewRoundRobinLimitDDL(nodes)

		// Exhaust all nodes
		rrl.Next() // A
		rrl.Next() // B

		if rrl.Len() != 0 {
			t.Errorf("Before reset: Len() = %d; expected 0", rrl.Len())
		}

		// Reset
		rrl.Reset(nodes)

		if rrl.Len() != 2 {
			t.Errorf("After reset: Len() = %d; expected 2", rrl.Len())
		}

		// Verify nodes work again
		result1 := rrl.Next()
		result2 := rrl.Next()
		if (result1 != "A" && result1 != "B") || (result2 != "A" && result2 != "B") {
			t.Errorf("After reset: got %s, %s; expected A and B", result1, result2)
		}
	})
}

// TestRoundRobinLimitDDLWithDifferentTypes tests with various value types
func TestRoundRobinLimitDDLWithDifferentTypes(t *testing.T) {
	t.Run("integer values", func(t *testing.T) {
		nodes := []*RoundRobinLimitDDLNode[int]{
			{Value: 1, Limit: 2},
			{Value: 2, Limit: 2},
		}
		rrl := NewRoundRobinLimitDDL(nodes)

		result1 := rrl.Next()
		result2 := rrl.Next()

		if (result1 != 1 && result1 != 2) || (result2 != 1 && result2 != 2) {
			t.Errorf("Next() returned invalid values: %d, %d", result1, result2)
		}
	})

	t.Run("struct values", func(t *testing.T) {
		type Server struct {
			Name string
			Port int
		}

		nodes := []*RoundRobinLimitDDLNode[Server]{
			{Value: Server{Name: "server1", Port: 8080}, Limit: 1},
			{Value: Server{Name: "server2", Port: 8081}, Limit: 1},
		}
		rrl := NewRoundRobinLimitDDL(nodes)

		result1 := rrl.Next()
		result2 := rrl.Next()

		if result1.Name == "" || result2.Name == "" {
			t.Errorf("Next() returned invalid struct values")
		}
	})
}

// TestRoundRobinLimitDDLEdgeCases tests edge cases
func TestRoundRobinLimitDDLEdgeCases(t *testing.T) {
	t.Run("large limit value", func(t *testing.T) {
		node := &RoundRobinLimitDDLNode[string]{Value: "A", Limit: 1000}
		rrl := NewRoundRobinLimitDDL([]*RoundRobinLimitDDLNode[string]{node})

		// Should be able to call Next() many times
		for i := 0; i < 500; i++ {
			result := rrl.Next()
			if result != "A" {
				t.Errorf("Next() call %d = %s; expected A", i+1, result)
				break
			}
		}
	})

	t.Run("release more than incremented", func(t *testing.T) {
		node := &RoundRobinLimitDDLNode[string]{Value: "A", Limit: 2}
		rrl := NewRoundRobinLimitDDL([]*RoundRobinLimitDDLNode[string]{node})

		// Release without calling Next()
		rrl.Release(node)

		// Should not go negative
		result := rrl.Next()
		if result != "A" {
			t.Errorf("Next() = %s; expected A", result)
		}
	})
}
