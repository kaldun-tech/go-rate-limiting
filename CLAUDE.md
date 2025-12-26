# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Rate limiting system implementation in Go with Redis support. This is a learning exercise focused on implementing distributed rate limiting algorithms and understanding system design trade-offs at scale.

## Development Commands

### Running Examples
- `go run examples/main.go` - Run example demonstrations of different algorithms
- Redis must be running on `localhost:6379` for Redis examples

### Testing
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output
- `go test -run TestTokenBucket ./ratelimiter` - Run specific test
- `go test -bench=. ./...` - Run benchmarks
- `go test -bench=BenchmarkTokenBucket ./ratelimiter` - Run specific benchmark

### Code Quality
- `go fmt ./...` - Format code
- `go vet ./...` - Run static analysis
- `go mod tidy` - Clean up dependencies

## Architecture

### Package Structure

```
ratelimiter/     - Core rate limiting algorithms (Token Bucket, Fixed Window, Sliding Window, Leaky Bucket)
storage/         - Storage abstraction layer (Redis, in-memory)
middleware/      - HTTP middleware for easy integration
examples/        - Usage examples and demonstrations
```

### Rate Limiting Algorithms

The project implements four algorithms, each with different trade-offs:

#### Token Bucket (Recommended)
- **Memory**: O(1) - stores only current tokens + last refill time
- **Performance**: O(1) operations
- **Allows bursts**: Up to bucket capacity
- **Best for**: Most APIs, allows reasonable traffic spikes
- **Implementation notes**:
  - Tokens refill at constant rate (Rate/Window)
  - Bucket capacity = BurstSize (defaults to Rate)
  - Each request consumes 1 token (or N for AllowN)

#### Fixed Window Counter
- **Memory**: O(1) - single counter per window
- **Performance**: O(1) - simple INCR operation
- **Boundary problem**: Can allow 2x rate at window boundaries
- **Best for**: Simple limits where small inaccuracies acceptable
- **Implementation notes**:
  - Window key includes timestamp: `key:window:1234567890`
  - Use Redis INCR + EXPIRE for atomic operations
  - Counter resets at fixed intervals

#### Sliding Window Log
- **Memory**: O(N) - stores timestamp for every request
- **Performance**: O(N) - must scan/remove old timestamps
- **Most accurate**: True sliding window, no boundary issues
- **Best for**: When accuracy is critical
- **Implementation notes**:
  - Use Redis Sorted Set (ZSET) with timestamp as score
  - ZREMRANGEBYSCORE to remove old entries
  - ZCARD to count current requests
  - 50x more memory than Token Bucket

#### Leaky Bucket
- **Memory**: O(1) - queue size + last leak time
- **Performance**: O(1) operations
- **Smooths bursts**: Enforces constant rate
- **Best for**: Protecting downstream services, traffic shaping
- **Implementation notes**:
  - Requests drain from queue at constant rate
  - New requests rejected if queue full
  - Calculate leaked requests based on elapsed time

### Storage Layer Design

**Storage Interface**: Abstraction that supports both Redis and in-memory backends

Key operations needed:
- `Get/Set` - Basic key-value operations
- `Increment/IncrementBy` - Atomic counters (for Fixed Window)
- `Delete/Expire` - Cleanup and TTL
- `ZAdd/ZRemRangeByScore/ZCount` - Sorted sets (for Sliding Window)
- `Eval` - Lua scripts (for atomic multi-step operations)

**Redis Implementation**:
- Used for distributed systems (multiple API servers)
- Atomic operations prevent race conditions
- Lua scripts ensure atomicity across multiple operations
- Sharding strategy: hash(user_id) % num_shards

**Memory Implementation**:
- Good for single server or testing
- Uses sync.RWMutex for thread safety
- Background goroutine cleans expired keys
- NOT suitable for distributed deployments

### HTTP Middleware Design

