package ratelimiter

import (
	"context"
	"fmt"

	"github.com/kaldun/go-rate-limiting/storage"
)

// SlidingWindow implements the sliding window log algorithm
// Keeps timestamp of each request and counts requests in last N seconds
// Most accurate but memory intensive - O(N) where N is requests in window
type SlidingWindow struct {
	storage storage.Storage
	config  Config
}

// NewSlidingWindow creates a new sliding window rate limiter
func NewSlidingWindow(storage storage.Storage, config Config) *SlidingWindow {
	return &SlidingWindow{
		storage: storage,
		config:  config,
	}
}

// Allow checks if a request should be allowed
func (sw *SlidingWindow) Allow(ctx context.Context, key string) (bool, error) {
	// TODO: Implement sliding window log algorithm
	// Steps:
	// 1. Get current timestamp
	// 2. Calculate window start time (now - window duration)
	// 3. Remove timestamps older than window start from storage
	// 4. Count remaining timestamps
	// 5. If count < rate limit, add new timestamp
	// 6. Return result
	//
	// Note: Use Redis Sorted Set (ZSET) for efficient timestamp storage
	// - Score = timestamp
	// - ZREMRANGEBYSCORE to remove old entries
	// - ZCARD to count entries
	// - ZADD to add new timestamp

	return false, fmt.Errorf("not implemented")
}

// AllowN checks if N requests should be allowed
func (sw *SlidingWindow) AllowN(ctx context.Context, key string, n int) (bool, error) {
	// TODO: Implement AllowN for sliding window
	// Add N timestamps instead of 1

	return false, fmt.Errorf("not implemented")
}

// Reset clears the rate limit state for a key
func (sw *SlidingWindow) Reset(ctx context.Context, key string) error {
	return sw.storage.Delete(ctx, key)
}

// AllowWithInfo returns detailed information about the rate limit check
func (sw *SlidingWindow) AllowWithInfo(ctx context.Context, key string) (*Result, error) {
	// TODO: Implement AllowWithInfo
	// Return Result with Allowed, Remaining, RetryAfter, ResetAt fields

	return nil, fmt.Errorf("not implemented")
}
