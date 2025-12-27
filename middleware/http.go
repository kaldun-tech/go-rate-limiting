package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kaldun-tech/go-rate-limiting/ratelimiter"
)

// KeyExtractor is a function that extracts the rate limiting key from a request
// Common implementations:
// - Extract from API key header
// - Extract from user ID (requires authentication)
// - Extract from IP address
// - Combine multiple factors
type KeyExtractor func(r *http.Request) string

// RateLimitMiddleware creates HTTP middleware for rate limiting
type RateLimitMiddleware struct {
	limiter         ratelimiter.RateLimiter
	keyExtractor    KeyExtractor
	onLimitExceeded func(w http.ResponseWriter, r *http.Request)
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(limiter ratelimiter.RateLimiter, keyExtractor KeyExtractor) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter:         limiter,
		keyExtractor:    keyExtractor,
		onLimitExceeded: defaultOnLimitExceeded,
	}
}

// WithOnLimitExceeded sets a custom handler for when rate limit is exceeded
func (m *RateLimitMiddleware) WithOnLimitExceeded(handler func(w http.ResponseWriter, r *http.Request)) *RateLimitMiddleware {
	m.onLimitExceeded = handler
	return m
}

// Middleware returns the HTTP middleware handler
func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement middleware logic
		// Steps:
		// 1. Extract key from request using keyExtractor
		// 2. Check rate limit using limiter.Allow()
		// 3. If allowed, set rate limit headers and call next handler
		// 4. If not allowed, call onLimitExceeded handler
		//
		// Standard rate limit headers to set:
		// - X-RateLimit-Limit: maximum requests allowed
		// - X-RateLimit-Remaining: requests remaining
		// - X-RateLimit-Reset: when the limit resets (Unix timestamp)
		// - Retry-After: seconds to wait before retrying (if limit exceeded)

		key := m.keyExtractor(r)
		_ = key // TODO: Remove this line after implementing

		// For now, just pass through
		next.ServeHTTP(w, r)
	})
}

// defaultOnLimitExceeded is the default handler for rate limit exceeded
func defaultOnLimitExceeded(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, `{"error":"rate limit exceeded","message":"Too many requests. Please try again later."}`)
}

// Common key extractors

// ExtractAPIKey extracts the API key from a header
func ExtractAPIKey(headerName string) KeyExtractor {
	return func(r *http.Request) string {
		apiKey := r.Header.Get(headerName)
		if apiKey == "" {
			return "anonymous"
		}
		return fmt.Sprintf("api_key:%s", apiKey)
	}
}

// ExtractIPAddress extracts the client IP address
func ExtractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return fmt.Sprintf("ip:%s", ip)
	}

	// Fall back to RemoteAddr
	return fmt.Sprintf("ip:%s", r.RemoteAddr)
}

// ExtractUserID extracts user ID from request context
// Assumes authentication middleware has set user ID in context
func ExtractUserID(contextKey string) KeyExtractor {
	return func(r *http.Request) string {
		userID := r.Context().Value(contextKey)
		if userID == nil {
			return "anonymous"
		}
		return fmt.Sprintf("user:%v", userID)
	}
}

// CombineExtractors combines multiple key extractors
// Useful for compound keys like "api_key:abc123:endpoint:/create"
func CombineExtractors(extractors ...KeyExtractor) KeyExtractor {
	return func(r *http.Request) string {
		var parts []string
		for _, extractor := range extractors {
			parts = append(parts, extractor(r))
		}
		result := ""
		for i, part := range parts {
			if i > 0 {
				result += ":"
			}
			result += part
		}
		return result
	}
}

// ExtractEndpoint extracts the request endpoint (path)
func ExtractEndpoint(r *http.Request) string {
	return fmt.Sprintf("endpoint:%s", r.URL.Path)
}

// Helper functions for setting rate limit headers

// SetRateLimitHeaders sets standard rate limit headers on the response
func SetRateLimitHeaders(w http.ResponseWriter, limit, remaining int, resetAt int64) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt, 10))
}

// SetRetryAfterHeader sets the Retry-After header
func SetRetryAfterHeader(w http.ResponseWriter, seconds int) {
	w.Header().Set("Retry-After", strconv.Itoa(seconds))
}
