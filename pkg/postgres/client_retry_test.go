package postgres_test

import (
	"testing"
	"time"

	kpostgres "github.com/KostLabs/klib/pkg/postgres"
)

func TestClientWithRetry_Connect(t *testing.T) {
	// Given
	cfg := testConfig(t)

	// When
	c, err := kpostgres.NewClientWithRetry(kpostgres.RetryConfig{
		Config:      cfg,
		MaxAttempts: 3,
		WaitBase:    10 * time.Millisecond,
		WaitMax:     50 * time.Millisecond,
	})

	// Then
	if err != nil {
		t.Fatalf("NewClientWithRetry: %v", err)
	}
	defer c.Close()

	if err := c.DB().Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestClientWithRetry_ExhaustsAttemptsOnBadAddress(t *testing.T) {
	// Given
	cfg := kpostgres.RetryConfig{
		Config: kpostgres.Config{
			Host:   "127.0.0.1",
			Port:   1,
			User:   "nobody",
			DBName: "none",
		},
		MaxAttempts: 2,
		WaitBase:    time.Millisecond,
		WaitMax:     5 * time.Millisecond,
	}

	// When
	_, err := kpostgres.NewClientWithRetry(cfg)

	// Then
	if err == nil {
		t.Fatal("expected error after exhausted attempts")
	}
}
