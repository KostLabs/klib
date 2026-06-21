package redis

import (
	"context"
	"time"
)

// RetryConfig holds configuration for ClientWithRetry.
type RetryConfig struct {
	Config
	// MaxAttempts is the total number of attempts before giving up. Defaults to 3.
	MaxAttempts int
	// WaitBase is the initial backoff duration between retries. Defaults to 100ms.
	WaitBase time.Duration
	// WaitMax is the maximum backoff duration between retries. Defaults to 2s.
	WaitMax time.Duration
}

// ClientWithRetry wraps Client with exponential-backoff retry logic.
type ClientWithRetry struct {
	inner    *Client
	attempts int
	waitBase time.Duration
	waitMax  time.Duration
}

// NewClientWithRetry returns a ClientWithRetry configured with cfg.
func NewClientWithRetry(cfg RetryConfig) *ClientWithRetry {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}

	if cfg.WaitBase == 0 {
		cfg.WaitBase = 100 * time.Millisecond
	}

	if cfg.WaitMax == 0 {
		cfg.WaitMax = 2 * time.Second
	}

	return &ClientWithRetry{
		inner:    NewClient(cfg.Config),
		attempts: cfg.MaxAttempts,
		waitBase: cfg.WaitBase,
		waitMax:  cfg.WaitMax,
	}
}

// Get returns the string value stored at key, retrying on error with exponential backoff.
func (retryClient *ClientWithRetry) Get(ctx context.Context, key string) (string, error) {
	return retry(retryClient.attempts, retryClient.waitBase, retryClient.waitMax, func() (string, error) {
		return retryClient.inner.Get(ctx, key)
	})
}

// Set stores value at key with the given TTL, retrying on error with exponential backoff.
func (retryClient *ClientWithRetry) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	_, err := retry(retryClient.attempts, retryClient.waitBase, retryClient.waitMax, func() (struct{}, error) {
		return struct{}{}, retryClient.inner.Set(ctx, key, value, ttl)
	})
	return err
}

// Close closes the underlying connection.
func (retryClient *ClientWithRetry) Close() error {
	return retryClient.inner.Close()
}

func retry[T any](attempts int, waitBase, waitMax time.Duration, fn func() (T, error)) (T, error) { //goverifier:ignore:any-type
	var (
		result T
		err    error
	)

	wait := waitBase
	for i := range attempts {
		if i > 0 {
			time.Sleep(wait)
			wait *= 2
			if wait > waitMax {
				wait = waitMax
			}
		}
		result, err = fn()
		if err == nil {
			return result, nil
		}
	}

	return result, err
}
