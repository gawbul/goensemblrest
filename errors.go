package goensemblrest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for errors.Is matching.
var (
	// ErrBadRequest is returned when Ensembl responds with HTTP 400 Bad Request.
	ErrBadRequest = errors.New("ensembl: bad request (400)")

	// ErrNotFound is returned when Ensembl responds with HTTP 404 Not Found.
	ErrNotFound = errors.New("ensembl: not found (404)")

	// ErrTimeout is returned when Ensembl responds with HTTP 408 Timeout or the request timed out.
	ErrTimeout = errors.New("ensembl: request timeout (408)")

	// ErrRateLimit is returned when Ensembl responds with HTTP 429 Too Many Requests.
	ErrRateLimit = errors.New("ensembl: rate limit exceeded (429)")

	// ErrInternalServer is returned when Ensembl responds with HTTP 500 Internal Server Error.
	ErrInternalServer = errors.New("ensembl: internal server error (500)")

	// ErrServiceUnavailable is returned when Ensembl responds with HTTP 503 Service Unavailable.
	ErrServiceUnavailable = errors.New("ensembl: service unavailable (503)")

	// ErrMaxRetriesReached is returned when the maximum number of retry attempts has been exceeded.
	ErrMaxRetriesReached = errors.New("ensembl: maximum retry attempts reached")
)

// HTTPStatusDescriptions maps Ensembl-specific HTTP status codes to names and descriptions.
var HTTPStatusDescriptions = map[int]struct {
	Name        string
	Description string
}{
	http.StatusOK: {
		Name:        "OK",
		Description: "Request was a success. Only process data from the service when you receive this code",
	},
	http.StatusBadRequest: {
		Name:        "Bad Request",
		Description: "Occurs during exceptional circumstances such as the service is unable to find an ID. If JSON, the object contains the error message",
	},
	http.StatusForbidden: {
		Name:        "Forbidden",
		Description: "You are submitting far too many requests and have been temporarily forbidden access. Wait and retry with a maximum of 15 requests per second",
	},
	http.StatusNotFound: {
		Name:        "Not Found",
		Description: "Indicates a badly formatted request. Check your URL",
	},
	http.StatusRequestTimeout: {
		Name:        "Timeout",
		Description: "The request was not processed in time. Wait and retry later",
	},
	http.StatusUnsupportedMediaType: {
		Name:        "Unsupported Media Type",
		Description: "The server is refusing to service the request because the entity format is not supported",
	},
	http.StatusTooManyRequests: {
		Name:        "Too Many Requests",
		Description: "You have been rate-limited; wait and retry",
	},
	http.StatusInternalServerError: {
		Name:        "Internal Server Error",
		Description: "Internal server error. Check your input or contact the Ensembl team if issue persists",
	},
	http.StatusServiceUnavailable: {
		Name:        "Service Unavailable",
		Description: "The service is temporarily down; retry after a pause",
	},
}

// APIError represents a structured error returned by the Ensembl REST API.
type APIError struct {
	// StatusCode is the HTTP response status code.
	StatusCode int `json:"status_code,omitempty"`

	// Message is the error message returned by the server or client description.
	Message string `json:"error,omitempty"`

	// RateReset is the value of the X-RateLimit-Reset header in epoch seconds, if present.
	RateReset *int `json:"rate_reset,omitempty"`

	// RateLimit is the value of the X-RateLimit-Limit header, if present.
	RateLimit *int `json:"rate_limit,omitempty"`

	// RateRemaining is the value of the X-RateLimit-Remaining header, if present.
	RateRemaining *int `json:"rate_remaining,omitempty"`

	// RetryAfter is the value of the Retry-After header in seconds, if present.
	RetryAfter *float64 `json:"retry_after,omitempty"`

	// RawBody is the raw response body from the server.
	RawBody []byte `json:"-"`
}

// Error formats the APIError into a human-readable string adhering to Go conventions.
func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}

	statusDesc, exists := HTTPStatusDescriptions[e.StatusCode]
	statusName := ""
	if exists {
		statusName = statusDesc.Name
	} else if e.StatusCode > 0 {
		statusName = http.StatusText(e.StatusCode)
	}

	msg := e.Message
	if msg == "" && exists {
		msg = statusDesc.Description
	}

	if e.StatusCode > 0 {
		if e.StatusCode == http.StatusTooManyRequests && e.RetryAfter != nil {
			return fmt.Sprintf("EnsEMBL REST API returned a %d (%s): %s (Rate limit hit: Retry after %d seconds)",
				e.StatusCode, statusName, msg, int(*e.RetryAfter))
		}
		return fmt.Sprintf("EnsEMBL REST API returned a %d (%s): %s", e.StatusCode, statusName, msg)
	}

	return fmt.Sprintf("EnsEMBL REST API error: %s", msg)
}

// Unwrap enables matching against sentinel errors with errors.Is.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusRequestTimeout:
		return ErrTimeout
	case http.StatusTooManyRequests:
		return ErrRateLimit
	case http.StatusInternalServerError:
		return ErrInternalServer
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return nil
	}
}

// parseErrorMessage extracts the error string from JSON if possible.
func parseErrorMessage(body []byte, defaultMsg string) string {
	if len(body) == 0 {
		return defaultMsg
	}

	var jsonErr struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &jsonErr); err == nil {
		if jsonErr.Error != "" {
			return jsonErr.Error
		}
		if jsonErr.Message != "" {
			return jsonErr.Message
		}
	}

	return string(body)
}
