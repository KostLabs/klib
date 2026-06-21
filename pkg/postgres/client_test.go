package postgres_test

import (
	"os"
	"strconv"
	"testing"

	kpostgres "github.com/KostLabs/klib/pkg/postgres"
)

// testConfig reads connection details from env vars.
// Tests are skipped when POSTGRES_HOST is not set.
func testConfig(t *testing.T) kpostgres.Config {
	t.Helper()
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		t.Skip("POSTGRES_HOST not set, skipping postgres integration tests")
	}
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if port == 0 {
		port = 5432
	}
	return kpostgres.Config{
		Host:     host,
		Port:     port,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
}

func TestClient_Connect(t *testing.T) {
	// Given
	cfg := testConfig(t)

	// When
	c, err := kpostgres.NewClient(cfg)

	// Then
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer c.Close()

	if err := c.DB().Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestClient_InvalidAddress(t *testing.T) {
	// Given
	cfg := kpostgres.Config{
		Host:   "127.0.0.1",
		Port:   1,
		User:   "nobody",
		DBName: "none",
	}

	// When
	_, err := kpostgres.NewClient(cfg)

	// Then
	if err == nil {
		t.Fatal("expected error connecting to invalid address")
	}
}
