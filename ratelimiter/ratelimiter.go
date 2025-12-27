package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter defines the interface that all rate limiting algorithms must implement
//
// Key concept: The "key" parameter identifies WHO is being rate limited.
// Each unique key gets its own independent rate limit state.
//
// Common key patterns:
//   - User-based:  "user:alice", "user:bob"
//   - API key:     "api_key:abc123"
//   - IP-based:    "ip:192.168.1.1"
//   - Combined:    "api_key:abc123:endpoint:/create"
//
// Each key maintains separate state, so different users/IPs/API keys
// don't interfere with each other's rate limits.
type RateLimiter interface {
	// Allow checks if a request should be allowed for the given key
	// Returns true if allowed, false if rate limit exceeded
	Allow(ctx context.Context, key string) (bool, error)

	// AllowN checks if N requests should be allowed for the given key
	// Useful for batch operations or requests with different costs
	AllowN(ctx context.Context, key string, n int) (bool, error)

	// Reset clears the rate limit state for a given key
	Reset(ctx context.Context, key string) error
}

// Config holds the configuration for rate limiters
type Config struct {
	// Rate is the number of requests allowed per window
	Rate int

	// Window is the time window for the rate limit
	Window time.Duration

	// BurstSize is the maximum burst size (only for Token Bucket)
	// If not set, defaults to Rate
	BurstSize int
}

// Result contains information about a rate limit check
type Result struct {
	// Allowed indicates if the request is allowed
	Allowed bool

	// Remaining is the number of requests remaining in the current window
	Remaining int

	// RetryAfter is the duration to wait before retrying (if not allowed)
	RetryAfter time.Duration

	// ResetAt is when the rate limit will reset
	ResetAt time.Time
}

// RateLimiterWithInfo is an extended interface that provides detailed information
type RateLimiterWithInfo interface {
	RateLimiter

	// AllowWithInfo returns detailed information about the rate limit check
	AllowWithInfo(ctx context.Context, key string) (*Result, error)
}
