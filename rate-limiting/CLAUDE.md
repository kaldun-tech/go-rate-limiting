# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Token Bucket Rate Limiter** - A clean, interview-ready implementation of the token bucket algorithm in Go.

**Purpose**: Interview preparation for technical coding challenges. Demonstrates understanding of:
- Rate limiting algorithms (token bucket)
- Concurrent programming (mutex, thread safety)
- Clean API design (Allow, AllowN, AllowWithInfo)
- System design concepts

**Focus**: This is a **single-server, in-memory implementation**. Not production-distributed systems code.

## Key Design Decisions

### Simplified Architecture (Option B)

This implementation uses:
- **In-memory storage**: `map[string]*bucket` with `sync.Mutex`
- **No external dependencies**: No Redis, no storage abstraction layer
- **No context.Context**: Simpler API for interview setting
- **No error returns**: In-memory operations don't fail

**Why this approach?**
- Interview-friendly: Focus on algorithm, not infrastructure
- Easy to explain and code on whiteboard/CoderPad
- Shows clean Go idioms without over-engineering

### What We're NOT Building

- ❌ Distributed rate limiting (Redis/multi-server)
- ❌ Multiple algorithms (Fixed Window, Sliding Window, Leaky Bucket)
- ❌ HTTP middleware layer
- ❌ Storage abstraction interfaces

**These are valuable for production**, but add complexity that obscures the core algorithm in an interview setting.

## Implementation Details

### Token Bucket Algorithm

```go
type TokenBucket struct {
    mu        sync.Mutex
    buckets   map[string]*bucket  // key -> bucket state
    rate      int                 // tokens per window
    burstSize int                 // max tokens in bucket
    window    time.Duration       // time window for rate
}

type bucket struct {
    tokens     float64    // current tokens (can be fractional)
    lastRefill time.Time  // last refill timestamp
}
```

**Key implementation notes:**

1. **Tokens are float64**: Because refill can produce fractional tokens
   - Example: 10 tokens/min = 0.167 tokens/sec
   - After 5.5 seconds: 0.9185 tokens added

2. **Requests (n) are int**: Because you can't make 2.5 requests
   - Clean API: `AllowN(key, 5)` not `AllowN(key, 5.0)`

3. **Refill calculation**:
   ```go
   refillRate = float64(rate) / window.Seconds()
   tokensToAdd = refillRate * elapsed.Seconds()
   newTokens = math.Min(currentTokens + tokensToAdd, float64(burstSize))
   ```

4. **Allow delegates to AllowN**: DRY principle
   ```go
   func (tb *TokenBucket) Allow(key string) bool {
       return tb.AllowN(key, 1)
   }
   ```

5. **buildResult helper**: Reduces duplication between early returns and normal flow

### Edge Cases to Handle

1. **First request** (`!exists`):
   - Create bucket with `burstSize - n` tokens
   - Check `n <= burstSize` first (reject if requesting more than capacity)

2. **Bucket empty** (`tokens < n`):
   - Don't consume tokens
   - Still update `lastRefill` to current time (important!)
   - Calculate `RetryAfter` duration

3. **Request exceeds capacity** (`n > burstSize`):
   - Always deny (can never have enough tokens)

4. **Zero or negative n**:
   - Technically allowed (useful for checking state without consuming)

### Thread Safety

Uses `sync.Mutex` for thread safety:
```go
func (tb *TokenBucket) AllowN(key string, n int) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    // ... implementation
}
```

**Why mutex not RWMutex?**
- Every call modifies state (updates `lastRefill` even on deny)
- No pure read operations, so RWMutex provides no benefit

## Development Commands

### Running Examples
```bash
go run examples/main.go
```

### Testing
```bash
go test ./...                         # All tests
go test -v ./ratelimiter             # Verbose output
go test -run TestTokenBucket ./ratelimiter  # Specific test
go test -bench=. ./ratelimiter       # Benchmarks
```

### Code Quality
```bash
go fmt ./...    # Format code
go vet ./...    # Static analysis
go mod tidy     # Clean dependencies
```

## Interview Strategy

### How to Approach This in an Interview

1. **Start with signature**:
   ```go
   func (tb *TokenBucket) AllowN(key string, n int) bool
   ```

2. **Explain the data structure**:
   - "Each key gets its own bucket with tokens and lastRefill time"
   - "Tokens refill continuously based on elapsed time"

