package dsa

import "sync"

// LRUEntry represents a key-value pair in the LRU cache
type LRUEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// LRUCache is a Least Recently Used cache with generic key and value types
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*DLLNode[LRUEntry[K, V]]
	list     *DLList[LRUEntry[K, V]]
	sync.RWMutex
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		panic("LRU cache capacity must be greater than 0")
	}
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*DLLNode[LRUEntry[K, V]], capacity),
		list:     NewDLList[LRUEntry[K, V]](),
	}
}

// Get retrieves a value from the cache by key
// Returns the value and true if found, zero value and false otherwise
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.RLock()
	element, found := c.cache[key]
	c.RUnlock()
	if !found {
		var zero V
		return zero, false
	}
	c.Lock()
	c.list.MoveToFront(element)
	c.Unlock()
	return element.Value.Value, true
}

// Put adds or updates a key-value pair in the cache
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.Lock()
	defer c.Unlock()
	if node, found := c.cache[key]; found {
		node.Value.Value = value
		c.list.MoveToFront(node)
		return
	}
	if c.list.Len() >= c.capacity {
		lastNode := c.list.Back()
		if lastNode != nil {
			c.list.Remove(lastNode)
			delete(c.cache, lastNode.Value.Key)
		}
	}
	entry := LRUEntry[K, V]{
		Key:   key,
		Value: value,
	}
	node := c.list.PushFront(entry)
	c.cache[key] = node
	return
}

// Remove removes a key-value pair from the cache
// Returns true if the key was found and removed, false otherwise
func (c *LRUCache[K, V]) Remove(key K) bool {
	c.Lock()
	defer c.Unlock()
	if node, found := c.cache[key]; found {
		c.list.Remove(node)
		delete(c.cache, key)
		return true
	}
	return false
}

// Len returns the current number of items in the cache
func (c *LRUCache[K, V]) Len() int {
	return c.list.Len()
}

// Clear removes all items from the cache
func (c *LRUCache[K, V]) Clear() {
	c.cache = make(map[K]*DLLNode[LRUEntry[K, V]], c.capacity)
	c.list = NewDLList[LRUEntry[K, V]]()
}

// Capacity returns the maximum capacity of the cache
func (c *LRUCache[K, V]) Capacity() int {
	return c.capacity
}

// Keys returns all keys in the cache in order from most to least recently used
func (c *LRUCache[K, V]) Keys() []K {
	keys := make([]K, 0, c.list.Len())
	for node := c.list.Front(); node != nil; node = node.Next() {
		keys = append(keys, node.Value.Key)
	}
	return keys
}

// Contains checks if a key exists in the cache without updating access order
func (c *LRUCache[K, V]) Contains(key K) bool {
	_, found := c.cache[key]
	return found
}

// Peek retrieves a value without updating the access order
// Returns the value and true if found, zero value and false otherwise
func (c *LRUCache[K, V]) Peek(key K) (V, bool) {
	if node, found := c.cache[key]; found {
		return node.Value.Value, true
	}
	var zero V
	return zero, false
}
