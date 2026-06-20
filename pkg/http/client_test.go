package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	khttp "github.com/KostLabs/klib/pkg/http"
)

func TestClient_Do(t *testing.T) {
	// Given
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})
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

func TestClient_DefaultTimeout(t *testing.T) {
	// Given / When
	c := khttp.NewClient(khttp.Config{})

	// Then
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
