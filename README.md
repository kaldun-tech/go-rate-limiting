# Token Bucket Rate Limiter in Go

A clean, thread-safe implementation of the Token Bucket rate limiting algorithm in Go. Built as interview preparation to demonstrate understanding of:
- Concurrent programming (mutex, thread safety)
- Algorithm implementation (token bucket)
- System design (rate limiting patterns)
- Clean API design

## What is Token Bucket?

The Token Bucket algorithm is one of the most popular rate limiting algorithms used by APIs like Stripe, GitHub, and AWS.

**How it works:**
1. A "bucket" holds tokens (starts full at `burstSize`)
2. Tokens refill at a constant rate (`rate` tokens per `window`)
3. Each request consumes tokens (usually 1 per request)
4. Request allowed if enough tokens available
5. Request rejected if bucket empty

**Key feature:** Allows traffic bursts up to bucket capacity while maintaining average rate limit.

## Installation

```bash
go get github.com/kaldun-tech/go-rate-limiting
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/kaldun-tech/go-rate-limiting/ratelimiter"
)

func main() {
    // Create limiter: 10 requests/second, burst up to 10
    limiter := ratelimiter.NewTokenBucket(10, time.Second, 0)

    // Check if request is allowed
    if limiter.Allow("user:alice") {
        fmt.Println("Request allowed!")
    } else {
        fmt.Println("Rate limited - please slow down")
    }
}
```

## API Reference

### Creating a Limiter

```go
func NewTokenBucket(rate int, window time.Duration, burstSize int) *TokenBucket
```

**Parameters:**
- `rate`: Number of tokens added per window (e.g., 100)
- `window`: Time window for rate (e.g., `time.Minute`)
- `burstSize`: Maximum tokens in bucket (e.g., 150). Set to `0` to default to `rate`

**Examples:**
```go
// 10 requests/second, burst=10
limiter := ratelimiter.NewTokenBucket(10, time.Second, 0)

// 100 requests/minute, burst=150 (allows temporary spike)
limiter := ratelimiter.NewTokenBucket(100, time.Minute, 150)

// 1000 requests/hour
limiter := ratelimiter.NewTokenBucket(1000, time.Hour, 0)
```

### Checking Rate Limits

#### Allow(key string) bool

Check if a single request should be allowed.

```go
allowed := limiter.Allow("user:123")
if !allowed {
    // Return HTTP 429 Too Many Requests
}
```

**The "key" identifies WHO is being rate limited:**
- `"user:alice"` - Per-user rate limiting
- `"api_key:abc123"` - Per API key
- `"ip:192.168.1.1"` - Per IP address

Each key maintains independent bucket state - users don't affect each other's limits!

#### AllowN(key string, n int) bool

Allow N tokens at once (for batch operations or weighted costs).

```go
// Normal request costs 1 token
allowed := limiter.AllowN("user:123", 1)

// Batch upload costs 10 tokens
allowed = limiter.AllowN("user:123", 10)

// Expensive AI query costs 50 tokens
allowed = limiter.AllowN("user:123", 50)
```

#### AllowWithInfo(key string, n int) *Result

Get detailed information about rate limit check.

```go
result := limiter.AllowWithInfo("user:123", 1)

fmt.Printf("Allowed: %v\n", result.Allowed)
fmt.Printf("Remaining: %d tokens\n", result.Remaining)
fmt.Printf("Retry after: %v\n", result.RetryAfter)  // If denied
fmt.Printf("Resets at: %v\n", result.ResetAt)
```

**Result fields:**
```go
type Result struct {
    Allowed    bool          // Whether request was allowed
    Remaining  int           // Tokens remaining in bucket
    RetryAfter time.Duration // How long to wait if denied
    ResetAt    time.Time     // When bucket refills to full
}
```

Use this to set HTTP headers:
```go
result := limiter.AllowWithInfo(apiKey, 1)
w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))
if !result.Allowed {
    w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
    http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
    return
}
```

#### Reset(key string)

Clear rate limit state for a key (useful for testing or admin overrides).

```go
limiter.Reset("user:123")
```

## Examples

Run the included examples:

```bash
go run examples/main.go
```

The examples demonstrate:
1. Basic rate limiting (allow/deny pattern)
2. Burst capacity (handling traffic spikes)
3. Detailed information (`AllowWithInfo`)
4. Weighted costs (`AllowN`)

## Understanding Keys

**Each unique key gets its own independent token bucket.**

