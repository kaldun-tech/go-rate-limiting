# Token Bucket Rate Limiter Interview Questions

Practice guide for demonstrating Token Bucket algorithm understanding in technical interviews.

## Core Implementation Questions

### 1. Design Token Bucket Rate Limiter
**Question:** Design a rate limiter using the token bucket algorithm.

**Requirements:**
- Limit requests to `rate` per `window` (e.g., 100 requests/minute)
- Allow bursts up to `burstSize`
- Per-key rate limiting (different users/IPs)
- Return true/false for allow/deny

**Key Points to Cover:**
- Each key has a bucket with tokens
- Tokens refill continuously over time
- Requests consume tokens
- Deny if not enough tokens
- Time: O(1), Space: O(number of keys)

**Data structure:**
```go
type TokenBucket struct {
    mu        sync.Mutex
    buckets   map[string]*bucket
    rate      int             // tokens per window
    burstSize int             // max tokens
    window    time.Duration
}

type bucket struct {
    tokens     float64
    lastRefill time.Time
}
```

**Follow-up:** Why use float64 for tokens but int for requests?

---

### 2. Token Refill Calculation
**Question:** How do you calculate tokens to add since last access?

**Key Points to Cover:**
```go
elapsed := now.Sub(b.lastRefill)
refillRate := float64(rate) / window.Seconds()
tokensToAdd := refillRate * elapsed.Seconds()
newTokens := math.Min(currentTokens + tokensToAdd, float64(burstSize))
```

**Important details:**
- Refill rate is fractional: 100 tokens/minute = 1.667 tokens/second
- Must use float64 to avoid losing precision
- Cap at burstSize - can't exceed capacity
- Update lastRefill to current time

**Common mistake:** Using integer division loses precision!
```go
// WRONG
refillRate := rate / int(window.Seconds())  // 10/60 = 0!

// RIGHT
refillRate := float64(rate) / window.Seconds()  // 10/60 = 0.167
```

**Follow-up:** What happens if the system clock moves backward?

---

### 3. Allow Operation
**Question:** Implement the core Allow(key) operation.

**Key Points to Cover:**
```go
func (tb *TokenBucket) Allow(key string) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    b, exists := tb.buckets[key]
    if !exists {
        // First request - create bucket
        if 1 > tb.burstSize {
            return false  // Requesting more than capacity
        }
        tb.buckets[key] = &bucket{
            tokens:     float64(tb.burstSize - 1),
            lastRefill: time.Now(),
        }
        return true
    }

    // Refill tokens
    now := time.Now()
    elapsed := now.Sub(b.lastRefill)
    refillRate := float64(tb.rate) / tb.window.Seconds()
    tokensToAdd := refillRate * elapsed.Seconds()
    b.tokens = math.Min(b.tokens+tokensToAdd, float64(tb.burstSize))
    b.lastRefill = now

    // Check and consume
    if b.tokens >= 1 {
        b.tokens -= 1
        return true
    }

    return false
}
```

