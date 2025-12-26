package ratelimiter

import (
	"context"
	"fmt"

	"github.com/kaldun/go-rate-limiting/storage"
)

// TokenBucket implements the token bucket algorithm
// Tokens are added at a constant rate, requests consume tokens
// Allows bursts up to bucket capacity
type TokenBucket struct {
	storage storage.Storage
	config  Config
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(storage storage.Storage, config Config) *TokenBucket {
	if config.BurstSize == 0 {
		config.BurstSize = config.Rate
	}
	return &TokenBucket{
		storage: storage,
		config:  config,
	}
}

// Allow checks if a request should be allowed
func (tb *TokenBucket) Allow(ctx context.Context, key string) (bool, error) {
	// TODO: Implement token bucket algorithm
	// Steps:
	// 1. Get current tokens and last refill time from storage
	// 2. Calculate tokens to add based on time elapsed
	// 3. Refill tokens (up to bucket capacity)
	// 4. Check if enough tokens available
	// 5. If yes, consume token and update storage
	// 6. Return result

	return false, fmt.Errorf("not implemented")
}

// AllowN checks if N requests should be allowed
func (tb *TokenBucket) AllowN(ctx context.Context, key string, n int) (bool, error) {
	// TODO: Implement AllowN for token bucket
	// Similar to Allow but consume N tokens instead of 1

	return false, fmt.Errorf("not implemented")
}

// Reset clears the rate limit state for a key
func (tb *TokenBucket) Reset(ctx context.Context, key string) error {
	return tb.storage.Delete(ctx, key)
}

// AllowWithInfo returns detailed information about the rate limit check
func (tb *TokenBucket) AllowWithInfo(ctx context.Context, key string) (*Result, error) {
	// TODO: Implement AllowWithInfo
	// Return Result with Allowed, Remaining, RetryAfter, ResetAt fields

	return nil, fmt.Errorf("not implemented")
}
