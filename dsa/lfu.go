package dsa

import "sync"

// LFUEntry represents a key-value pair with frequency count in the LFU cache
type LFUEntry[K comparable, V any] struct {
	Key       K
	Value     V
	Frequency int
}

// LFUCache is a Least Frequently Used cache with generic key and value types
type LFUCache[K comparable, V any] struct {
	capacity     int
	minFrequency int
	cache        map[K]*DLLNode[LFUEntry[K, V]]
	frequencies  map[int]*DLList[LFUEntry[K, V]]
	sync.RWMutex
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	if capacity <= 0 {
		panic("LFU cache capacity must be greater than 0")
	}
	return &LFUCache[K, V]{
		capacity:     capacity,
		minFrequency: 0,
		cache:        make(map[K]*DLLNode[LFUEntry[K, V]], capacity),
		frequencies:  make(map[int]*DLList[LFUEntry[K, V]]),
	}
}

// Get retrieves a value from the cache by key
// Returns the value and true if found, zero value and false otherwise
func (c *LFUCache[K, V]) Get(key K) (V, bool) {
	c.RLock()
	node, found := c.cache[key]
	c.RUnlock()

	if !found {
		var zero V
		return zero, false
	}

	c.Lock()
	defer c.Unlock()

	// Remove from current frequency list
	oldFreq := node.Value.Frequency
	c.removeFromFrequency(node, oldFreq)

	// Increment frequency
	node.Value.Frequency++
	newFreq := node.Value.Frequency

	// Add to new frequency list
	c.addToFrequency(node, newFreq)

	// Update minFrequency if needed
	if c.frequencies[oldFreq].Len() == 0 {
		delete(c.frequencies, oldFreq)
		if oldFreq == c.minFrequency {
			c.minFrequency++
		}
	}

	return node.Value.Value, true
}

// Put adds or updates a key-value pair in the cache
func (c *LFUCache[K, V]) Put(key K, value V) {
	c.Lock()
	defer c.Unlock()

	// If key exists, update value and frequency
	if node, found := c.cache[key]; found {
		oldFreq := node.Value.Frequency
		node.Value.Value = value

		// Remove from current frequency list
		c.removeFromFrequency(node, oldFreq)

		// Increment frequency
		node.Value.Frequency++
		newFreq := node.Value.Frequency

		// Add to new frequency list
		c.addToFrequency(node, newFreq)

		// Update minFrequency if needed
		if c.frequencies[oldFreq].Len() == 0 {
			delete(c.frequencies, oldFreq)
			if oldFreq == c.minFrequency {
				c.minFrequency++
			}
		}
		return
	}

	// If at capacity, remove least frequently used (and least recently used if tied)
	if len(c.cache) >= c.capacity {
		c.evict()
	}

	// Add new entry with frequency 1
	entry := LFUEntry[K, V]{
		Key:       key,
		Value:     value,
		Frequency: 1,
	}

	// Ensure frequency list exists
	if c.frequencies[1] == nil {
		c.frequencies[1] = NewDLList[LFUEntry[K, V]]()
	}

	// Add to front of frequency 1 list (most recent)
	node := c.frequencies[1].PushFront(entry)
	c.cache[key] = node
	c.minFrequency = 1
}

// evict removes the least frequently used item (and least recently used if tied)
func (c *LFUCache[K, V]) evict() {
	if c.minFrequency == 0 || c.frequencies[c.minFrequency] == nil {
		return
	}

	// Remove from back of minimum frequency list (least recently used)
	freqList := c.frequencies[c.minFrequency]
	if freqList.Len() == 0 {
		return
	}

	back := freqList.Back()
	if back != nil {
		freqList.Remove(back)
		delete(c.cache, back.Value.Key)
		if freqList.Len() == 0 {
			delete(c.frequencies, c.minFrequency)
			c.minFrequency++
			return
		}
	}
}

// removeFromFrequency removes a node from its frequency list
func (c *LFUCache[K, V]) removeFromFrequency(node *DLLNode[LFUEntry[K, V]], freq int) {
	if list, exists := c.frequencies[freq]; exists {
		list.Remove(node)
	}
}

// addToFrequency adds a node to the front of its frequency list
func (c *LFUCache[K, V]) addToFrequency(node *DLLNode[LFUEntry[K, V]], freq int) {
	if c.frequencies[freq] == nil {
		c.frequencies[freq] = NewDLList[LFUEntry[K, V]]()
	}
	c.frequencies[freq].PushFront(node.Value)
	c.cache[node.Value.Key] = c.frequencies[freq].Front()
}

// Remove removes a key-value pair from the cache
// Returns true if the key was found and removed, false otherwise
func (c *LFUCache[K, V]) Remove(key K) bool {
	c.Lock()
	defer c.Unlock()

	if node, found := c.cache[key]; found {
		freq := node.Value.Frequency
		c.removeFromFrequency(node, freq)
		delete(c.cache, key)

		// Update minFrequency if needed
		if freq == c.minFrequency && c.frequencies[freq].Len() == 0 {
			// Find next minimum frequency
			c.updateMinFrequency()
		}
		return true
	}
	return false
}

// updateMinFrequency finds the next minimum frequency
func (c *LFUCache[K, V]) updateMinFrequency() {
	if len(c.cache) == 0 {
		c.minFrequency = 0
		return
	}

	// Find the smallest frequency that has items
	for freq := c.minFrequency; freq <= c.capacity; freq++ {
		if list, exists := c.frequencies[freq]; exists && list.Len() > 0 {
			c.minFrequency = freq
			return
		}
	}
	c.minFrequency = 0
}

// Len returns the current number of items in the cache
func (c *LFUCache[K, V]) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.cache)
}

// Clear removes all items from the cache
func (c *LFUCache[K, V]) Clear() {
	c.Lock()
	defer c.Unlock()
	c.cache = make(map[K]*DLLNode[LFUEntry[K, V]], c.capacity)
	c.frequencies = make(map[int]*DLList[LFUEntry[K, V]])
	c.minFrequency = 0
}

// Capacity returns the maximum capacity of the cache
func (c *LFUCache[K, V]) Capacity() int {
	return c.capacity
}

// Contains checks if a key exists in the cache without updating frequency
func (c *LFUCache[K, V]) Contains(key K) bool {
	c.RLock()
	defer c.RUnlock()
	_, found := c.cache[key]
	return found
}

// Peek retrieves a value without updating the frequency
// Returns the value and true if found, zero value and false otherwise
func (c *LFUCache[K, V]) Peek(key K) (V, bool) {
	c.RLock()
	defer c.RUnlock()

	if node, found := c.cache[key]; found {
		return node.Value.Value, true
	}
	var zero V
	return zero, false
}

// GetFrequency returns the access frequency of a key
// Returns the frequency and true if found, 0 and false otherwise
func (c *LFUCache[K, V]) GetFrequency(key K) (int, bool) {
	c.RLock()
	defer c.RUnlock()

	if node, found := c.cache[key]; found {
		return node.Value.Frequency, true
	}
	return 0, false
}

// Keys returns all keys in the cache
func (c *LFUCache[K, V]) Keys() []K {
	c.RLock()
	defer c.RUnlock()

	keys := make([]K, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys
}
