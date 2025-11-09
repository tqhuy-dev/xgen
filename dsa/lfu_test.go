package dsa

import (
	"testing"
)

// TestNewLFUCache tests the NewLFUCache constructor
func TestNewLFUCache(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		expectPanic bool
	}{
		{
			name:        "valid capacity",
			capacity:    5,
			expectPanic: false,
		},
		{
			name:        "capacity of 1",
			capacity:    1,
			expectPanic: false,
		},
		{
			name:        "large capacity",
			capacity:    1000,
			expectPanic: false,
		},
		{
			name:        "zero capacity should panic",
			capacity:    0,
			expectPanic: true,
		},
		{
			name:        "negative capacity should panic",
			capacity:    -1,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewLFUCache() should have panicked")
					}
				}()
				NewLFUCache[int, int](tt.capacity)
			} else {
				cache := NewLFUCache[int, int](tt.capacity)
				if cache == nil {
					t.Errorf("NewLFUCache() returned nil")
				}
				if cache.Capacity() != tt.capacity {
					t.Errorf("NewLFUCache() capacity = %d; expected %d", cache.Capacity(), tt.capacity)
				}
				if cache.Len() != 0 {
					t.Errorf("NewLFUCache() initial length = %d; expected 0", cache.Len())
				}
			}
		})
	}
}

// TestLFUCachePut tests the Put method
func TestLFUCachePut(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		operations  []struct{ key, value int }
		expectedLen int
	}{
		{
			name:     "put single item",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
			},
			expectedLen: 1,
		},
		{
			name:     "put multiple items within capacity",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			expectedLen: 3,
		},
		{
			name:     "put items exceeding capacity",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
				{4, 400},
			},
			expectedLen: 3,
		},
		{
			name:     "update existing key",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{1, 150}, // Update key 1
			},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)
			for _, op := range tt.operations {
				cache.Put(op.key, op.value)
			}

			if cache.Len() != tt.expectedLen {
				t.Errorf("Put() cache length = %d; expected %d", cache.Len(), tt.expectedLen)
			}
		})
	}
}

// TestLFUCacheGet tests the Get method
func TestLFUCacheGet(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		putOperations []struct{ key, value int }
		getKey        int
		expectedValue int
		expectedFound bool
	}{
		{
			name:     "get existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			getKey:        2,
			expectedValue: 200,
			expectedFound: true,
		},
		{
			name:     "get non-existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
			getKey:        3,
			expectedValue: 0,
			expectedFound: false,
		},
		{
			name:          "get from empty cache",
			capacity:      3,
			putOperations: []struct{ key, value int }{},
			getKey:        1,
			expectedValue: 0,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)
			for _, op := range tt.putOperations {
				cache.Put(op.key, op.value)
			}

			value, found := cache.Get(tt.getKey)
			if found != tt.expectedFound {
				t.Errorf("Get(%d) found = %v; expected %v", tt.getKey, found, tt.expectedFound)
			}
			if value != tt.expectedValue {
				t.Errorf("Get(%d) value = %d; expected %d", tt.getKey, value, tt.expectedValue)
			}
		})
	}
}

