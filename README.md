# Go Rate Limiting

A flexible and efficient rate limiting library for Go with support for multiple algorithms and Redis-based distributed rate limiting.

## Features

- **Multiple Algorithms**
  - Token Bucket (recommended for most APIs - allows bursts)
  - Fixed Window Counter (simple, fast)
  - Sliding Window Log (most accurate, memory intensive)
  - Leaky Bucket (smooths traffic)

- **Flexible Storage Backends**
  - In-memory (single server)
  - Redis (distributed systems)

- **HTTP Middleware**
  - Easy integration with standard `net/http`
  - Flexible key extraction (API key, IP, User ID, etc.)
  - Standard rate limit headers

## Installation

```bash
go get github.com/kaldun/go-rate-limiting
```

## Quick Start

### Basic Usage with In-Memory Storage

```go
package main

import (
    "context"
    "time"

    "github.com/kaldun/go-rate-limiting/ratelimiter"
    "github.com/kaldun/go-rate-limiting/storage"
)

func main() {
    // Create storage
    store := storage.NewMemoryStorage()
    defer store.Close()

    // Create rate limiter: 10 requests per second
    config := ratelimiter.Config{
        Rate:   10,
        Window: time.Second,
    }
    limiter := ratelimiter.NewTokenBucket(store, config)

    // Check if request is allowed
    ctx := context.Background()
    allowed, err := limiter.Allow(ctx, "user:123")
    if err != nil {
        // handle error
    }

    if allowed {
        // process request
    } else {
        // reject with 429 Too Many Requests
    }
}
```

### With Redis (Distributed Systems)

```go
import (
    "github.com/redis/go-redis/v9"
)

// Create Redis storage
store, err := storage.NewRedisStorage(&redis.Options{
    Addr: "localhost:6379",
})
if err != nil {
    // handle error
}
defer store.Close()

// Use same as above
limiter := ratelimiter.NewFixedWindow(store, config)
```

### HTTP Middleware

```go
import (
    "net/http"
    "github.com/kaldun/go-rate-limiting/middleware"
)

// Create middleware
rateLimitMW := middleware.NewRateLimitMiddleware(
    limiter,
    middleware.ExtractAPIKey("X-API-Key"),
)

// Wrap your handler
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Success!")
})

http.Handle("/api/endpoint", rateLimitMW.Middleware(handler))
```

## Algorithms Comparison

| Algorithm | Memory | Performance | Accuracy | Bursts | Use Case |
|-----------|--------|-------------|----------|--------|----------|
| Token Bucket | O(1) | O(1) | Good | Yes | Most APIs |
| Fixed Window | O(1) | O(1) | Good | Boundary issue | Simple limits |
| Sliding Window | O(N) | O(N) | Excellent | No | Strict accuracy |
| Leaky Bucket | O(1) | O(1) | Good | Smoothed | Traffic shaping |

**Recommendation**: Token Bucket for most use cases - allows reasonable bursts while maintaining rate limits.

## Project Structure

```
go-rate-limiting/
├── ratelimiter/        # Rate limiting algorithms
│   ├── ratelimiter.go  # Interfaces and types
│   ├── tokenbucket.go
│   ├── fixedwindow.go
│   ├── slidingwindow.go
│   └── leakybucket.go
├── storage/            # Storage backends
│   ├── storage.go      # Interface
│   ├── redis.go
│   └── memory.go
├── middleware/         # HTTP middleware
│   └── http.go
└── examples/           # Usage examples
    └── main.go
```

## Development

```bash
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestTokenBucket ./ratelimiter

# Run benchmarks
go test -bench=. ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## TODO

The core algorithm implementations are left as exercises. See TODO comments in:
- `ratelimiter/*.go` - Implement Allow/AllowN/AllowWithInfo methods
- `storage/redis.go` - Implement Redis operations
- `storage/memory.go` - Implement in-memory operations
- `middleware/http.go` - Implement middleware logic
- `*_test.go` - Implement test cases

## License

MIT
