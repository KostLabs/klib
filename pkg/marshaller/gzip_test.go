package marshaller_test

import (
	"testing"

	"github.com/KostLabs/klib/pkg/marshaller"
)

type sample struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestMarshalGzip_RoundTrip(t *testing.T) {
	// Given
	in := sample{ID: 1, Name: "test"}

	// When
	data, err := marshaller.MarshalGzip(in)

	// Then
	if err != nil {
		t.Fatalf("MarshalGzip: unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("MarshalGzip: expected non-empty output")
	}

	var out sample
	if err := marshaller.UnmarshalGzip(data, &out); err != nil {
		t.Fatalf("UnmarshalGzip: unexpected error: %v", err)
	}
	if out.ID != in.ID || out.Name != in.Name {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", out, in)
	}
}

func TestMarshalGzip_ProducesCompressedBytes(t *testing.T) {
	// Given
	in := sample{ID: 42, Name: "compressed"}

	// When
	data, err := marshaller.MarshalGzip(in)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// gzip magic bytes: 0x1f 0x8b
	if len(data) < 2 || data[0] != 0x1f || data[1] != 0x8b {
		t.Fatal("expected output to start with gzip magic bytes")
	}
}

func TestUnmarshalGzip_InvalidData(t *testing.T) {
	// Given
	garbage := []byte("this is not gzip data")

	// When
	var out sample
	err := marshaller.UnmarshalGzip(garbage, &out)

	// Then
	if err == nil {
		t.Fatal("expected error for invalid gzip, got nil")
	}
}

func TestMarshalGzip_EmptySlice(t *testing.T) {
	// Given
	in := []string{}

	// When
	data, err := marshaller.MarshalGzip(in)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out []string
	if err := marshaller.UnmarshalGzip(data, &out); err != nil {
		t.Fatalf("UnmarshalGzip: unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %v", out)
	}
}
