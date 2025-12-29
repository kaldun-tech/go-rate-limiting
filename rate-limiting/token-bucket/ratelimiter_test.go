package tokenbucket

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_Allow(t *testing.T) {
	limiter := NewTokenBucket(10, time.Second, 0)
	key := "user:alice"

	// First 10 requests should succeed (bucket starts full)
	for i := 0; i < 10; i++ {
		if !limiter.Allow(key) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 11th request should fail (bucket empty)
	if limiter.Allow(key) {
		t.Error("Request 11 should be denied - bucket is empty")
	}
}

func TestTokenBucket_AllowN(t *testing.T) {
	limiter := NewTokenBucket(100, time.Minute, 100)
	key := "user:bob"

	// Single request (1 token)
	if !limiter.AllowN(key, 1) {
		t.Error("Single request should be allowed")
	}

	// Batch operation (10 tokens)
	if !limiter.AllowN(key, 10) {
		t.Error("Batch request (10 tokens) should be allowed")
	}

	// Expensive operation (50 tokens)
	if !limiter.AllowN(key, 50) {
		t.Error("Expensive request (50 tokens) should be allowed")
	}

	// Check remaining: started with 100, used 1 + 10 + 50 = 61, should have 39 left
	// Another expensive operation should fail
	if limiter.AllowN(key, 50) {
		t.Error("Second expensive request should be denied (only ~39 tokens remaining)")
	}
}

func TestTokenBucket_ExceedsCapacity(t *testing.T) {
	limiter := NewTokenBucket(10, time.Second, 10)
	key := "user:charlie"

	// Request more than bucket capacity
	if limiter.AllowN(key, 20) {
		t.Error("Request exceeding bucket capacity should be denied")
	}

	// Bucket should still be full for normal requests
	if !limiter.Allow(key) {
		t.Error("Normal request should still be allowed")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	// 10 tokens/second
	limiter := NewTokenBucket(10, time.Second, 10)
	key := "user:david"

	// Drain the bucket
	for i := 0; i < 10; i++ {
		limiter.Allow(key)
	}

	// Should be denied now
	if limiter.Allow(key) {
		t.Error("Should be denied after draining bucket")
	}

	// Wait 500ms (should refill ~5 tokens)
	time.Sleep(500 * time.Millisecond)

	// Should be able to make a few more requests
	allowed := 0
	for i := 0; i < 10; i++ {
		if limiter.Allow(key) {
			allowed++
		}
	}

	// Should have gotten approximately 5 tokens back (±1 for timing variance)
	if allowed < 4 || allowed > 6 {
		t.Errorf("Expected ~5 tokens after 500ms, got %d", allowed)
	}
}

func TestTokenBucket_IndependentKeys(t *testing.T) {
	limiter := NewTokenBucket(5, time.Second, 5)

	// Drain Alice's bucket
	for i := 0; i < 5; i++ {
		limiter.Allow("user:alice")
	}

	// Alice should be denied
	if limiter.Allow("user:alice") {
		t.Error("Alice should be rate limited")
	}

	// Bob should still have full bucket (independent keys)
	if !limiter.Allow("user:bob") {
		t.Error("Bob should not be affected by Alice's rate limit")
	}
}

func TestTokenBucket_AllowWithInfo(t *testing.T) {
	limiter := NewTokenBucket(10, time.Second, 15)
	key := "user:eve"

	// First request
	result := limiter.AllowWithInfo(key, 1)
	if !result.Allowed {
		t.Error("First request should be allowed")
	}
	if result.Remaining != 14 {
		t.Errorf("Expected 14 remaining tokens, got %d", result.Remaining)
	}

	// Drain the bucket
	for limiter.Allow(key) {
		// Keep draining
	}

	// Check denied result
	result = limiter.AllowWithInfo(key, 1)
	if result.Allowed {
		t.Error("Request should be denied when bucket is empty")
	}
	if result.Remaining != 0 {
		t.Errorf("Expected 0 remaining tokens, got %d", result.Remaining)
	}
	if result.RetryAfter == 0 {
		t.Error("RetryAfter should be set when denied")
	}
	if result.ResetAt.IsZero() {
		t.Error("ResetAt should be set")
	}
}

func TestTokenBucket_Burst(t *testing.T) {
	// Rate: 5/sec, Burst: 10 (allows temporary spike)
	limiter := NewTokenBucket(5, time.Second, 10)
	key := "user:frank"

	// Should handle burst of 10
	for i := 0; i < 10; i++ {
		if !limiter.Allow(key) {
			t.Errorf("Burst request %d should be allowed (burst capacity = 10)", i+1)
		}
	}

	// 11th should fail
	if limiter.Allow(key) {
		t.Error("Request beyond burst capacity should be denied")
	}
}

func TestTokenBucket_Reset(t *testing.T) {
	limiter := NewTokenBucket(5, time.Second, 5)
	key := "user:grace"

	// Drain bucket
	for i := 0; i < 5; i++ {
		limiter.Allow(key)
	}

	// Should be denied
	if limiter.Allow(key) {
		t.Error("Should be denied after draining")
	}

	// Reset
	limiter.Reset(key)

	// Should have full bucket again
	if !limiter.Allow(key) {
		t.Error("Should be allowed after reset")
	}
}

func TestTokenBucket_ConcurrentAccess(t *testing.T) {
	limiter := NewTokenBucket(100, time.Second, 100)
	key := "user:concurrent"

	var wg sync.WaitGroup
	allowedCount := 0
	var mu sync.Mutex

	// Launch 150 concurrent requests (more than capacity)
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow(key) {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should have allowed exactly 100 (bucket capacity)
	if allowedCount != 100 {
		t.Errorf("Expected exactly 100 allowed requests, got %d", allowedCount)
	}
}

func BenchmarkTokenBucket_Allow(b *testing.B) {
	limiter := NewTokenBucket(1000000, time.Second, 1000000)
	key := "bench:user:1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(key)
	}
}

func BenchmarkTokenBucket_AllowN(b *testing.B) {
	limiter := NewTokenBucket(1000000, time.Second, 1000000)
	key := "bench:user:2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.AllowN(key, 5)
	}
}

func BenchmarkTokenBucket_AllowWithInfo(b *testing.B) {
	limiter := NewTokenBucket(1000000, time.Second, 1000000)
	key := "bench:user:3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.AllowWithInfo(key, 1)
	}
}

func BenchmarkTokenBucket_Concurrent(b *testing.B) {
	limiter := NewTokenBucket(1000000, time.Second, 1000000)
	keys := []string{"user:1", "user:2", "user:3", "user:4"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%len(keys)]
			limiter.Allow(key)
			i++
		}
	})
}
