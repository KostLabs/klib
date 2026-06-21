package http

import (
	"net/http"
	"time"
)

// RetryConfig holds configuration for ClientWithRetry.
type RetryConfig struct {
	Config
	// MaxAttempts is the total number of attempts before giving up. Defaults to 3.
	MaxAttempts int
	// WaitBase is the initial backoff duration between retries. Defaults to 500ms.
	WaitBase time.Duration
	// WaitMax is the maximum backoff duration between retries. Defaults to 5s.
	WaitMax time.Duration
	// RetryOn is the set of HTTP status codes that trigger a retry.
	// Defaults to 429, 502, 503, 504 when nil.
	RetryOn []int
}

// ClientWithRetry wraps Client with exponential-backoff retry logic.
type ClientWithRetry struct {
	inner    *http.Client
	attempts int
	waitBase time.Duration
	waitMax  time.Duration
	retryOn  map[int]struct{}
}

// NewClientWithRetry returns a ClientWithRetry configured with cfg.
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

// Do executes the request, retrying on configured status codes with exponential backoff.
func (retryClient *ClientWithRetry) Do(req *http.Request) (*http.Response, error) {
	wait := retryClient.waitBase
	var (
		resp *http.Response
		err  error
	)

	for attempt := range retryClient.attempts {
		if attempt > 0 {
			time.Sleep(wait)
			wait *= 2
			if wait > retryClient.waitMax {
				wait = retryClient.waitMax
			}
		}

		resp, err = retryClient.inner.Do(req)
		if err != nil {
			continue
		}
		if _, shouldRetry := retryClient.retryOn[resp.StatusCode]; !shouldRetry {
			return resp, nil
		}

		resp.Body.Close()
	}

	return resp, err
}
