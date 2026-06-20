package redis

import (
	"context"
	"time"
)

type RetryConfig struct {
	Config
	MaxAttempts int
	WaitBase    time.Duration
	WaitMax     time.Duration
}

type ClientWithRetry struct {
	inner    *Client
	attempts int
	waitBase time.Duration
	waitMax  time.Duration
}

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

func (c *ClientWithRetry) Get(ctx context.Context, key string) (string, error) {
	return retry(c.attempts, c.waitBase, c.waitMax, func() (string, error) {
		return c.inner.Get(ctx, key)
	})
}

func (c *ClientWithRetry) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := retry(c.attempts, c.waitBase, c.waitMax, func() (struct{}, error) {
		return struct{}{}, c.inner.Set(ctx, key, value, ttl)
	})
	return err
}

func (c *ClientWithRetry) Close() error {
	return c.inner.Close()
}

func retry[T any](attempts int, waitBase, waitMax time.Duration, fn func() (T, error)) (T, error) {
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
