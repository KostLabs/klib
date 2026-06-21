package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	kredis "github.com/KostLabs/klib/pkg/redis"
)

func newTestClient(t *testing.T) (*kredis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	c := kredis.NewClient(kredis.Config{Addr: mr.Addr()})
	return c, mr
}

func TestClient_SetAndGet(t *testing.T) {
	// Given
	c, _ := newTestClient(t)
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

func TestClient_GetMissingKey(t *testing.T) {
	// Given
	c, _ := newTestClient(t)
	defer c.Close()

	// When
	_, err := c.Get(context.Background(), "missing")

	// Then
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestClient_TTLExpiry(t *testing.T) {
	// Given
	c, mr := newTestClient(t)
	defer c.Close()

	ctx := context.Background()
	if err := c.Set(ctx, "expiring", []byte("val"), time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// When
	mr.FastForward(2 * time.Second)

	// Then
	_, err := c.Get(ctx, "expiring")
	if err == nil {
		t.Fatal("expected error after TTL expiry")
	}
}
