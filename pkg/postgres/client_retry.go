package postgres

import (
	"fmt"
	"time"
)

type RetryConfig struct {
	Config
	MaxAttempts int
	WaitBase    time.Duration
	WaitMax     time.Duration
}

func NewClientWithRetry(cfg RetryConfig) (*Client, error) {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}

	if cfg.WaitBase == 0 {
		cfg.WaitBase = 500 * time.Millisecond
	}

	if cfg.WaitMax == 0 {
		cfg.WaitMax = 5 * time.Second
	}

	wait := cfg.WaitBase
	var lastErr error
	for i := range cfg.MaxAttempts {
		if i > 0 {
			time.Sleep(wait)
			wait *= 2
			if wait > cfg.WaitMax {
				wait = cfg.WaitMax
			}
		}
		client, err := NewClient(cfg.Config)
		if err == nil {
			return client, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("postgres connect after %d attempts: %w", cfg.MaxAttempts, lastErr)
}
