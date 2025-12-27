package storage

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// MemoryStorage implements Storage using in-memory maps
// Useful for testing and single-server deployments
// NOT suitable for distributed systems
type MemoryStorage struct {
	mu     sync.RWMutex
	data   map[string]string
	expiry map[string]time.Time
	// For sorted sets
	sortedSets map[string][]SortedSetMember
}

// NewMemoryStorage creates a new in-memory storage backend
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		data:       make(map[string]string),
		expiry:     make(map[string]time.Time),
		sortedSets: make(map[string][]SortedSetMember),
	}

	// Start cleanup goroutine
	go storage.cleanupExpired()

	return storage
}

// cleanupExpired removes expired keys periodically
func (m *MemoryStorage) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for key, exp := range m.expiry {
			if now.After(exp) {
				delete(m.data, key)
				delete(m.expiry, key)
				delete(m.sortedSets, key)
			}
		}
		m.mu.Unlock()
	}
}

// Get retrieves a value from memory
func (m *MemoryStorage) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	// Check whether expiry is set and in the past
	if expiry, ok := m.expiry[key]; ok {
		if time.Now().After(expiry) {
			return "", fmt.Errorf("Entry expired at %v", expiry)
		}
	}
	// Data found
	return data, nil
}

// Set stores a value in memory with optional expiration
func (m *MemoryStorage) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value

	// Set expiration if specified
	if 0 < expiration {
		m.expiry[key] = time.Now().Add(expiration)
	}
	return nil
}

// Increment atomically increments a counter in memory
func (m *MemoryStorage) Increment(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Parse current value as int64
	var current int64 = 0
	if s, ok := m.data[key]; ok {
		var err error
		current, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	// Increment and store as string
	current++
	m.data[key] = strconv.FormatInt(current, 10)
	return current, nil
}

// IncrementBy atomically increments a counter by n in memory
func (m *MemoryStorage) IncrementBy(ctx context.Context, key string, n int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Parse current value as int64
	var current int64 = 0
	if s, ok := m.data[key]; ok {
		var err error
		current, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	// Increment and store as string
	current += n
	m.data[key] = strconv.FormatInt(current, 10)
	return current, nil
}

// Delete removes a key from memory
func (m *MemoryStorage) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove data, expiration, and sortedSets if present
	delete(m.data, key)
	delete(m.expiry, key)
	delete(m.sortedSets, key)
	return nil
}

// Expire sets an expiration on a key in memory
// Strict behavior: error if key doesn't exist
func (m *MemoryStorage) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("key does not exist")
	}

	m.expiry[key] = time.Now().Add(expiration)
	return nil
}

// GetMultiple retrieves multiple values from memory
func (m *MemoryStorage) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	// TODO: Implement GetMultiple
	return nil, fmt.Errorf("not implemented")
}

// SetMultiple stores multiple key-value pairs in memory
func (m *MemoryStorage) SetMultiple(ctx context.Context, items map[string]string, expiration time.Duration) error {
	// TODO: Implement SetMultiple
	return fmt.Errorf("not implemented")
}

// ZAdd adds members with scores to a sorted set
func (m *MemoryStorage) ZAdd(ctx context.Context, key string, members ...SortedSetMember) error {
	// TODO: Implement ZAdd
	// Add members to sortedSets[key] and maintain sorted order
	return fmt.Errorf("not implemented")
}

// ZRemRangeByScore removes members with scores in the given range
func (m *MemoryStorage) ZRemRangeByScore(ctx context.Context, key string, min, max float64) error {
	// TODO: Implement ZRemRangeByScore
	return fmt.Errorf("not implemented")
}

// ZCount counts members with scores in the given range
func (m *MemoryStorage) ZCount(ctx context.Context, key string, min, max float64) (int64, error) {
	// TODO: Implement ZCount
	return 0, fmt.Errorf("not implemented")
}

// ZCard returns the number of members in a sorted set
func (m *MemoryStorage) ZCard(ctx context.Context, key string) (int64, error) {
	// TODO: Implement ZCard
	return 0, fmt.Errorf("not implemented")
}

// Eval is not supported for in-memory storage
func (m *MemoryStorage) Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Eval not supported for in-memory storage")
}

// Close closes the storage (cleanup goroutine will be stopped by GC)
func (m *MemoryStorage) Close() error {
	return nil
}
