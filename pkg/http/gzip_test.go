package http_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"testing"

	khttp "github.com/KostLabs/klib/pkg/http"
)

func TestDecompressResponse_PlainBody(t *testing.T) {
	// Given
	body := "hello world"
	resp := &http.Response{
		Header: http.Header{},
		Body:   io.NopCloser(strings.NewReader(body)),
	}

	// When
	reader, err := khttp.DecompressResponse(resp)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer reader.Close()

	got, _ := io.ReadAll(reader)
	if string(got) != body {
		t.Fatalf("expected %q, got %q", body, got)
	}
}

func TestDecompressResponse_GzipBody(t *testing.T) {
	// Given
	original := "compressed content"
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte(original))
	_ = gz.Close()

	resp := &http.Response{
		Header: http.Header{"Content-Encoding": []string{"gzip"}},
		Body:   io.NopCloser(&buf),
	}

	// When
	reader, err := khttp.DecompressResponse(resp)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer reader.Close()

	got, _ := io.ReadAll(reader)
	if string(got) != original {
		t.Fatalf("expected %q, got %q", original, got)
	}
}

func TestDecompressResponse_InvalidGzipBody(t *testing.T) {
	// Given
	resp := &http.Response{
		Header: http.Header{"Content-Encoding": []string{"gzip"}},
		Body:   io.NopCloser(strings.NewReader("not gzip")),
	}

	// When
	_, err := khttp.DecompressResponse(resp)

	// Then
	if err == nil {
		t.Fatal("expected error for invalid gzip, got nil")
	}
}
