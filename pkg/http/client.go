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
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.inner.Do(req)
}