// TestLFUCacheFrequency tests frequency tracking
func TestLFUCacheFrequency(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		operations        []string // "put:key:value" or "get:key"
		checkKey          int
		expectedFrequency int
		expectedFound     bool
	}{
		{
			name:     "frequency after single put",
			capacity: 3,
			operations: []string{
				"put:1:100",
			},
			checkKey:          1,
			expectedFrequency: 1,
			expectedFound:     true,
		},
		{
			name:     "frequency after get",
			capacity: 3,
			operations: []string{
				"put:1:100",
				"get:1",
			},
			checkKey:          1,
			expectedFrequency: 2,
			expectedFound:     true,
		},
		{
			name:     "frequency after multiple gets",
			capacity: 3,
			operations: []string{
				"put:1:100",
				"get:1",
				"get:1",
				"get:1",
			},
			checkKey:          1,
			expectedFrequency: 4,
			expectedFound:     true,
		},
		{
			name:     "frequency after update",
			capacity: 3,
			operations: []string{
				"put:1:100",
				"get:1",
				"put:1:150", // Update increases frequency
			},
			checkKey:          1,
			expectedFrequency: 3,
			expectedFound:     true,
		},
		{
			name:     "non-existing key",
			capacity: 3,
			operations: []string{
				"put:1:100",
			},
			checkKey:          2,
			expectedFrequency: 0,
			expectedFound:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)

			for _, op := range tt.operations {
				var key, value int
				if op[:3] == "put" {
					_, _ = parseOperation(op, &key, &value)
					cache.Put(key, value)
				} else if op[:3] == "get" {
					_, _ = parseOperation(op, &key, nil)
					cache.Get(key)
				}
			}

			freq, found := cache.GetFrequency(tt.checkKey)
			if found != tt.expectedFound {
				t.Errorf("GetFrequency(%d) found = %v; expected %v", tt.checkKey, found, tt.expectedFound)
			}
			if freq != tt.expectedFrequency {
				t.Errorf("GetFrequency(%d) = %d; expected %d", tt.checkKey, freq, tt.expectedFrequency)
			}
		})
	}
}

// TestLFUCacheEviction tests eviction behavior
func TestLFUCacheEviction(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		operations   []string
		shouldContain map[int]bool
	}{
		{
			name:     "evict least frequently used",
			capacity: 2,
			operations: []string{
				"put:1:100",
				"put:2:200",
				"get:1",     // 1 has freq 2, 2 has freq 1
				"put:3:300", // Should evict 2 (least frequent)
			},
			shouldContain: map[int]bool{
				1: true,
				2: false,
				3: true,
			},
		},
		{
			name:     "evict least recently used when frequencies are equal",
			capacity: 2,
			operations: []string{
				"put:1:100", // freq 1
				"put:2:200", // freq 1
				"put:3:300", // Should evict 1 (oldest with freq 1)
			},
			shouldContain: map[int]bool{
				1: false,
				2: true,
				3: true,
			},
		},
		{
			name:     "frequent access prevents eviction",
			capacity: 2,
			operations: []string{
				"put:1:100",
				"put:2:200",
				"get:2", // 2 now has freq 2
				"get:2", // 2 now has freq 3
				"put:3:300", // Should evict 1 (freq 1)
			},
			shouldContain: map[int]bool{
				1: false,
				2: true,
				3: true,
			},
		},
		{
			name:     "update increases frequency",
			capacity: 2,
			operations: []string{
				"put:1:100",
				"put:2:200",
				"put:1:150", // Update increases freq of 1
				"put:3:300", // Should evict 2
			},
			shouldContain: map[int]bool{
				1: true,
				2: false,
				3: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)

			for _, op := range tt.operations {
				var key, value int
				if op[:3] == "put" {
					_, _ = parseOperation(op, &key, &value)
					cache.Put(key, value)
				} else if op[:3] == "get" {
					_, _ = parseOperation(op, &key, nil)
					cache.Get(key)
				}
			}

			for key, shouldExist := range tt.shouldContain {
				exists := cache.Contains(key)
				if exists != shouldExist {
					t.Errorf("Eviction test: key %d exists = %v; expected %v", key, exists, shouldExist)
				}
			}
		})
	}
}

// TestLFUCacheRemove tests the Remove method
func TestLFUCacheRemove(t *testing.T) {
	tests := []struct {
		name            string
		capacity        int
		putOperations   []struct{ key, value int }
		removeKey       int
		expectedRemoved bool
		expectedLen     int
	}{
		{
			name:     "remove existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			removeKey:       2,
			expectedRemoved: true,
			expectedLen:     2,
		},
		{
			name:     "remove non-existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
			removeKey:       3,
			expectedRemoved: false,
			expectedLen:     2,
		},
		{
			name:            "remove from empty cache",
			capacity:        3,
			putOperations:   []struct{ key, value int }{},
			removeKey:       1,
			expectedRemoved: false,
			expectedLen:     0,
		},
		{
			name:     "remove last item",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
			},
			removeKey:       1,
			expectedRemoved: true,
			expectedLen:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)
			for _, op := range tt.putOperations {
				cache.Put(op.key, op.value)
			}

			removed := cache.Remove(tt.removeKey)
			if removed != tt.expectedRemoved {
				t.Errorf("Remove(%d) = %v; expected %v", tt.removeKey, removed, tt.expectedRemoved)
			}

			if cache.Len() != tt.expectedLen {
				t.Errorf("Remove() cache length = %d; expected %d", cache.Len(), tt.expectedLen)
			}
		})
	}
}

