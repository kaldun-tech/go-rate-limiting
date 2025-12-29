package ratelimiter

import "time"

// Result contains information about a rate limit check
type Result struct {
	// Allowed indicates if the request is allowed
	Allowed bool

	// Remaining is the number of requests remaining in the current window
	Remaining int

	// RetryAfter is the duration to wait before retrying (if not allowed)
	RetryAfter time.Duration

	// ResetAt is when the rate limit will reset to full capacity
	ResetAt time.Time
}
