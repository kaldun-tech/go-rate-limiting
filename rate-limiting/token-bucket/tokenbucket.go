package tokenbucket

import (
	"math"
	"sync"
	"time"
)

// TokenBucket implements the token bucket algorithm
// Tokens are added at a constant rate, requests consume tokens
// Allows bursts up to bucket capacity
//
// Key concept: Each unique "key" gets its own independent token bucket.
// The key identifies WHO is being rate limited (user ID, API key, IP address, etc.)
// Each key maintains independent state - users don't share buckets.
type TokenBucket struct {
	mu        sync.Mutex
	buckets   map[string]*bucket // Key is bucket state
	rate      int                // Tokens per window
	burstSize int                // Max tokens in bucket
	window    time.Duration      // Time window for rate
}

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate int, window time.Duration, burstSize int) *TokenBucket {
	// Default burst size to rate if not specified
	if burstSize == 0 {
		burstSize = rate
	}
	return &TokenBucket{
		buckets:   make(map[string]*bucket),
		rate:      rate,
		burstSize: burstSize,
		window:    window,
	}
}

// Allow checks if a request should be allowed for the given key using one token
// The key identifies the entity being rate limited (user ID, API key, IP, etc.)
// Returns true if the request is allowed, false if rate limit exceeded
func (tb *TokenBucket) Allow(key string) bool {
	return tb.AllowN(key, 1)
}

// AllowN checks if N requests should be allowed
func (tb *TokenBucket) AllowN(key string, n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Get current tokens and last refill time from storage
	b, exists := tb.buckets[key]
	if !exists {
		// Check that n is in bounds of burstSize
		if n <= tb.burstSize {
			// First request for key -> create bucket with full tokens minus n
			tb.buckets[key] = &bucket{
				tokens:     float64(tb.burstSize - n),
				lastRefill: time.Now(),
			}
			return true
		}
		// Requested more than bucket capacity
		return false
	}

	// Calculate tokens to add based on time elapsed
	elapsed := time.Since(b.lastRefill)
	refillRate := float64(tb.rate) / tb.window.Seconds()
	tokensToAdd := refillRate * elapsed.Seconds()

	// Refill tokens (up to bucket capacity = burst size)
	b.tokens = math.Min(b.tokens+tokensToAdd, float64(tb.burstSize))
	b.lastRefill = time.Now()

	// Check and consume
	if float64(n) <= b.tokens {
		b.tokens -= float64(n)
		return true
	}

	// Not enough tokens
	return false
}

// Reset clears the rate limit state for a key
func (tb *TokenBucket) Reset(key string) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	delete(tb.buckets, key)
}

// Build a result given a bucket, whether it is allowed, and requested tokens
func (tb *TokenBucket) buildResult(b *bucket, allowed bool, requestedTokens int) *Result {
	if b == nil {
		// No bucket!
		return &Result{Allowed: false, Remaining: 0}
	}

	result := &Result{
		Allowed:   allowed,
		Remaining: int(b.tokens),
	}
	refillRate := float64(tb.rate) / tb.window.Seconds()

	// Calculate RetryAfter if denied
	if !allowed {
		tokensNeeded := float64(requestedTokens) - b.tokens
		result.RetryAfter = time.Duration(tokensNeeded * float64(time.Second) / refillRate)
	}

	// Calculate ResetAt = when bucket is full
	tokensToFull := float64(tb.burstSize) - b.tokens
	result.ResetAt = time.Now().Add(time.Duration(tokensToFull * float64(time.Second) / refillRate))

	return result
}

// AllowWithInfo returns detailed information about the rate limit check
// Return Result with Allowed, Remaining, RetryAfter, ResetAt fields
func (tb *TokenBucket) AllowWithInfo(key string, n int) *Result {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Get current tokens and last refill time from storage
	b, exists := tb.buckets[key]
	if !exists {
		// Check that n is in bounds of burstSize
		allowed := false
		if n <= tb.burstSize {
			// First request for key -> create bucket with full tokens minus n
			tb.buckets[key] = &bucket{
				tokens:     float64(tb.burstSize - n),
				lastRefill: time.Now(),
			}
			allowed = true
		}

		return tb.buildResult(tb.buckets[key], allowed, n)
	}

	// Calculate tokens to add based on time elapsed
	elapsed := time.Since(b.lastRefill)
	refillRate := float64(tb.rate) / tb.window.Seconds()
	tokensToAdd := refillRate * elapsed.Seconds()

	// Refill tokens (up to bucket capacity = burst size)
	b.tokens = math.Min(b.tokens+tokensToAdd, float64(tb.burstSize))
	b.lastRefill = time.Now()

	// Check and consume
	allowed := false
	if float64(n) <= b.tokens {
		b.tokens -= float64(n)
		allowed = true
	}

	return tb.buildResult(b, allowed, n)
}
