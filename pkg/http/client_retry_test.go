package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	khttp "github.com/KostLabs/klib/pkg/http"
)

func TestClientWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	// Given
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := khttp.NewClientWithRetry(khttp.RetryConfig{
		Config:      khttp.Config{Timeout: 5 * time.Second},
		MaxAttempts: 3,
		WaitBase:    time.Millisecond,
	})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// When
	resp, err := c.Do(req)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestClientWithRetry_RetriesOnConfiguredStatusCode(t *testing.T) {
	// Given
	var calls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := khttp.NewClientWithRetry(khttp.RetryConfig{
		Config:      khttp.Config{Timeout: 5 * time.Second},
		MaxAttempts: 3,
		WaitBase:    time.Millisecond,
		WaitMax:     10 * time.Millisecond,
	})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// When
	resp, err := c.Do(req)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 after retries, got %d", resp.StatusCode)
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestClientWithRetry_ExhaustsAttempts(t *testing.T) {
	// Given
	var calls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := khttp.NewClientWithRetry(khttp.RetryConfig{
		Config:      khttp.Config{Timeout: 5 * time.Second},
		MaxAttempts: 3,
		WaitBase:    time.Millisecond,
		WaitMax:     10 * time.Millisecond,
	})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// When
	resp, err := c.Do(req)

	// Then
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502 after exhausted retries, got %d", resp.StatusCode)
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestClientWithRetry_DoesNotRetryOnSuccess(t *testing.T) {
	// Given
	var calls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	c := khttp.NewClientWithRetry(khttp.RetryConfig{
		Config:      khttp.Config{Timeout: 5 * time.Second},
		MaxAttempts: 3,
		WaitBase:    time.Millisecond,
	})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// When
	resp, err := c.Do(req)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if calls.Load() != 1 {
		t.Fatalf("expected exactly 1 call, got %d", calls.Load())
	}
}