// TestLFUCacheClear tests the Clear method
func TestLFUCacheClear(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		putOperations []struct{ key, value int }
	}{
		{
			name:     "clear non-empty cache",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
		},
		{
			name:          "clear empty cache",
			capacity:      3,
			putOperations: []struct{ key, value int }{},
		},
		{
			name:     "clear cache at capacity",
			capacity: 2,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)
			for _, op := range tt.putOperations {
				cache.Put(op.key, op.value)
			}

			cache.Clear()

			if cache.Len() != 0 {
				t.Errorf("Clear() cache length = %d; expected 0", cache.Len())
			}
			if cache.Capacity() != tt.capacity {
				t.Errorf("Clear() changed capacity to %d; expected %d", cache.Capacity(), tt.capacity)
			}
		})
	}
}

// TestLFUCacheContains tests the Contains method
func TestLFUCacheContains(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		putOperations []struct{ key, value int }
		checkKey      int
		expectedFound bool
	}{
		{
			name:     "contains existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
			checkKey:      1,
			expectedFound: true,
		},
		{
			name:     "contains non-existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
			checkKey:      3,
			expectedFound: false,
		},
		{
			name:          "contains in empty cache",
			capacity:      3,
			putOperations: []struct{ key, value int }{},
			checkKey:      1,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)
			for _, op := range tt.putOperations {
				cache.Put(op.key, op.value)
			}

			found := cache.Contains(tt.checkKey)
			if found != tt.expectedFound {
				t.Errorf("Contains(%d) = %v; expected %v", tt.checkKey, found, tt.expectedFound)
			}
		})
	}
}

// TestLFUCachePeek tests the Peek method
func TestLFUCachePeek(t *testing.T) {
	tests := []struct {
		name               string
		capacity           int
		putOperations      []struct{ key, value int }
		peekKey            int
		expectedValue      int
		expectedFound      bool
		expectedFreqBefore int
		expectedFreqAfter  int
	}{
		{
			name:     "peek existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			peekKey:            1,
			expectedValue:      100,
			expectedFound:      true,
			expectedFreqBefore: 1,
			expectedFreqAfter:  1, // Peek doesn't change frequency
		},
		{
			name:     "peek non-existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
			peekKey:            3,
			expectedValue:      0,
			expectedFound:      false,
			expectedFreqBefore: 0,
			expectedFreqAfter:  0,
		},
		{
			name:               "peek from empty cache",
			capacity:           3,
			putOperations:      []struct{ key, value int }{},
			peekKey:            1,
			expectedValue:      0,
			expectedFound:      false,
			expectedFreqBefore: 0,
			expectedFreqAfter:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLFUCache[int, int](tt.capacity)
			for _, op := range tt.putOperations {
				cache.Put(op.key, op.value)
			}

			// Check frequency before peek
			freqBefore, _ := cache.GetFrequency(tt.peekKey)
			if freqBefore != tt.expectedFreqBefore {
				t.Errorf("Peek() frequency before = %d; expected %d", freqBefore, tt.expectedFreqBefore)
			}

			value, found := cache.Peek(tt.peekKey)
			if found != tt.expectedFound {
				t.Errorf("Peek(%d) found = %v; expected %v", tt.peekKey, found, tt.expectedFound)
			}
			if value != tt.expectedValue {
				t.Errorf("Peek(%d) value = %d; expected %d", tt.peekKey, value, tt.expectedValue)
			}

			// Verify frequency didn't change
			freqAfter, _ := cache.GetFrequency(tt.peekKey)
			if freqAfter != tt.expectedFreqAfter {
				t.Errorf("Peek() changed frequency to %d; expected %d", freqAfter, tt.expectedFreqAfter)
			}
		})
	}
}

