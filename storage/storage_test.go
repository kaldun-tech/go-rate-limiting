package storage

import (
	"context"
	"testing"
	"time"
)

// TestMemoryStorage tests the in-memory storage implementation
func TestMemoryStorage(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()

	t.Run("get and set", func(t *testing.T) {
		t.Skip("TODO: Implement test")
		// Test basic get/set operations
	})

	t.Run("increment", func(t *testing.T) {
		t.Skip("TODO: Implement test")
		// Test increment operations
	})

	t.Run("expiration", func(t *testing.T) {
		t.Skip("TODO: Implement test")
		// Test that keys expire correctly
	})

	t.Run("sorted sets", func(t *testing.T) {
		t.Skip("TODO: Implement test")
		// Test ZAdd, ZRemRangeByScore, ZCount, ZCard
	})

	_ = store
	_ = ctx
}

// TestRedisStorage tests the Redis storage implementation
// Requires Redis to be running on localhost:6379
func TestRedisStorage(t *testing.T) {
	// TODO: Implement Redis storage tests
	// Use t.Skip() if Redis is not available
	t.Skip("TODO: Implement test - requires Redis")
}

// BenchmarkMemoryStorage benchmarks the in-memory storage
func BenchmarkMemoryStorage(b *testing.B) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()
	key := "bench:key"

	b.Run("Get", func(b *testing.B) {
		store.Set(ctx, key, "value", time.Minute)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = store.Get(ctx, key)
		}
	})

	b.Run("Set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = store.Set(ctx, key, "value", time.Minute)
		}
	})

	b.Run("Increment", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = store.Increment(ctx, key)
		}
	})
}