3. **Walk through algorithm**:
   - Calculate elapsed time since last refill
   - Calculate tokens to add (refillRate * elapsed)
   - Cap at burstSize
   - Check if enough tokens, consume if yes

4. **Implement, explaining as you go**

5. **Test with examples**:
   - First request (bucket creation)
   - Burst (drain bucket)
   - Refill (wait and retry)
   - Exceeding capacity

6. **Discuss extensions** (if time):
   - `AllowWithInfo` for HTTP headers
   - Distributed version with Redis
   - Alternative algorithms (Fixed Window, Sliding Window)

### Common Interview Questions

**Q: "How would this work with multiple servers?"**
A: "Current implementation is single-server. For distributed, I'd use Redis with Lua scripts for atomic operations. The algorithm is the same, but storage becomes shared state."

**Q: "What if the system clock changes?"**
A: "Good catch! Production systems use monotonic time or Redis TIME command for consistency. For this implementation, we accept system clock as-is since time-travel isn't a realistic threat in the interview context."

**Q: "Why not use a goroutine to refill tokens?"**
A: "Lazy refill is more efficient - we only calculate refill when needed, not continuously. Saves CPU and avoids timing issues."

**Q: "How do you handle different cost per request?"**
A: "`AllowN(key, n)` - some operations cost more tokens. Like 1 for reads, 10 for writes, 100 for AI inference."

**Q: "What about memory cleanup?"**
A: "Could add TTL/LRU eviction for unused keys. In practice, most systems have bounded user sets. For interview, this is acceptable unless explicitly asked."

## Testing Strategy

**What to test:**

1. ✅ Basic allow/deny
2. ✅ Burst handling (initial full bucket)
3. ✅ Refill over time
4. ✅ AllowN with different costs
5. ✅ Edge case: n > burstSize
6. ✅ Edge case: first request for key
7. ✅ Thread safety (concurrent access)
8. ✅ Result fields accuracy (Remaining, RetryAfter, ResetAt)

**Example test structure:**
```go
func TestTokenBucket_BasicAllow(t *testing.T) {
    limiter := NewTokenBucket(10, time.Second, 0)

    // First 10 should succeed
    for i := 0; i < 10; i++ {
        if !limiter.Allow("user1") {
            t.Errorf("Request %d should be allowed", i)
        }
    }

    // 11th should fail
    if limiter.Allow("user1") {
        t.Error("Request 11 should be denied")
    }
}
```

## Real-World Context (Discussion Points)

While this implementation is interview-focused, here's how production systems extend it:

### Distributed Systems

```
[Load Balancer]
    ↓
[API Gateway Cluster] (multiple instances)
    ↓
[Redis Cluster] (shared state)
```

**Changes needed:**
- Replace `map[string]*bucket` with Redis
- Use Lua scripts for atomic refill+consume
- Use Redis TIME command for consistent timestamps
- Handle Redis failures (fail open vs. fail closed)

### HTTP Integration

```go
func RateLimitMiddleware(limiter *TokenBucket) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := extractKey(r) // From API key, IP, user ID, etc.

            result := limiter.AllowWithInfo(key, 1)

            // Set headers
            w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
            w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

            if !result.Allowed {
                w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

## Common Pitfalls to Avoid

1. **Not updating lastRefill on deny**: Must update even when denying to track refill time correctly

2. **Integer division**: `rate / window` loses precision - use `float64(rate) / window.Seconds()`

3. **Forgetting to cap at burstSize**: Tokens can't exceed bucket capacity

4. **Race conditions**: All bucket access must be under mutex

5. **Not handling n > burstSize**: Should always deny if requesting more than capacity

## What Makes This Interview-Ready?

✅ **Clear structure**: Easy to explain and code step-by-step
✅ **Good defaults**: `burstSize=0` → defaults to rate
✅ **Progressive implementation**: Allow → AllowN → AllowWithInfo
✅ **Real-world applicability**: Actually used by Stripe, GitHub, AWS
✅ **Extension points**: Natural segue to distributed systems discussion
✅ **Clean Go idioms**: Mutex, defer, pointer receivers, zero values

This shows you can:
- Implement algorithms correctly
- Write clean, maintainable code
- Think about edge cases
- Design good APIs
- Understand concurrency
- Discuss system design trade-offs
