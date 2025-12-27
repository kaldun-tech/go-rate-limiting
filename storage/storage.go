package storage

import (
	"context"
	"time"
)

// Storage defines the interface for storing rate limit state
// Implementations can use Redis, in-memory cache, or other backends
//
// Note: Storage is a generic key-value store. Rate limiters compose keys
// from the user-provided key (e.g., "user:alice") plus algorithm-specific
// suffixes (e.g., "user:alice:tokens", "user:alice:last_refill").
type Storage interface {
	// Get retrieves a value from storage
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value with optional expiration
	Set(ctx context.Context, key string, value string, expiration time.Duration) error

	// Increment atomically increments a counter and returns the new value
	// If the key doesn't exist, it should be created with value 1
	Increment(ctx context.Context, key string) (int64, error)

	// IncrementBy atomically increments a counter by n and returns the new value
	IncrementBy(ctx context.Context, key string, n int64) (int64, error)

	// Delete removes a key from storage
	Delete(ctx context.Context, key string) error

	// Expire sets an expiration on a key
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// GetMultiple retrieves multiple values in a single operation
	GetMultiple(ctx context.Context, keys []string) (map[string]string, error)

	// SetMultiple stores multiple key-value pairs in a single operation
	SetMultiple(ctx context.Context, items map[string]string, expiration time.Duration) error

	// Close closes the storage connection
	Close() error
}

// SortedSetStorage extends Storage with sorted set operations
// Useful for sliding window log implementation
type SortedSetStorage interface {
	Storage

	// ZAdd adds members with scores to a sorted set
	ZAdd(ctx context.Context, key string, members ...SortedSetMember) error

	// ZRemRangeByScore removes members with scores in the given range
	ZRemRangeByScore(ctx context.Context, key string, min, max float64) error

	// ZCount counts members with scores in the given range
	ZCount(ctx context.Context, key string, min, max float64) (int64, error)

	// ZCard returns the number of members in a sorted set
	ZCard(ctx context.Context, key string) (int64, error)
}

// SortedSetMember represents a member in a sorted set
type SortedSetMember struct {
	Score  float64
	Member string
}

// ScriptStorage extends Storage with Lua script execution
// Useful for atomic operations in Redis
type ScriptStorage interface {
	Storage

	// Eval executes a Lua script
	Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error)
}
