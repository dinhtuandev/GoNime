// Package httpx provides a shared HTTP client with retries and timeouts.
package httpx

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	maxAttempts   = 4
	maxBackoff    = 4 * time.Second
	maxRetryAfter = 15 * time.Second
)

var client = &http.Client{Timeout: 30 * time.Second}

// Option customises a single request.
type Option func(*http.Request)

// WithHeader sets a request header.
func WithHeader(key, value string) Option {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}

// Get performs GET with retries on 429 and 5xx responses. The returned body
// must be closed by the caller. A non-retryable status is returned as-is.
func Get(ctx context.Context, rawURL string, opts ...Option) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		for _, opt := range opts {
			opt(req)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt == maxAttempts-1 {
				break
			}
			if err := sleep(ctx, backoff(attempt, 0)); err != nil {
				return nil, err
			}
			continue
		}

		if !retryable(resp.StatusCode) {
			return resp, nil
		}
		if attempt == maxAttempts-1 {
			return resp, nil
		}

		delay := backoff(attempt, retryAfter(resp))
		resp.Body.Close()
		if err := sleep(ctx, delay); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("http get %s: %w", rawURL, lastErr)
}

func retryable(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500
}

// retryAfter parses the Retry-After header in seconds or HTTP-date form.
func retryAfter(resp *http.Response) time.Duration {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func backoff(attempt int, serverDelay time.Duration) time.Duration {
	d := serverDelay
	if d <= 0 {
		d = 300 * time.Millisecond * time.Duration(math.Pow(2, float64(attempt)))
	}
	if d > maxRetryAfter {
		d = maxRetryAfter
	}
	if d > maxBackoff && serverDelay <= 0 {
		d = maxBackoff
	}
	return d
}

func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// ErrNoResponse reports that every attempt failed before a response arrived.
var ErrNoResponse = errors.New("no response")
