package http

import (
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

type Config struct {
	Timeout time.Duration
}

type Client struct {
	inner *http.Client
}

func NewClient(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}

	return &Client{inner: &http.Client{Timeout: cfg.Timeout}}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.inner.Do(req)
}
