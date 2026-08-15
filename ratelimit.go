package goensemblrest

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimitInfo holds the latest rate limiting metadata parsed from Ensembl response headers.
type RateLimitInfo struct {
	// Reset is the timestamp in seconds when the rate limit window resets.
	Reset *int `json:"rate_reset,omitempty"`

	// Limit is the maximum number of requests allowed within the period.
	Limit *int `json:"rate_limit,omitempty"`

	// Remaining is the number of remaining requests allowed in the current window.
	Remaining *int `json:"rate_remaining,omitempty"`

	// Period is the duration in seconds of the rate limit window.
	Period *int `json:"rate_period,omitempty"`

	// RetryAfter is the number of seconds to wait before retrying when rate limited.
	RetryAfter *float64 `json:"retry_after,omitempty"`
}

// rateLimiter implements a thread-safe sliding window rate limiter.
type rateLimiter struct {
	mu         sync.Mutex
	timestamps []time.Time
	reqsPerSec int
	window     time.Duration
	info       RateLimitInfo
}

// newRateLimiter creates a new sliding-window rate limiter.
func newRateLimiter(reqsPerSec int, window time.Duration) *rateLimiter {
	if reqsPerSec <= 0 {
		reqsPerSec = 15
	}
	if window <= 0 {
		window = time.Second
	}
	return &rateLimiter{
		reqsPerSec: reqsPerSec,
		window:     window,
		timestamps: make([]time.Time, 0, reqsPerSec),
	}
}

// Wait blocks until the caller is permitted to make a request under the rate limit,
// or returns early if the context is cancelled.
func (rl *rateLimiter) Wait(ctx context.Context) error {
	for {
		rl.mu.Lock()
		now := time.Now()

		// Prune timestamps older than the sliding window
		cutoff := now.Add(-rl.window)
		validIdx := 0
		for i, t := range rl.timestamps {
			if t.After(cutoff) {
				validIdx = i
				break
			}
			if i == len(rl.timestamps)-1 {
				validIdx = len(rl.timestamps)
			}
		}
		if validIdx > 0 {
			rl.timestamps = rl.timestamps[validIdx:]
		}

		// If under capacity, record timestamp and proceed
		if len(rl.timestamps) < rl.reqsPerSec {
			rl.timestamps = append(rl.timestamps, now)
			rl.mu.Unlock()
			return nil
		}

		// Compute required sleep duration until the oldest entry expires
		oldest := rl.timestamps[0]
		toSleep := rl.window - now.Sub(oldest)
		rl.mu.Unlock()

		if toSleep <= 0 {
			continue
		}

		// Sleep or abort on context cancellation
		timer := time.NewTimer(toSleep)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			// Loop around to re-check window
		}
	}
}

// UpdateFromHeaders extracts rate limit information from HTTP response headers.
func (rl *rateLimiter) UpdateFromHeaders(headers http.Header) RateLimitInfo {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if resetStr := headers.Get("X-RateLimit-Reset"); resetStr != "" {
		if v, err := strconv.Atoi(resetStr); err == nil {
			rl.info.Reset = &v
		}
	}
	if limitStr := headers.Get("X-RateLimit-Limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			rl.info.Limit = &v
		}
	}
	if remStr := headers.Get("X-RateLimit-Remaining"); remStr != "" {
		if v, err := strconv.Atoi(remStr); err == nil {
			rl.info.Remaining = &v
		}
	}
	if periodStr := headers.Get("X-RateLimit-Period"); periodStr != "" {
		if v, err := strconv.Atoi(periodStr); err == nil {
			rl.info.Period = &v
		}
	}
	if retryStr := headers.Get("Retry-After"); retryStr != "" {
		if v, err := strconv.ParseFloat(retryStr, 64); err == nil {
			rl.info.RetryAfter = &v
		}
	}

	return rl.info
}

// GetInfo returns a copy of the current rate limit metadata.
func (rl *rateLimiter) GetInfo() RateLimitInfo {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.info
}
