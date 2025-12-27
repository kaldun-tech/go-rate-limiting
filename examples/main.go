package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kaldun-tech/go-rate-limiting/middleware"
	"github.com/kaldun-tech/go-rate-limiting/ratelimiter"
	"github.com/kaldun-tech/go-rate-limiting/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Example 1: In-memory rate limiter (single server)
	fmt.Println("Example 1: In-memory rate limiter")
	exampleInMemory()

	// Example 2: Redis rate limiter (distributed)
	fmt.Println("\nExample 2: Redis rate limiter")
	exampleRedis()

	// Example 3: HTTP middleware
	fmt.Println("\nExample 3: HTTP middleware")
	exampleHTTPMiddleware()
}

func exampleInMemory() {
	// Create in-memory storage
	store := storage.NewMemoryStorage()
	defer store.Close()

	// Create token bucket rate limiter: 10 requests per second with burst of 20
	config := ratelimiter.Config{
		Rate:      10,
		Window:    time.Second,
		BurstSize: 20,
	}
	limiter := ratelimiter.NewTokenBucket(store, config)

	// Test rate limiting
	ctx := context.Background()
	key := "user:123"

	// Simulate requests
	for i := 0; i < 15; i++ {
		allowed, err := limiter.Allow(ctx, key)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		if allowed {
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: RATE LIMITED\n", i+1)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func exampleRedis() {
	// Create Redis storage
	store, err := storage.NewRedisStorage(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
	if err != nil {
		log.Printf("Failed to connect to Redis (is it running?): %v", err)
		log.Println("Skipping Redis example...")
		return
	}
	defer store.Close()

	// Create fixed window rate limiter: 100 requests per minute
	config := ratelimiter.Config{
		Rate:   100,
		Window: time.Minute,
	}
	limiter := ratelimiter.NewFixedWindow(store, config)

	// Test rate limiting
	ctx := context.Background()
	key := "user:456"

	// Simulate requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		if allowed {
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: RATE LIMITED\n", i+1)
		}
	}
}

func exampleHTTPMiddleware() {
	// Create in-memory storage
	store := storage.NewMemoryStorage()

	// Create rate limiter: 5 requests per 10 seconds
	config := ratelimiter.Config{
		Rate:   5,
		Window: 10 * time.Second,
	}
	limiter := ratelimiter.NewTokenBucket(store, config)

	// Create middleware with IP-based rate limiting
	rateLimitMW := middleware.NewRateLimitMiddleware(
		limiter,
		middleware.ExtractIPAddress,
	)

	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello! This endpoint is rate limited.\n")
	})

	// Wrap handler with rate limiting middleware
	http.Handle("/api/hello", rateLimitMW.Middleware(handler))

	// Start server
	fmt.Println("Starting server on :8080")
	fmt.Println("Try: curl http://localhost:8080/api/hello")
	fmt.Println("Make multiple requests to see rate limiting in action")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
