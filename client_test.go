package goensemblrest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientInitialization(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		client, err := NewClient()
		if err != nil {
			t.Fatalf("unexpected error creating client: %v", err)
		}
		defer client.Close()

		if client.BaseURL() != DefaultBaseURL {
			t.Errorf("expected BaseURL %q, got %q", DefaultBaseURL, client.BaseURL())
		}
		if client.UserAgent() != DefaultUserAgent {
			t.Errorf("expected UserAgent %q, got %q", DefaultUserAgent, client.UserAgent())
		}
		if client.Timeout() != DefaultTimeout {
			t.Errorf("expected Timeout %v, got %v", DefaultTimeout, client.Timeout())
		}
		if client.MaxAttempts() != DefaultMaxAttempts {
			t.Errorf("expected MaxAttempts %d, got %d", DefaultMaxAttempts, client.MaxAttempts())
		}
	})

	t.Run("custom options", func(t *testing.T) {
		customHTTP := &http.Client{Timeout: 10 * time.Second}
		client, err := NewClient(
			WithBaseURL("https://custom.rest.org/"),
			WithTimeout(15*time.Second),
			WithMaxAttempts(3),
			WithRateLimit(10, 2*time.Second),
			WithUserAgent("CustomApp/1.0"),
			WithHeader("X-Custom", "TestHeader"),
			WithHTTPClient(customHTTP),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer client.Close()

		if client.BaseURL() != "https://custom.rest.org" {
			t.Errorf("expected trimmed BaseURL https://custom.rest.org, got %q", client.BaseURL())
		}
		if client.Timeout() != 15*time.Second {
			t.Errorf("expected Timeout 15s, got %v", client.Timeout())
		}
		if client.MaxAttempts() != 3 {
			t.Errorf("expected MaxAttempts 3, got %d", client.MaxAttempts())
		}
		if client.UserAgent() != "CustomApp/1.0" {
			t.Errorf("expected UserAgent CustomApp/1.0, got %q", client.UserAgent())
		}
	})

	t.Run("invalid options return error", func(t *testing.T) {
		if _, err := NewClient(WithBaseURL("::invalid-url")); err == nil {
			t.Error("expected error with invalid base url, got nil")
		}
		if _, err := NewClient(WithHTTPClient(nil)); err == nil {
			t.Error("expected error with nil http client, got nil")
		}
		if _, err := NewClient(WithTimeout(-1 * time.Second)); err == nil {
			t.Error("expected error with negative timeout, got nil")
		}
		if _, err := NewClient(WithRateLimit(0, time.Second)); err == nil {
			t.Error("expected error with 0 reqsPerSec, got nil")
		}
		if _, err := NewClient(WithMaxAttempts(0)); err == nil {
			t.Error("expected error with 0 max attempts, got nil")
		}
	})
}

func TestRateLimiter(t *testing.T) {
	t.Run("sliding window rate limit", func(t *testing.T) {
		rl := newRateLimiter(3, 200*time.Millisecond)
		ctx := context.Background()

		start := time.Now()
		// 3 immediate calls should complete quickly
		for i := 0; i < 3; i++ {
			if err := rl.Wait(ctx); err != nil {
				t.Fatalf("unexpected error on call %d: %v", i, err)
			}
		}
		initialDuration := time.Since(start)
		if initialDuration > 50*time.Millisecond {
			t.Errorf("initial 3 calls took too long: %v", initialDuration)
		}

		// 4th call must wait for sliding window to open
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on 4th call: %v", err)
		}
		totalDuration := time.Since(start)
		if totalDuration < 150*time.Millisecond {
			t.Errorf("expected rate limiter to throttle, but completed in %v", totalDuration)
		}
	})

	t.Run("context cancellation during rate limit wait", func(t *testing.T) {
		rl := newRateLimiter(1, 5*time.Second)
		ctx, cancel := context.WithCancel(context.Background())

		// First request takes capacity
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("first wait failed: %v", err)
		}

		// Cancel context quickly
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		// Second request should abort on cancellation
		err := rl.Wait(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})

	t.Run("parse rate limit headers", func(t *testing.T) {
		rl := newRateLimiter(15, time.Second)
		headers := http.Header{}
		headers.Set("X-RateLimit-Reset", "1600000000")
		headers.Set("X-RateLimit-Limit", "15")
		headers.Set("X-RateLimit-Remaining", "14")
		headers.Set("X-RateLimit-Period", "1")
		headers.Set("Retry-After", "5.5")

		info := rl.UpdateFromHeaders(headers)
		if info.Reset == nil || *info.Reset != 1600000000 {
			t.Errorf("expected Reset 1600000000, got %v", info.Reset)
		}
		if info.Limit == nil || *info.Limit != 15 {
			t.Errorf("expected Limit 15, got %v", info.Limit)
		}
		if info.Remaining == nil || *info.Remaining != 14 {
			t.Errorf("expected Remaining 14, got %v", info.Remaining)
		}
		if info.Period == nil || *info.Period != 1 {
			t.Errorf("expected Period 1, got %v", info.Period)
		}
		if info.RetryAfter == nil || *info.RetryAfter != 5.5 {
			t.Errorf("expected RetryAfter 5.5, got %v", info.RetryAfter)
		}
	})
}