```go
limiter := ratelimiter.NewTokenBucket(10, time.Second, 0)

// Alice gets 10 req/sec
limiter.Allow("user:alice")
limiter.Allow("user:alice")

// Bob ALSO gets 10 req/sec (separate bucket!)
limiter.Allow("user:bob")

// Can combine factors for finer control
limiter.Allow("user:alice:endpoint:/api/create")
limiter.Allow("user:alice:endpoint:/api/read")
```

**Common key patterns:**
- `"user:{userID}"` - Per-user limits (best for authenticated APIs)
- `"api_key:{key}"` - Per API key (common for public APIs)
- `"ip:{address}"` - Per IP (fallback for unauthenticated endpoints)
- `"user:{userID}:resource:{id}"` - Per-user per-resource

## Thread Safety

The implementation is **thread-safe** and can be called concurrently from multiple goroutines:

```go
limiter := ratelimiter.NewTokenBucket(100, time.Second, 0)

// Safe to call from multiple goroutines
go func() { limiter.Allow("user:alice") }()
go func() { limiter.Allow("user:bob") }()
go func() { limiter.Allow("user:alice") }()
```

Internally uses `sync.Mutex` to protect shared state.

## Algorithm Details

**Refill calculation:**
```
refillRate = rate / window.Seconds()
tokensToAdd = refillRate * elapsedTime.Seconds()
newTokens = min(currentTokens + tokensToAdd, burstSize)
```

**Example:** Rate=100/min, window=1min
- Refill rate: 100/60 = 1.67 tokens/second
- After 10 seconds: +16.7 tokens
- Capped at `burstSize` (won't exceed bucket capacity)

**Comparison with other algorithms:**

| Algorithm | Memory | Bursts | Accuracy | Use Case |
|-----------|--------|--------|----------|----------|
| **Token Bucket** | O(1) | ✅ Yes | Good | Most APIs (Stripe, GitHub) |
| Fixed Window | O(1) | ❌ Boundary issue | Fair | Simple counters |
| Sliding Window | O(N) | ❌ No | Excellent | Financial/security |
| Leaky Bucket | O(1) | 🔄 Smoothed | Good | Traffic shaping |

**Token Bucket is recommended for most APIs** because:
- Efficient O(1) time and space
- Allows reasonable bursts (good UX)
- Simple to implement and understand
- Industry-proven (Stripe, AWS, GitHub all use it)

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./ratelimiter

# Run specific test
go test -run TestTokenBucket ./ratelimiter

# Format and vet code
go fmt ./...
go vet ./...
```

## Project Structure

```
go-rate-limiting/
├── ratelimiter/
│   ├── ratelimiter.go     # Result type definition
│   ├── tokenbucket.go     # Token bucket implementation
│   └── ratelimiter_test.go
├── examples/
│   └── main.go            # Usage examples
├── README.md
└── go.mod
```

## Real-World Usage

**In an HTTP API:**
```go
func apiHandler(w http.ResponseWriter, r *http.Request) {
    apiKey := r.Header.Get("X-API-Key")

    result := limiter.AllowWithInfo(apiKey, 1)

    // Set rate limit headers
    w.Header().Set("X-RateLimit-Limit", "100")
    w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
    w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

    if !result.Allowed {
        w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    // Process request...
}
```

**Different costs for different operations:**
```go
// Light read: 1 token
limiter.AllowN(userID, 1)

// Batch upload: 10 tokens
limiter.AllowN(userID, 10)

// AI/ML inference: 100 tokens
limiter.AllowN(userID, 100)
```

## Production Considerations

This implementation is perfect for:
- ✅ Single-server applications
- ✅ Interview coding challenges
- ✅ Learning rate limiting algorithms
- ✅ Local development/testing

For distributed systems (multiple API servers), you'd need:
- Shared storage (Redis, Memcached)
- Atomic operations (Lua scripts in Redis)
- Consistent timestamps across servers

See [CLAUDE.md](./CLAUDE.md) for distributed system design notes.

## License

MIT

## Interview Tips

When implementing this in an interview:

1. **Start simple:** Implement `AllowN` first, then `Allow` delegates to it
2. **Explain the math:** Walk through refill calculation clearly
3. **Handle edge cases:** First request, bucket empty, n > burstSize
4. **Think about thread safety:** Explain mutex usage
5. **Discuss trade-offs:** Compare with other algorithms
6. **Extend thoughtfully:** Show how to add `AllowWithInfo` for headers

**Common interview questions:**
- *"How do you handle bursts?"* → Bucket capacity allows temporary spikes
- *"What if multiple servers?"* → Need Redis for shared state
- *"How to prevent race conditions?"* → Mutex (single server) or atomic ops (Redis)
- *"When to use this vs. fixed window?"* → Use token bucket for better UX (allows bursts)
