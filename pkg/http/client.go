// Package http provides a thin wrapper around net/http with configurable
// timeout and optional exponential-backoff retry logic.
package http

import (
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

// Config holds configuration for Client.
type Config struct {
	// Timeout is the maximum time allowed for a single request.
	// Defaults to 30 seconds when zero.
	Timeout time.Duration
}

// Client is a thin wrapper around net/http.Client.
type Client struct {
	inner *http.Client
}

// NewClient returns a Client configured with cfg.
func NewClient(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}

	return &Client{inner: &http.Client{Timeout: cfg.Timeout}}
}

// Do executes the request and returns the response.
// Any provided opts are applied to req before the call is made.
func (client *Client) Do(req *http.Request, opts ...RequestOption) (*http.Response, error) {
	applyOptions(req, opts)
	return client.inner.Do(req)
}

// applyOptions applies all opts to req in order.
func applyOptions(req *http.Request, opts []RequestOption) {
	for _, opt := range opts {
		opt(req)
	}
}