// TestLFUCacheGenericTypes tests the cache with different types
func TestLFUCacheGenericTypes(t *testing.T) {
	t.Run("string key string value", func(t *testing.T) {
		cache := NewLFUCache[string, string](3)
		cache.Put("name", "Alice")
		cache.Put("city", "NYC")

		value, found := cache.Get("name")
		if !found {
			t.Errorf("Get(name) not found")
		}
		if value != "Alice" {
			t.Errorf("Get(name) = %s; expected Alice", value)
		}
	})

	t.Run("string key struct value", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		cache := NewLFUCache[string, User](2)
		cache.Put("user1", User{ID: 1, Name: "Alice"})
		cache.Put("user2", User{ID: 2, Name: "Bob"})

		value, found := cache.Get("user1")
		if !found {
			t.Errorf("Get(user1) not found")
		}
		if value.Name != "Alice" {
			t.Errorf("Get(user1).Name = %s; expected Alice", value.Name)
		}
	})

	t.Run("int key pointer value", func(t *testing.T) {
		cache := NewLFUCache[int, *int](2)
		val1 := 100
		val2 := 200
		cache.Put(1, &val1)
		cache.Put(2, &val2)

		value, found := cache.Get(1)
		if !found {
			t.Errorf("Get(1) not found")
		}
		if *value != 100 {
			t.Errorf("*Get(1) = %d; expected 100", *value)
		}
	})
}

// TestLFUCacheCapacityOne tests edge case with capacity of 1
func TestLFUCacheCapacityOne(t *testing.T) {
	cache := NewLFUCache[int, int](1)

	cache.Put(1, 100)
	if cache.Len() != 1 {
		t.Errorf("cache length = %d; expected 1", cache.Len())
	}

	cache.Put(2, 200)
	if cache.Len() != 1 {
		t.Errorf("cache length = %d; expected 1", cache.Len())
	}

	// Should only contain key 2
	_, found := cache.Get(1)
	if found {
		t.Errorf("Get(1) found = true; expected false (should be evicted)")
	}

	value, found := cache.Get(2)
	if !found {
		t.Errorf("Get(2) found = false; expected true")
	}
	if value != 200 {
		t.Errorf("Get(2) = %d; expected 200", value)
	}
}

// TestLFUCacheComplexScenario tests a complex scenario
func TestLFUCacheComplexScenario(t *testing.T) {
	cache := NewLFUCache[int, int](3)

	// Put 3 items
	cache.Put(1, 100)
	cache.Put(2, 200)
	cache.Put(3, 300)

	// Access pattern: 1 (3 times), 2 (2 times), 3 (1 time)
	cache.Get(1)
	cache.Get(1)
	cache.Get(1)
	cache.Get(2)
	cache.Get(2)

	// Verify frequencies
	freq1, _ := cache.GetFrequency(1)
	freq2, _ := cache.GetFrequency(2)
	freq3, _ := cache.GetFrequency(3)

	if freq1 != 4 { // 1 put + 3 gets
		t.Errorf("Frequency of key 1 = %d; expected 4", freq1)
	}
	if freq2 != 3 { // 1 put + 2 gets
		t.Errorf("Frequency of key 2 = %d; expected 3", freq2)
	}
	if freq3 != 1 { // 1 put only
		t.Errorf("Frequency of key 3 = %d; expected 1", freq3)
	}

	// Add new item, should evict key 3 (lowest frequency)
	cache.Put(4, 400)

	if cache.Contains(3) {
		t.Errorf("Key 3 should have been evicted")
	}
	if !cache.Contains(1) || !cache.Contains(2) || !cache.Contains(4) {
		t.Errorf("Keys 1, 2, and 4 should be in cache")
	}
}

