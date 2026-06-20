package http

import (
	"net/http"
	"time"
)

type RetryConfig struct {
	Config
	MaxAttempts int
	WaitBase    time.Duration
	WaitMax     time.Duration
	// RetryOn is a set of HTTP status codes that trigger a retry.
	// Defaults to 429, 502, 503, 504 when nil.
	RetryOn []int
}

type ClientWithRetry struct {
	inner    *http.Client
	attempts int
	waitBase time.Duration
	waitMax  time.Duration
	retryOn  map[int]struct{}
}

func NewClientWithRetry(cfg RetryConfig) *ClientWithRetry {
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}

	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}

	if cfg.WaitBase == 0 {
		cfg.WaitBase = 500 * time.Millisecond
	}

	if cfg.WaitMax == 0 {
		cfg.WaitMax = 5 * time.Second
	}

	codes := cfg.RetryOn
	if len(codes) == 0 {
		codes = []int{http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}
	}

	retryOn := make(map[int]struct{}, len(codes))
	for _, code := range codes {
		retryOn[code] = struct{}{}
	}

	return &ClientWithRetry{
		inner:    &http.Client{Timeout: cfg.Timeout},
		attempts: cfg.MaxAttempts,
		waitBase: cfg.WaitBase,
		waitMax:  cfg.WaitMax,
		retryOn:  retryOn,
	}
}

func (c *ClientWithRetry) Do(req *http.Request) (*http.Response, error) {
	wait := c.waitBase
	var (
		resp *http.Response
		err  error
	)

	for attempt := range c.attempts {
		if attempt > 0 {
			time.Sleep(wait)
			wait *= 2
			if wait > c.waitMax {
				wait = c.waitMax
			}
		}

		resp, err = c.inner.Do(req)
		if err != nil {
			continue
		}
		if _, shouldRetry := c.retryOn[resp.StatusCode]; !shouldRetry {
			return resp, nil
		}
		resp.Body.Close()
	}

	return resp, err
}
