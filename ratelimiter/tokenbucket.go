package ratelimiter

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/kaldun-tech/go-rate-limiting/storage"
)

const TOKEN_SUFFIX = "tokens"
const LAST_REFILL_SUFFIX = ":last_refill"

// TokenBucket implements the token bucket algorithm
// Tokens are added at a constant rate, requests consume tokens
// Allows bursts up to bucket capacity
//
// Key concept: Each unique "key" gets its own independent token bucket.
// The key identifies WHO is being rate limited (user ID, API key, IP address, etc.)
//
// Storage pattern per key:
//   - key + ":tokens"      -> current token count (float64 as string)
//   - key + ":last_refill" -> last refill timestamp (Unix seconds as string)
//
// Examples:
//   - User rate limiting:  key = "user:alice"
//   - API key limiting:    key = "api_key:abc123"
//   - IP-based limiting:   key = "ip:192.168.1.1"
//
// Each key maintains independent state - users don't share buckets.
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

// Allow checks if a request should be allowed for the given key
// The key identifies the entity being rate limited (user ID, API key, IP, etc.)
// Returns true if the request is allowed, false if rate limit exceeded
func (tb *TokenBucket) Allow(ctx context.Context, key string) (bool, error) {
	// Get current tokens and last refill time from storage
	curTokens, err := tb.storage.Get(ctx, key+TOKEN_SUFFIX)
	if err != nil {
		return false, err
	}
	s, err := tb.storage.Get(ctx, key+LAST_REFILL_SUFFIX)
	if err != nil {
		return false, err
	}
	// Parse Unix timestamp
	unixTime, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false, err
	}
	// Convert to time.Time
	lastRefillTime := time.Unix(unixTime, 0)

	// Calculate tokens to add based on time elapsed
	elapsed := time.Since(lastRefillTime)
	refillRate := float64(tb.config.Rate) / tb.config.Window.Seconds()
	tokensToAdd := refillRate * elapsed.Seconds()

	// 3. Refill tokens (up to bucket capacity = burst size)
	newTokens := curTokens + tokensToAdd
	finalTokens := math.Min(newTokens, float64(tb.config.BurstSize))

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
