package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	kredis "github.com/KostLabs/klib/pkg/redis"
)

func newTestRetryClient(t *testing.T) (*kredis.ClientWithRetry, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	c := kredis.NewClientWithRetry(kredis.RetryConfig{
		Config:      kredis.Config{Addr: mr.Addr()},
		MaxAttempts: 3,
		WaitBase:    time.Millisecond,
		WaitMax:     5 * time.Millisecond,
	})
	return c, mr
}

func TestClientWithRetry_SetAndGet(t *testing.T) {
	// Given
	c, _ := newTestRetryClient(t)
	defer c.Close()

	ctx := context.Background()

	// When
	err := c.Set(ctx, "key", []byte("value"), time.Minute)

	// Then
	if err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := c.Get(ctx, "key")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "value" {
		t.Fatalf("expected %q, got %q", "value", got)
	}
}

func TestClientWithRetry_RetriesOnTransientError(t *testing.T) {
	// Given
	mr := miniredis.RunT(t)
	c := kredis.NewClientWithRetry(kredis.RetryConfig{
		Config:      kredis.Config{Addr: mr.Addr()},
		MaxAttempts: 3,
		WaitBase:    time.Millisecond,
		WaitMax:     5 * time.Millisecond,
	})
	defer c.Close()

	ctx := context.Background()
	if err := c.Set(ctx, "k", []byte("v"), time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// When
	mr.SetError("ERR server busy")
	go func() {
		time.Sleep(3 * time.Millisecond)
		mr.SetError("")
	}()

	// Then
	got, err := c.Get(ctx, "k")
	if err != nil {
		t.Fatalf("Get after retry: %v", err)
	}
	if got != "v" {
		t.Fatalf("expected %q, got %q", "v", got)
	}
}
