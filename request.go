package goensemblrest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	paramRegex = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

	knownTransientErrors = []string{
		"something bad has happened",
		"Something went wrong while fetching from LDFeatureContainerAdaptor",
		"timeout",
	}
)

// resolvePath interpolates path parameters into a URL template, preserving colons while escaping other characters.
func resolvePath(template string, params map[string]string) (string, error) {
	var missingParams []string

	resolved := paramRegex.ReplaceAllStringFunc(template, func(match string) string {
		key := strings.Trim(match, "{}")
		val, ok := params[key]
		if !ok || val == "" {
			missingParams = append(missingParams, key)
			return match
		}
		// Escape path component while preserving colon ':'
		escaped := url.PathEscape(val)
		escaped = strings.ReplaceAll(escaped, "%3A", ":")
		return escaped
	})

	if len(missingParams) > 0 {
		return "", fmt.Errorf("mandatory param %q not specified", missingParams[0])
	}

	return resolved, nil
}

// buildURL joins the base URL with the resolved path and appends query parameters.
func (c *Client) buildURL(resolvedPath string, query url.Values) (*url.URL, error) {
	c.mu.RLock()
	baseURL := c.baseURL
	c.mu.RUnlock()

	basePath := strings.TrimRight(baseURL.Path, "/")
	cleanPath := strings.TrimLeft(resolvedPath, "/")
	fullPath := basePath + "/" + cleanPath

	u := *baseURL
	u.Path = fullPath
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	return &u, nil
}

// isTransientError checks if a status code or body message is known to be transient and retryable.
func isTransientError(statusCode int, msg string) bool {
	if statusCode == http.StatusInternalServerError || statusCode == http.StatusRequestTimeout {
		return true
	}
	if statusCode == http.StatusBadRequest {
		lowerMsg := strings.ToLower(msg)
		for _, known := range knownTransientErrors {
			if strings.Contains(lowerMsg, strings.ToLower(known)) {
				return true
			}
		}
	}
	return false
}

// executeRequest sends the HTTP request with rate-limiting, exponential backoff retries, and unmarshaling.
func (c *Client) executeRequest(
	ctx context.Context,
	method string,
	pathTemplate string,
	pathParams map[string]string,
	bodyData any,
	target any,
	reqCfg RequestConfig,
) error {
	// 1. Resolve path variables
	resolvedPath, err := resolvePath(pathTemplate, pathParams)
	if err != nil {
		return err
	}

	// 2. Build full URL with query parameters
	reqURL, err := c.buildURL(resolvedPath, reqCfg.QueryParams)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	// 3. Serialize request body if provided
	var bodyBytes []byte
	if bodyData != nil {
		switch v := bodyData.(type) {
		case []byte:
			bodyBytes = v
		case string:
			bodyBytes = []byte(v)
		default:
			var err error
			bodyBytes, err = json.Marshal(bodyData)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
		}
	}

	c.mu.RLock()
	maxAttempts := c.maxAttempts
	wallTime := c.wallTime
	userAgent := c.userAgent
	extraHeaders := c.extraHeaders.Clone()
	c.mu.RUnlock()

	var lastErr error
	var lastAPIError *APIError

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Enforce sliding window rate limit
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}

		// Prepare HTTP request
		var bodyReader io.Reader
		if len(bodyBytes) > 0 {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		httpReq, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %w", err)
		}

		// Apply headers
		for k, vs := range extraHeaders {
			for _, v := range vs {
				httpReq.Header.Add(k, v)
			}
		}
		for k, v := range reqCfg.Headers {
			httpReq.Header.Set(k, v)
		}

		httpReq.Header.Set("User-Agent", userAgent)
		if reqCfg.ContentType != "" {
			httpReq.Header.Set("Content-Type", reqCfg.ContentType)
			httpReq.Header.Set("Accept", reqCfg.ContentType)
		} else {
			httpReq.Header.Set("Content-Type", DefaultContentType)
			httpReq.Header.Set("Accept", DefaultContentType)
		}

		// Execute HTTP call
		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			// Network or context failure
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			lastErr = err
			lastAPIError = &APIError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    fmt.Sprintf("network connection error: %v", err),
			}
		} else {
			defer resp.Body.Close()

			// Parse and update rate limit state from response headers
			rateInfo := c.limiter.UpdateFromHeaders(resp.Header)

			respBody, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				lastErr = readErr
				continue
			}

			// Handle successful response (2xx)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return unmarshalResponse(respBody, reqCfg.ContentType, target)
			}

			// Parse error message
			errorMsg := parseErrorMessage(respBody, http.StatusText(resp.StatusCode))

			lastAPIError = &APIError{
				StatusCode:    resp.StatusCode,
				Message:       errorMsg,
				RateReset:     rateInfo.Reset,
				RateLimit:     rateInfo.Limit,
				RateRemaining: rateInfo.Remaining,
				RetryAfter:    rateInfo.RetryAfter,
				RawBody:       respBody,
			}
			lastErr = lastAPIError

			// Check if error is not retryable
			if !isTransientError(resp.StatusCode, errorMsg) && resp.StatusCode != http.StatusTooManyRequests {
				return lastAPIError
			}
		}

		// Compute backoff duration before retry
		if attempt < maxAttempts {
			var sleepDuration time.Duration
			if lastAPIError != nil && lastAPIError.RetryAfter != nil && *lastAPIError.RetryAfter > 0 {
				sleepDuration = time.Duration(*lastAPIError.RetryAfter * float64(time.Second))
			} else {
				sleepDuration = time.Duration(attempt) * (wallTime * 2)
				if sleepDuration < 10*time.Millisecond {
					sleepDuration = 10 * time.Millisecond
				}
			}

			timer := time.NewTimer(sleepDuration)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}

	if lastAPIError != nil {
		return fmt.Errorf("%w: %s (attempts: %d)", ErrMaxRetriesReached, lastAPIError.Error(), maxAttempts)
	}

	return fmt.Errorf("%w: %v (attempts: %d)", ErrMaxRetriesReached, lastErr, maxAttempts)
}

// unmarshalResponse decodes the response body into the target value.
func unmarshalResponse(body []byte, contentType string, target any) error {
	if target == nil {
		return nil
	}

	switch v := target.(type) {
	case *string:
		*v = string(body)
		return nil
	case *[]byte:
		*v = body
		return nil
	default:
		// Attempt JSON unmarshal
		if err := json.Unmarshal(body, target); err != nil {
			// If JSON decode fails and target is an interface, return string
			if ifacePtr, ok := target.(*any); ok {
				*ifacePtr = string(body)
				return nil
			}
			return fmt.Errorf("failed to decode response as JSON: %w", err)
		}
		return nil
	}
}
