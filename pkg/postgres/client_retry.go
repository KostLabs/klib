package postgres

import (
	"fmt"
	"time"
)

// RetryConfig holds configuration for NewClientWithRetry.
type RetryConfig struct {
	Config
	// MaxAttempts is the total number of connection attempts before giving up. Defaults to 3.
	MaxAttempts int
	// WaitBase is the initial backoff duration between attempts. Defaults to 500ms.
	WaitBase time.Duration
	// WaitMax is the maximum backoff duration between attempts. Defaults to 5s.
	WaitMax time.Duration
}

// NewClientWithRetry attempts to open a PostgreSQL connection up to
// cfg.MaxAttempts times, applying exponential backoff between attempts.
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