**Edge cases:**
- First request (bucket doesn't exist)
- Bucket is full (cap at burstSize)
- Not enough tokens (deny but still update lastRefill!)
- Zero elapsed time (same timestamp)

**Follow-up:** Why update lastRefill even when denying?

---

### 4. AllowN Operation
**Question:** Extend to allow requests with different costs.

**Use case:** Read = 1 token, Write = 5 tokens, AI inference = 100 tokens

**Key differences from Allow:**
```go
func (tb *TokenBucket) AllowN(key string, n int) bool {
    // ... same refill logic ...

    if b.tokens >= float64(n) {
        b.tokens -= float64(n)
        return true
    }
    return false
}
```

**Important edge case:**
```go
if n > tb.burstSize {
    return false  // Can NEVER have enough tokens
}
```

**Follow-up:** Should we allow n=0? (Yes, useful for checking state without consuming)

---

### 5. Thread Safety
**Question:** How do you make this thread-safe?

**Answer:**
```go
type TokenBucket struct {
    mu sync.Mutex  // Protects buckets map and individual buckets
    // ...
}

func (tb *TokenBucket) Allow(key string) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    // ... entire operation under lock ...
}
```

**Why mutex not RWMutex?**
- Every operation modifies state (updates lastRefill)
- No pure read operations
- RWMutex provides no benefit

**Alternative:** Per-key mutexes for better concurrency
```go
type bucket struct {
    mu         sync.Mutex
    tokens     float64
    lastRefill time.Time
}
```

**Trade-off:** More memory, more complexity, better parallelism

**Follow-up:** How many goroutines can execute Allow() simultaneously with global mutex?

---

### 6. AllowWithInfo Operation
**Question:** Return detailed information for HTTP headers.

**Use case:** HTTP 429 responses need Retry-After header

```go
type Result struct {
    Allowed    bool
    Remaining  int
    RetryAfter time.Duration
    ResetAt    time.Time
}

func (tb *TokenBucket) AllowWithInfo(key string, n int) Result {
    // ... same logic as AllowN ...

    if allowed {
        return Result{
            Allowed:   true,
            Remaining: int(b.tokens),
            ResetAt:   now.Add(tb.window),
        }
    } else {
        tokensNeeded := float64(n) - b.tokens
        retryAfter := time.Duration(tokensNeeded/refillRate) * time.Second
        return Result{
            Allowed:    false,
            Remaining:  int(b.tokens),
            RetryAfter: retryAfter,
            ResetAt:    now.Add(tb.window),
        }
    }
}
```

**HTTP integration:**
```go
w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))
if !result.Allowed {
    w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
    http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
}
```

**Follow-up:** What headers does GitHub/Stripe use?

---

## Conceptual Questions

### 7. Token Bucket vs Alternatives
**Question:** What are the different rate limiting algorithms?

**Algorithms:**

1. **Token Bucket** (What you implemented)
   - Tokens refill continuously
   - Allows bursts up to capacity
   - Smooth refill rate

2. **Leaky Bucket**
   - Requests enter bucket, processed at fixed rate
   - Smooths out bursts
   - No bursting allowed

3. **Fixed Window**
   - Count requests in time windows (e.g., 1:00-1:01)
   - Simple but allows 2x burst at boundaries

4. **Sliding Window Log**
   - Track timestamp of each request
   - Accurate but memory intensive

5. **Sliding Window Counter**
   - Hybrid: weighted count from previous window
   - Good balance of accuracy and efficiency

**Comparison:**

| Algorithm | Bursts | Accuracy | Memory | Complexity |
|-----------|--------|----------|--------|------------|
| Token Bucket | Yes | Good | O(keys) | Medium |
| Leaky Bucket | No | Good | O(keys) | Medium |
| Fixed Window | 2x at boundary | Poor | O(keys) | Simple |
| Sliding Log | No | Perfect | O(keys*rate) | Complex |
| Sliding Counter | Limited | Good | O(keys) | Medium |

**Follow-up:** When would you choose Token Bucket over Leaky Bucket?

---

### 8. Burst Handling
**Question:** Explain burst size and why it matters.

**Answer:**
- **Burst size** = max tokens in bucket = max requests at once
- **Rate** = refill rate = sustained throughput

**Example:** 100 requests/hour, burst of 10
- Can handle 10 requests instantly (burst)
- Then limited to ~1.67 requests/minute (rate)

**Use cases:**
- API allows batch operations (burst needed)
- Handle traffic spikes gracefully
- Balance responsiveness vs. abuse prevention

**Setting burst size:**
- `burstSize = rate` → No bursting (same as leaky bucket)
- `burstSize > rate` → Allow bursts
- `burstSize = 2 * rate` → Common choice

**Follow-up:** What if burst = 1?

---

### 9. Distributed Rate Limiting
**Question:** How would you implement this across multiple servers?

**Approaches:**

**1. Shared State (Redis)**
```lua
-- Lua script for atomic refill + consume
local key = KEYS[1]
local rate = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local now = redis.call('TIME')
-- ... refill logic ...
-- ... consume if enough tokens ...
return {allowed, remaining}
```

**Pros:** Accurate, centralized
**Cons:** Network latency, Redis is SPOF

**2. Local Caches + Eventual Consistency**
- Each server has local token bucket
- Divide rate by number of servers
- Accept over-limit during server changes

**Pros:** Fast, no network
**Cons:** Inaccurate, can exceed limit

**3. Sticky Sessions**
- Route same user to same server
- Each server has independent buckets

**Pros:** No coordination needed
**Cons:** Uneven load, doesn't survive server restart

**Follow-up:** How does AWS API Gateway handle this?

---

### 10. Memory Management
**Question:** Buckets grow unbounded. How to clean up inactive keys?

**Solutions:**

**1. TTL-based eviction**
```go
type bucket struct {
    tokens     float64
    lastRefill time.Time
    lastAccess time.Time  // Track access
}

// Background goroutine
func (tb *TokenBucket) cleanup() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        tb.mu.Lock()
        for key, b := range tb.buckets {
            if time.Since(b.lastAccess) > 24*time.Hour {
                delete(tb.buckets, key)
            }
        }
        tb.mu.Unlock()
    }
}
```

**2. LRU eviction**
- Limit total keys
- Evict least recently used

**3. External storage (Redis)**
- Redis handles TTL automatically
- No cleanup needed

**Follow-up:** What's the right TTL value?

---

## Coding Challenges

### 11. Implement Leaky Bucket
**Question:** Implement the leaky bucket algorithm.

**Key difference:**
```go
type LeakyBucket struct {
    queue    []time.Time  // Request timestamps
    capacity int
    rate     time.Duration
}

func (lb *LeakyBucket) Allow() bool {
    now := time.Now()

    // Leak: Remove old requests
    cutoff := now.Add(-lb.rate)
    lb.queue = filterAfter(lb.queue, cutoff)

    // Check capacity
    if len(lb.queue) < lb.capacity {
        lb.queue = append(lb.queue, now)
        return true
    }
    return false
}
```

**Difference:** Queue of requests vs. token count

---

### 12. Hierarchical Rate Limiting
**Question:** Implement per-user AND per-organization limits.

**Example:** 100 req/min per user, 1000 req/min per org

**Hint:**
```go
func (tb *TokenBucket) AllowMulti(userKey, orgKey string, n int) bool {
    // Must pass BOTH limits
    userOK := tb.Allow(userKey, n)
    if !userOK {
        return false
    }

    orgOK := tb.Allow(orgKey, n)
    if !orgOK {
        // Refund user tokens!
        tb.refund(userKey, n)
        return false
    }

    return true
}
```

**Follow-up:** How do you handle refunds?

---

### 13. Cost-Based Rate Limiting
**Question:** Different endpoints have different costs.

**Example:**
- GET /users: 1 token
- POST /users: 5 tokens
- POST /ai/generate: 100 tokens

**Implementation:**
```go
const (
    CostRead    = 1
    CostWrite   = 5
    CostAI      = 100
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    cost := getCost(r.URL.Path, r.Method)
    if !limiter.AllowN(getKey(r), cost) {
        http.Error(w, "Rate limit exceeded", 429)
        return
    }
    // ... handle request ...
}
```

---

### 14. Dynamic Rate Limits
**Question:** Rate limits change based on user tier (free/pro/enterprise).

**Implementation:**
```go
type RateLimiter struct {
    limiters map[string]*TokenBucket  // tier -> limiter
}

func (rl *RateLimiter) Allow(user User) bool {
    limiter := rl.limiters[user.Tier]
    return limiter.Allow(user.ID)
}
```

**Follow-up:** How do you handle tier upgrades?

---

### 15. Request Scheduling
**Question:** Instead of deny, return when request can proceed.

```go
func (tb *TokenBucket) Wait(key string) time.Duration {
    for {
        if tb.Allow(key) {
            return 0
        }
        waitTime := tb.calculateWaitTime(key)
        time.Sleep(waitTime)
    }
}
```

**Use case:** Background jobs that can wait

**Follow-up:** How would you implement with channels?

---

## Common Pitfalls

### Mistake 1: Integer Division Loses Precision
```go
// WRONG - 10 tokens/60 seconds = 0!
refillRate := tb.rate / int(tb.window.Seconds())

// RIGHT - 10 tokens/60 seconds = 0.167
refillRate := float64(tb.rate) / tb.window.Seconds()
```

---

### Mistake 2: Not Updating lastRefill on Deny
```go
// WRONG - lastRefill not updated when denying
if b.tokens >= float64(n) {
    b.tokens -= float64(n)
    b.lastRefill = now
    return true
}
return false

// RIGHT - always update lastRefill
b.lastRefill = now
if b.tokens >= float64(n) {
    b.tokens -= float64(n)
    return true
}
return false
```

**Why?** Next call needs accurate elapsed time for refill calculation.

---

### Mistake 3: Not Capping at Burst Size
```go
// WRONG - tokens can exceed capacity
b.tokens += tokensToAdd

// RIGHT - cap at burst size
b.tokens = math.Min(b.tokens + tokensToAdd, float64(tb.burstSize))
```

---

### Mistake 4: Forgetting First Request Case
```go
// WRONG - crashes on first request
b := tb.buckets[key]  // nil!
b.tokens += tokensToAdd  // panic!

// RIGHT - handle first request
b, exists := tb.buckets[key]
if !exists {
    tb.buckets[key] = &bucket{
        tokens:     float64(tb.burstSize - n),
        lastRefill: now,
    }
    return n <= tb.burstSize
}
```

---

## Design Discussion Questions

### 16. Real-World Examples
**Question:** How do major APIs implement rate limiting?

**GitHub:**
- 5,000 requests/hour for authenticated
- 60 requests/hour for unauthenticated
- Returns `X-RateLimit-*` headers

**Stripe:**
- Different limits per endpoint
- Token bucket with rolling window
- Separate limits for test vs live mode

**Twitter:**
- 15-minute windows
- Different limits per endpoint
- User vs app authentication

**AWS:**
- Per-service limits
- Burst capacity (token bucket)
- Request quota + burst quota

---

### 17. HTTP Status Codes
**Question:** What status code for rate limit exceeded?

**Answer:** `429 Too Many Requests`

**Headers to include:**
```
HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640000000
Retry-After: 30
```

**Follow-up:** Should you return different codes for different limits?

---

### 18. Testing Strategy
**Question:** How do you test rate limiters?

**Approaches:**
1. **Freeze time** - Control time.Now() with test clock
2. **Small windows** - 10 tokens/100ms instead of 100/sec
3. **Parallel requests** - Test thread safety
4. **Boundary conditions** - Exactly at limit
5. **Long-running** - Test refill over time

**Example:**
```go
func TestRefill(t *testing.T) {
    limiter := NewTokenBucket(10, time.Second, 10)

    // Drain bucket
    for i := 0; i < 10; i++ {
        limiter.Allow("user1")
    }

    // Should be denied
    if limiter.Allow("user1") {
        t.Error("Should be denied")
    }

    // Wait for refill
    time.Sleep(1 * time.Second)

    // Should be allowed
    if !limiter.Allow("user1") {
        t.Error("Should be allowed after refill")
    }
}
```

---

## Practice Drill

**Warm-up (5 minutes each):**
1. Explain token bucket algorithm in plain English
2. Calculate refill rate for 100 req/min
3. Draw bucket state over time with 5 req burst

**Medium (15 minutes each):**
4. Implement Allow() on whiteboard
5. Add AllowN() with cost parameter
6. Calculate RetryAfter duration

**Advanced (20+ minutes):**
7. Implement distributed version with Redis
8. Design multi-tier rate limiting system
9. Add request scheduling (wait instead of deny)

---

## Key Takeaways for Interviews

1. **Continuous refill** - Tokens add based on elapsed time, not periodic timer
2. **Float64 tokens** - Fractional tokens from refill calculation
3. **Always update lastRefill** - Even when denying request
4. **Cap at burst size** - Can't exceed capacity
5. **Thread safety** - Mutex protects shared state
6. **First request** - Initialize bucket with burst tokens
7. **Edge case** - Request > capacity always denied
8. **Real-world usage** - GitHub, Stripe, AWS, Twitter

**Interview Tips:**
- Start with simple Allow(), then extend to AllowN
- Draw timeline showing token refill
- Discuss trade-offs vs other algorithms
- Know distributed implementation with Redis
- Understand real-world applications

Good luck!
