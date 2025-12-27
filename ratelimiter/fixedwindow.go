package ratelimiter

import (
	"context"
	"fmt"

	"github.com/kaldun-tech/go-rate-limiting/storage"
)

// FixedWindow implements the fixed window counter algorithm
// Counts requests in fixed time windows (e.g., 1:00-1:59, 2:00-2:59)
// Simple but has burst problem at window boundaries
type FixedWindow struct {
	storage storage.Storage
	config  Config
}

// NewFixedWindow creates a new fixed window rate limiter
func NewFixedWindow(storage storage.Storage, config Config) *FixedWindow {
	return &FixedWindow{
		storage: storage,
		config:  config,
	}
}

// Allow checks if a request should be allowed
func (fw *FixedWindow) Allow(ctx context.Context, key string) (bool, error) {
	// TODO: Implement fixed window algorithm
	// Steps:
	// 1. Calculate current window start time
	// 2. Create window-specific key (e.g., "user:123:window:1234567890")
	// 3. Increment counter in storage
	// 4. Set expiration on first increment
	// 5. Check if counter exceeds rate limit
	// 6. Return result

	return false, fmt.Errorf("not implemented")
}

// AllowN checks if N requests should be allowed
func (fw *FixedWindow) AllowN(ctx context.Context, key string, n int) (bool, error) {
	// TODO: Implement AllowN for fixed window
	// Similar to Allow but increment by N instead of 1

	return false, fmt.Errorf("not implemented")
}

// Reset clears the rate limit state for a key
func (fw *FixedWindow) Reset(ctx context.Context, key string) error {
	// Note: For fixed window, you may want to reset all windows for a key
	return fw.storage.Delete(ctx, key)
}

// AllowWithInfo returns detailed information about the rate limit check
func (fw *FixedWindow) AllowWithInfo(ctx context.Context, key string) (*Result, error) {
	// TODO: Implement AllowWithInfo
	// Return Result with Allowed, Remaining, RetryAfter, ResetAt fields

	return nil, fmt.Errorf("not implemented")
}
