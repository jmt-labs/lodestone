package ingest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultMaxRetries     = 3
	defaultInitialBackoff = 200 * time.Millisecond
	defaultMaxBackoff     = 5 * time.Second
)

type httpStatusError struct {
	Source string
	Status int
	Body   string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("%s: status %d: %s", e.Source, e.Status, e.Body)
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var herr *httpStatusError
	if errors.As(err, &herr) {
		return herr.Status >= 500 || herr.Status == http.StatusTooManyRequests
	}
	return true
}

type retryConfig struct {
	maxRetries     int
	initialBackoff time.Duration
	maxBackoff     time.Duration
	sleep          func(time.Duration)
}

func defaultRetryConfig(sleep func(time.Duration)) retryConfig {
	if sleep == nil {
		sleep = time.Sleep
	}
	return retryConfig{
		maxRetries:     defaultMaxRetries,
		initialBackoff: defaultInitialBackoff,
		maxBackoff:     defaultMaxBackoff,
		sleep:          sleep,
	}
}

func retryFetch[T any](ctx context.Context, cfg retryConfig, source string, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error
	backoff := cfg.initialBackoff
	for attempt := 0; attempt < cfg.maxRetries; attempt++ {
		if attempt > 0 {
			if err := ctx.Err(); err != nil {
				return zero, err
			}
			cfg.sleep(backoff)
			backoff *= 2
			if backoff > cfg.maxBackoff {
				backoff = cfg.maxBackoff
			}
		}
		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !isRetryable(err) {
			return zero, err
		}
	}
	return zero, fmt.Errorf("%s: max retries exceeded: %w", source, lastErr)
}
