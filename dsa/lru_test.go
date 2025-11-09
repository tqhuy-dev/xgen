package dsa

import (
	"fmt"
	"testing"
)

// TestNewLRUCache tests the NewLRUCache constructor
func TestNewLRUCache(t *testing.T) {
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
						t.Errorf("NewLRUCache() should have panicked")
					}
				}()
				NewLRUCache[int, int](tt.capacity)
			} else {
				cache := NewLRUCache[int, int](tt.capacity)
				if cache == nil {
					t.Errorf("NewLRUCache() returned nil")
				}
				if cache.Capacity() != tt.capacity {
					t.Errorf("NewLRUCache() capacity = %d; expected %d", cache.Capacity(), tt.capacity)
				}
				if cache.Len() != 0 {
					t.Errorf("NewLRUCache() initial length = %d; expected 0", cache.Len())
				}
			}
		})
	}
}

// TestLRUCachePut tests the Put method
func TestLRUCachePut(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		operations   []struct{ key, value int }
		expectedLen  int
		expectedKeys []int
	}{
		{
			name:     "put single item",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
			},
			expectedLen:  1,
			expectedKeys: []int{1},
		},
		{
			name:     "put multiple items within capacity",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			expectedLen:  3,
			expectedKeys: []int{3, 2, 1},
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
			expectedLen:  3,
			expectedKeys: []int{4, 3, 2}, // 1 should be evicted
		},
		{
			name:     "update existing key",
			capacity: 3,
			operations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{1, 150}, // Update key 1
			},
			expectedLen:  2,
			expectedKeys: []int{1, 2}, // 1 moved to front
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[int, int](tt.capacity)
			for _, op := range tt.operations {
				cache.Put(op.key, op.value)
			}

			if cache.Len() != tt.expectedLen {
				t.Errorf("Put() cache length = %d; expected %d", cache.Len(), tt.expectedLen)
			}

			keys := cache.Keys()
			if len(keys) != len(tt.expectedKeys) {
				t.Errorf("Put() keys length = %d; expected %d", len(keys), len(tt.expectedKeys))
				return
			}

			for i, key := range keys {
				if key != tt.expectedKeys[i] {
					t.Errorf("Put() keys[%d] = %d; expected %d", i, key, tt.expectedKeys[i])
				}
			}
		})
	}
}

// TestLRUCacheGet tests the Get method
func TestLRUCacheGet(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		putOperations []struct{ key, value int }
		getKey        int
		expectedValue int
		expectedFound bool
		expectedOrder []int // Order after get
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
			expectedOrder: []int{2, 3, 1}, // 2 moved to front
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
			expectedOrder: []int{2, 1}, // Order unchanged
		},
		{
			name:          "get from empty cache",
			capacity:      3,
			putOperations: []struct{ key, value int }{},
			getKey:        1,
			expectedValue: 0,
			expectedFound: false,
			expectedOrder: []int{},
		},
		{
			name:     "get updates access order",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			getKey:        1,
			expectedValue: 100,
			expectedFound: true,
			expectedOrder: []int{1, 3, 2}, // 1 moved to front
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[int, int](tt.capacity)
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

			keys := cache.Keys()
			if len(keys) != len(tt.expectedOrder) {
				t.Errorf("Get() keys length = %d; expected %d", len(keys), len(tt.expectedOrder))
				return
			}
			for i, key := range keys {
				if key != tt.expectedOrder[i] {
					t.Errorf("Get() keys[%d] = %d; expected %d", i, key, tt.expectedOrder[i])
				}
			}
		})
	}
}

// TestLRUCacheRemove tests the Remove method
func TestLRUCacheRemove(t *testing.T) {
	tests := []struct {
		name            string
		capacity        int
		putOperations   []struct{ key, value int }
		removeKey       int
		expectedRemoved bool
		expectedLen     int
		expectedKeys    []int
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
			expectedKeys:    []int{3, 1},
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
			expectedKeys:    []int{2, 1},
		},
		{
			name:            "remove from empty cache",
			capacity:        3,
			putOperations:   []struct{ key, value int }{},
			removeKey:       1,
			expectedRemoved: false,
			expectedLen:     0,
			expectedKeys:    []int{},
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
			expectedKeys:    []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[int, int](tt.capacity)
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

			keys := cache.Keys()
			if len(keys) != len(tt.expectedKeys) {
				t.Errorf("Remove() keys length = %d; expected %d", len(keys), len(tt.expectedKeys))
				return
			}
			for i, key := range keys {
				if key != tt.expectedKeys[i] {
					t.Errorf("Remove() keys[%d] = %d; expected %d", i, key, tt.expectedKeys[i])
				}
			}
		})
	}
}

