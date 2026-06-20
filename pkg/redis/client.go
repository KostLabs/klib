package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Client struct {
	inner *redis.Client
}

func NewClient(cfg Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Client{inner: rdb}
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.inner.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.inner.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Close() error {
	return c.inner.Close()
}
