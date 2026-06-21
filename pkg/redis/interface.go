package redis

import (
	"context"
	"time"
)

// Cmdable is the interface implemented by Client and ClientWithRetry.
// Use this interface in constructors to satisfy the dependency-inversion principle.
type Cmdable interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Close() error
}
