package datastructures

// LRUCache implements a Least Recently Used cache with O(1) get and put operations.
// This is a common interview question that tests understanding of hash maps and doubly linked lists.
//
// Typical interview requirements:
// - Get(key) - Get value from cache, return -1 if not exists
// - Put(key, value) - Insert or update key-value pair
// - When cache reaches capacity, evict least recently used item before inserting new item
// - Both operations should be O(1)
//
// Implementation approach:
// - HashMap for O(1) lookup
// - Doubly linked list to track recency (most recent at head, least recent at tail)
// - Move nodes to head on access
// - Evict from tail when full
type LRUCache struct {
	capacity int
	// TODO: Add fields for HashMap and doubly linked list
}

// NewLRUCache creates a new LRU cache with given capacity
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		// TODO: Initialize data structures
	}
}

// Get retrieves a value from the cache
// Returns the value if key exists, -1 otherwise
// Marks the key as recently used
func (c *LRUCache) Get(key int) int {
	// TODO: Implement
	return -1
}

// Put inserts or updates a key-value pair
// Evicts least recently used item if at capacity
func (c *LRUCache) Put(key, value int) {
	// TODO: Implement
}
