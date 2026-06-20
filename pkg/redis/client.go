// Package redis provides a thin wrapper around go-redis with configurable
// retry logic via ClientWithRetry.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds connection parameters for a Redis client.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Client is a thin wrapper around go-redis.
type Client struct {
	inner *redis.Client
}

// NewClient returns a Client connected to the Redis instance described by cfg.
func NewClient(cfg Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Client{inner: rdb}
}

// Get returns the string value stored at key, or an error if the key does not exist.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.inner.Get(ctx, key).Result()
}

// Set stores value at key with the given TTL. A TTL of 0 means no expiry.
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.inner.Set(ctx, key, value, ttl).Err()
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.inner.Close()
}