func TestURLResolution(t *testing.T) {
	t.Run("resolve path variables and preserve colons", func(t *testing.T) {
		path, err := resolvePath("/sequence/region/{{species}}/{{region}}", map[string]string{
			"species": "homo sapiens",
			"region":  "X:1000..2000:1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "/sequence/region/homo%20sapiens/X:1000..2000:1"
		if path != expected {
			t.Errorf("expected %q, got %q", expected, path)
		}
	})

	t.Run("missing mandatory parameter error", func(t *testing.T) {
		_, err := resolvePath("/archive/id/{{id}}", map[string]string{})
		if err == nil || !strings.Contains(err.Error(), "mandatory param \"id\" not specified") {
			t.Errorf("expected missing param error, got %v", err)
		}
	})
}

func TestErrorsAndRetries(t *testing.T) {
	t.Run("404 Not Found error mapping", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID not found"})
		}))
		defer server.Close()

		client, err := NewClient(WithBaseURL(server.URL), WithMaxAttempts(1))
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		var res map[string]any
		err = client.GetArchiveByID(context.Background(), "NONEXISTENT", &res)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got %v", err)
		}

		var apiErr *APIError
		if errors.As(err, &apiErr) {
			if apiErr.StatusCode != http.StatusNotFound {
				t.Errorf("expected status code 404, got %d", apiErr.StatusCode)
			}
			if apiErr.Message != "ID not found" {
				t.Errorf("expected message 'ID not found', got %q", apiErr.Message)
			}
		} else {
			t.Errorf("expected errors.As(err, &apiErr) to succeed, got %v", err)
		}
	})

	t.Run("400 Bad Request error mapping", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid format"})
		}))
		defer server.Close()

		client, err := NewClient(WithBaseURL(server.URL), WithMaxAttempts(1))
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		var res map[string]any
		err = client.GetArchiveByID(context.Background(), "BAD", &res)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrBadRequest) {
			t.Errorf("expected errors.Is(err, ErrBadRequest) to be true, got %v", err)
		}
	})

	t.Run("retry transient 500 error then succeed", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := atomic.AddInt32(&attempts, 1)
			if count == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Temporary database glitch"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"id": "ENSG00000157764"})
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithMaxAttempts(3),
			WithRateLimit(100, 10*time.Millisecond),
		)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		var res map[string]string
		err = client.GetArchiveByID(context.Background(), "ENSG00000157764", &res)
		if err != nil {
			t.Fatalf("expected retry to succeed, got %v", err)
		}

		if atomic.LoadInt32(&attempts) != 2 {
			t.Errorf("expected 2 attempts, got %d", atomic.LoadInt32(&attempts))
		}
		if res["id"] != "ENSG00000157764" {
			t.Errorf("expected ID ENSG00000157764, got %q", res["id"])
		}
	})

	t.Run("exhaust retries on 500 error", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&attempts, 1)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Permanent crash"})
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithMaxAttempts(3),
			WithRateLimit(100, 10*time.Millisecond),
		)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		var res map[string]string
		err = client.GetArchiveByID(context.Background(), "ENSG00000157764", &res)
		if err == nil {
			t.Fatal("expected max retries error, got nil")
		}

		if !errors.Is(err, ErrMaxRetriesReached) {
			t.Errorf("expected ErrMaxRetriesReached, got %v", err)
		}
		if atomic.LoadInt32(&attempts) != 3 {
			t.Errorf("expected exactly 3 attempts, got %d", atomic.LoadInt32(&attempts))
		}
	})

	t.Run("retry on 429 rate limit with Retry-After", func(t *testing.T) {
		var attempts int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := atomic.AddInt32(&attempts, 1)
			if count == 1 {
				w.Header().Set("Retry-After", "0.05")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithMaxAttempts(3),
			WithRateLimit(100, 10*time.Millisecond),
		)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		var res map[string]string
		err = client.GetArchiveByID(context.Background(), "TEST", &res)
		if err != nil {
			t.Fatalf("expected rate limit retry to succeed, got %v", err)
		}

		if atomic.LoadInt32(&attempts) != 2 {
			t.Errorf("expected 2 attempts, got %d", atomic.LoadInt32(&attempts))
		}
	})
}

func TestDynamicCallAndEndpoints(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	endpoints := client.Endpoints()
	if len(endpoints) < 100 {
		t.Errorf("expected at least 100 endpoints, got %d", len(endpoints))
	}

	// Verify key endpoints exist in catalog
	for _, expectedName := range []string{
		"getArchiveById",
		"getArchiveByMultipleIds",
		"getLookupById",
		"getSequenceById",
		"getInfoPing",
		"getHomologyById",
		"getGA4GHBeacon",
	} {
		if _, ok := endpoints[expectedName]; !ok {
			t.Errorf("expected endpoint %q in table", expectedName)
		}
	}

	// Test unknown endpoint error
	err = client.Call(context.Background(), "unknownNonExistentApi", nil, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "unknown Ensembl REST API endpoint") {
		t.Errorf("expected unknown endpoint error, got %v", err)
	}
}