// TestLRUCacheClear tests the Clear method
func TestLRUCacheClear(t *testing.T) {
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
			cache := NewLRUCache[int, int](tt.capacity)
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

// TestLRUCacheContains tests the Contains method
func TestLRUCacheContains(t *testing.T) {
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
			cache := NewLRUCache[int, int](tt.capacity)
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

// TestLRUCachePeek tests the Peek method
func TestLRUCachePeek(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		putOperations []struct{ key, value int }
		peekKey       int
		expectedValue int
		expectedFound bool
		expectedOrder []int // Order should not change after peek
	}{
		{
			name:     "peek existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
				{3, 300},
			},
			peekKey:       1,
			expectedValue: 100,
			expectedFound: true,
			expectedOrder: []int{3, 2, 1}, // Order unchanged
		},
		{
			name:     "peek non-existing key",
			capacity: 3,
			putOperations: []struct{ key, value int }{
				{1, 100},
				{2, 200},
			},
			peekKey:       3,
			expectedValue: 0,
			expectedFound: false,
			expectedOrder: []int{2, 1},
		},
		{
			name:          "peek from empty cache",
			capacity:      3,
			putOperations: []struct{ key, value int }{},
			peekKey:       1,
			expectedValue: 0,
			expectedFound: false,
			expectedOrder: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[int, int](tt.capacity)
			for _, op := range tt.putOperations {
				cache.Put(op.key, op.value)
			}

			value, found := cache.Peek(tt.peekKey)
			if found != tt.expectedFound {
				t.Errorf("Peek(%d) found = %v; expected %v", tt.peekKey, found, tt.expectedFound)
			}
			if value != tt.expectedValue {
				t.Errorf("Peek(%d) value = %d; expected %d", tt.peekKey, value, tt.expectedValue)
			}

			// Verify order didn't change
			keys := cache.Keys()
			if len(keys) != len(tt.expectedOrder) {
				t.Errorf("Peek() changed keys length to %d; expected %d", len(keys), len(tt.expectedOrder))
				return
			}
			for i, key := range keys {
				if key != tt.expectedOrder[i] {
					t.Errorf("Peek() changed order: keys[%d] = %d; expected %d", i, key, tt.expectedOrder[i])
				}
			}
		})
	}
}

// TestLRUCacheEviction tests eviction behavior
func TestLRUCacheEviction(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		operations   []string // "put:key:value" or "get:key"
		expectedKeys []int
		expectedLen  int
	}{
		{
			name:     "evict least recently used",
			capacity: 2,
			operations: []string{
				"put:1:100",
				"put:2:200",
				"put:3:300", // Should evict 1
			},
			expectedKeys: []int{3, 2},
			expectedLen:  2,
		},
		{
			name:     "get prevents eviction",
			capacity: 2,
			operations: []string{
				"put:1:100",
				"put:2:200",
				"get:1",     // Makes 1 most recently used
				"put:3:300", // Should evict 2
			},
			expectedKeys: []int{3, 1},
			expectedLen:  2,
		},
		{
			name:     "update prevents eviction",
			capacity: 2,
			operations: []string{
				"put:1:100",
				"put:2:200",
				"put:1:150", // Update makes 1 most recently used
				"put:3:300", // Should evict 2
			},
			expectedKeys: []int{3, 1},
			expectedLen:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[int, int](tt.capacity)

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

			if cache.Len() != tt.expectedLen {
				t.Errorf("Eviction test cache length = %d; expected %d", cache.Len(), tt.expectedLen)
			}

			keys := cache.Keys()
			if len(keys) != len(tt.expectedKeys) {
				t.Errorf("Eviction test keys length = %d; expected %d", len(keys), len(tt.expectedKeys))
				return
			}
			for i, key := range keys {
				if key != tt.expectedKeys[i] {
					t.Errorf("Eviction test keys[%d] = %d; expected %d", i, key, tt.expectedKeys[i])
				}
			}
		})
	}
}

// TestLRUCacheGenericTypes tests the cache with different types
func TestLRUCacheGenericTypes(t *testing.T) {
	t.Run("string key string value", func(t *testing.T) {
		cache := NewLRUCache[string, string](3)
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

		cache := NewLRUCache[string, User](2)
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
		cache := NewLRUCache[int, *int](2)
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

// TestLRUCacheCapacityOne tests edge case with capacity of 1
func TestLRUCacheCapacityOne(t *testing.T) {
	cache := NewLRUCache[int, int](1)

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

// Helper function to parse operation strings
func parseOperation(op string, key *int, value *int) (string, error) {
	var opType string
	if op[:3] == "put" {
		opType = "put"
		_, err := fmt.Sscanf(op, "put:%d:%d", key, value)
		return opType, err
	} else if op[:3] == "get" {
		opType = "get"
		_, err := fmt.Sscanf(op, "get:%d", key)
		return opType, err
	}
	return "", fmt.Errorf("unknown operation")
}
