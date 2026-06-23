package http_test

import (
	"net/http"
	"testing"

	khttp "github.com/KostLabs/klib/pkg/http"
)

func TestWithHeaders_SetsAllHeaders(t *testing.T) {
	// Given
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	opt := khttp.WithHeaders(map[string]string{
		"X-Foo": "bar",
		"X-Baz": "qux",
	})

	// When
	opt(req)

	// Then
	if got := req.Header.Get("X-Foo"); got != "bar" {
		t.Fatalf("expected X-Foo=bar, got %q", got)
	}
	if got := req.Header.Get("X-Baz"); got != "qux" {
		t.Fatalf("expected X-Baz=qux, got %q", got)
	}
}

func TestWithHeaders_OverwritesExistingHeader(t *testing.T) {
	// Given
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Set("X-Foo", "old")
	opt := khttp.WithHeaders(map[string]string{"X-Foo": "new"})

	// When
	opt(req)

	// Then
	if got := req.Header.Get("X-Foo"); got != "new" {
		t.Fatalf("expected X-Foo=new, got %q", got)
	}
}

func TestWithBearerToken_SetsAuthorizationHeader(t *testing.T) {
	// Given
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	opt := khttp.WithBearerToken("my-secret-token")

	// When
	opt(req)

	// Then
	if got := req.Header.Get("Authorization"); got != "Bearer my-secret-token" {
		t.Fatalf("expected Bearer my-secret-token, got %q", got)
	}
}

func TestWithBearerToken_OverwritesExistingToken(t *testing.T) {
	// Given
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer old-token")
	opt := khttp.WithBearerToken("new-token")

	// When
	opt(req)

	// Then
	if got := req.Header.Get("Authorization"); got != "Bearer new-token" {
		t.Fatalf("expected Bearer new-token, got %q", got)
	}
}
