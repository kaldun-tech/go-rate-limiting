package ratelimiter

import (
	"context"
	"fmt"

	"github.com/kaldun/go-rate-limiting/storage"
)

// LeakyBucket implements the leaky bucket algorithm
// Requests enter a queue, processed at constant rate
// Good for smoothing traffic and enforcing steady rate
type LeakyBucket struct {
	storage storage.Storage
	config  Config
}

// NewLeakyBucket creates a new leaky bucket rate limiter
func NewLeakyBucket(storage storage.Storage, config Config) *LeakyBucket {
	return &LeakyBucket{
		storage: storage,
		config:  config,
	}
}

// Allow checks if a request should be allowed
func (lb *LeakyBucket) Allow(ctx context.Context, key string) (bool, error) {
	// TODO: Implement leaky bucket algorithm
	// Steps:
	// 1. Get current queue size and last leak time
	// 2. Calculate requests leaked since last check (based on leak rate)
	// 3. Update queue size (current - leaked)
	// 4. Check if queue has space for new request
	// 5. If yes, add request to queue
	// 6. Update storage with new queue size and leak time
	// 7. Return result

	return false, fmt.Errorf("not implemented")
}

// AllowN checks if N requests should be allowed
func (lb *LeakyBucket) AllowN(ctx context.Context, key string, n int) (bool, error) {
	// TODO: Implement AllowN for leaky bucket
	// Add N requests to queue instead of 1

	return false, fmt.Errorf("not implemented")
}

// Reset clears the rate limit state for a key
func (lb *LeakyBucket) Reset(ctx context.Context, key string) error {
	return lb.storage.Delete(ctx, key)
}

// AllowWithInfo returns detailed information about the rate limit check
func (lb *LeakyBucket) AllowWithInfo(ctx context.Context, key string) (*Result, error) {
	// TODO: Implement AllowWithInfo
	// Return Result with Allowed, Remaining, RetryAfter, ResetAt fields

	return nil, fmt.Errorf("not implemented")
}