**Key Extraction Strategies**:
- `ExtractAPIKey` - Best for authenticated APIs (unique per client)
- `ExtractIPAddress` - Fallback for unauthenticated endpoints (can be spoofed/shared)
- `ExtractUserID` - For per-user limits (requires auth middleware)
- `CombineExtractors` - Compound keys like "api_key:abc:endpoint:/create"

**Standard Headers Set**:
- `X-RateLimit-Limit` - Maximum requests allowed
- `X-RateLimit-Remaining` - Requests remaining in window
- `X-RateLimit-Reset` - When limit resets (Unix timestamp)
- `Retry-After` - Seconds to wait (when limit exceeded)
- HTTP 429 status code when rate limited

### Distributed System Considerations

**Architecture at Scale**:
```
[Load Balancer]
    ↓
[API Gateway Cluster] (stateless, auto-scaling)
    ↓
[Redis Cluster] (sharded by user ID)
    ↓
[Backend Services]
```

**Race Condition Prevention**:
- Multiple API gateway instances may check same user's limit simultaneously
- Use Redis atomic operations (INCR, Lua scripts)
- Lua scripts execute atomically - all operations in one transaction

**Failure Modes**:
- Redis unavailable: Fail open (allow) vs fail closed (reject)
- Network partition: Local fallback with conservative limits
- Redis shard down: Failover via Redis Sentinel/Cluster

**Performance Optimizations**:
- Local caching: Check memory cache before Redis
- Batch operations: Group multiple checks (for async workloads)
- Approximate counting: Count-Min Sketch for less critical limits

### Algorithm Selection Guide

**Choose Token Bucket when**:
- Building a typical REST API
- Want to allow reasonable bursts
- Need good performance (O(1))
- Memory efficiency matters

**Choose Fixed Window when**:
- Simplicity is priority
- Small boundary inaccuracies acceptable
- Need absolute best performance
- Easy debugging/monitoring needed

**Choose Sliding Window when**:
- Accuracy is critical (financial, security)
- Can afford memory cost
- No burst allowance acceptable
- Willing to accept O(N) performance

**Choose Leaky Bucket when**:
- Protecting downstream services
- Need to smooth traffic
- Constant rate more important than bursts
- Traffic shaping is the goal

## Implementation Notes

### TODO Items

Core implementations left as exercises (marked with TODO comments):

1. **Algorithm implementations** (`ratelimiter/*.go`):
   - Implement `Allow()`, `AllowN()`, `AllowWithInfo()` methods
   - Calculate token refills, window boundaries, timestamps
   - Return detailed Result objects with remaining/retry info

2. **Storage backends** (`storage/*.go`):
   - Redis: Implement all operations using go-redis client
   - Memory: Implement with proper mutex locking and expiry
   - Handle errors and edge cases

3. **HTTP Middleware** (`middleware/http.go`):
   - Extract keys from requests
   - Check rate limits
   - Set appropriate headers
   - Handle limit exceeded responses

4. **Tests** (`*_test.go`):
   - Unit tests for each algorithm
   - Storage backend tests
   - Middleware integration tests
   - Benchmarks for performance comparison

### Key Questions Answered

**Where does rate limiting happen?**
- API Gateway layer (recommended) - protects backend, centralized control
- Can also implement at application layer for per-resource limits

**What's the limiting key?**
- API key for authenticated requests (best, unique per client)
- IP address as fallback (can be shared/spoofed)
- Consider combining factors for stricter control

**How to handle distributed systems?**
- Redis provides shared state across gateway instances
- Use atomic operations (INCR, Lua scripts) to prevent race conditions
- Shard Redis by user ID for horizontal scaling

**What happens when limit exceeded?**
- Return HTTP 429 Too Many Requests
- Set Retry-After header
- For background jobs: queuing makes sense
- For real-time APIs: reject immediately

### Common Pitfalls

- **Clock skew**: Use Redis TIME command for consistent timestamps across servers
- **Race conditions**: Always use atomic operations, never read-modify-write
- **Memory leaks**: Ensure old data is cleaned up (EXPIRE, sorted set trimming)
- **Testing**: Don't rely on time.Sleep for timing-sensitive tests - use mocking
