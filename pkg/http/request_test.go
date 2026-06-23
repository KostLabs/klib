package http_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	khttp "github.com/KostLabs/klib/pkg/http"
)

type payload struct {
	Name string `json:"name"`
}

func gzipJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, _ := json.Marshal(v)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write(b)
	_ = gz.Close()
	return buf.Bytes()
}

func TestPostJSON_Success(t *testing.T) {
	// Given
	want := payload{Name: "result"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.PostJSON(context.Background(), srv.URL, payload{Name: "input"}, &got)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != want.Name {
		t.Fatalf("expected name=%q, got %q", want.Name, got.Name)
	}
}

func TestPostJSON_GzipResponse(t *testing.T) {
	// Given
	want := payload{Name: "gzipped"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(gzipJSON(t, want))
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.PostJSON(context.Background(), srv.URL, payload{}, &got)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != want.Name {
		t.Fatalf("expected name=%q, got %q", want.Name, got.Name)
	}
}

func TestPostJSON_NonSuccessStatus(t *testing.T) {
	// Given
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad"}`))
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.PostJSON(context.Background(), srv.URL, payload{}, &got)

	// Then
	var statusErr khttp.ErrUnexpectedStatus
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected ErrUnexpectedStatus, got %T: %v", err, err)
	}
	if statusErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusErr.StatusCode)
	}
}

func TestGetJSON_Success(t *testing.T) {
	// Given
	want := payload{Name: "fetched"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.GetJSON(context.Background(), srv.URL, &got)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != want.Name {
		t.Fatalf("expected name=%q, got %q", want.Name, got.Name)
	}
}

func TestGetJSON_GzipResponse(t *testing.T) {
	// Given
	want := payload{Name: "gzip-get"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(gzipJSON(t, want))
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.GetJSON(context.Background(), srv.URL, &got)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != want.Name {
		t.Fatalf("expected name=%q, got %q", want.Name, got.Name)
	}
}

func TestGetJSON_NonSuccessStatus(t *testing.T) {
	// Given
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.GetJSON(context.Background(), srv.URL, &got)

	// Then
	var statusErr khttp.ErrUnexpectedStatus
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected ErrUnexpectedStatus, got %T: %v", err, err)
	}
	if statusErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusErr.StatusCode)
	}
}

func TestGetJSON_AppliesRequestOptions(t *testing.T) {
	// Given
	var receivedToken string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken = r.Header.Get("Authorization")
		_ = json.NewEncoder(w).Encode(payload{})
	}))
	defer srv.Close()

	c := khttp.NewClient(khttp.Config{Timeout: 5 * time.Second})

	// When
	var got payload
	err := c.GetJSON(context.Background(), srv.URL, &got, khttp.WithBearerToken("tok123"))

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedToken != "Bearer tok123" {
		t.Fatalf("expected Authorization=Bearer tok123, got %q", receivedToken)
	}
}
