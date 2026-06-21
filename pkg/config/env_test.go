package config_test

import (
	"testing"
	"time"

	"github.com/KostLabs/klib/pkg/config"
)

func TestLoadEnv_StringField(t *testing.T) {
	// Given
	type cfg struct {
		APIKey string `env:"TEST_API_KEY"`
	}

	t.Setenv("TEST_API_KEY", "secret")
	target := cfg{}

	// When
	err := config.LoadEnv(&target)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target.APIKey != "secret" {
		t.Errorf("expected 'secret', got %q", target.APIKey)
	}
}

func TestLoadEnv_IntField(t *testing.T) {
	// Given
	type cfg struct {
		Port int `env:"TEST_PORT"`
	}

	t.Setenv("TEST_PORT", "8080")
	target := cfg{}

	// When
	err := config.LoadEnv(&target)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target.Port != 8080 {
		t.Errorf("expected 8080, got %d", target.Port)
	}
}

func TestLoadEnv_BoolField(t *testing.T) {
	// Given
	type cfg struct {
		Debug bool `env:"TEST_DEBUG"`
	}

	t.Setenv("TEST_DEBUG", "true")
	target := cfg{}

	// When
	err := config.LoadEnv(&target)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !target.Debug {
		t.Error("expected Debug to be true")
	}
}

func TestLoadEnv_DurationField(t *testing.T) {
	// Given
	type cfg struct {
		Timeout time.Duration `env:"TEST_TIMEOUT"`
	}

	t.Setenv("TEST_TIMEOUT", "30s")
	target := cfg{}

	// When
	err := config.LoadEnv(&target)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target.Timeout != 30*time.Second {
		t.Errorf("expected 30s, got %v", target.Timeout)
	}
}

func TestLoadEnv_NestedStruct(t *testing.T) {
	// Given
	type redisCfg struct {
		Addr     string `env:"TEST_REDIS_ADDR"`
		Password string `env:"TEST_REDIS_PASSWORD"`
	}
	type cfg struct {
		Redis redisCfg
	}

	t.Setenv("TEST_REDIS_ADDR", "localhost:6379")
	t.Setenv("TEST_REDIS_PASSWORD", "hunter2")
	target := cfg{}

	// When
	err := config.LoadEnv(&target)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target.Redis.Addr != "localhost:6379" {
		t.Errorf("expected 'localhost:6379', got %q", target.Redis.Addr)
	}
	if target.Redis.Password != "hunter2" {
		t.Errorf("expected 'hunter2', got %q", target.Redis.Password)
	}
}

func TestLoadEnv_UnsetVarPreservesDefault(t *testing.T) {
	// Given
	type cfg struct {
		APIKey string `env:"TEST_UNSET_KEY_XYZ"`
	}

	target := cfg{APIKey: "default"}

	// When
	err := config.LoadEnv(&target)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target.APIKey != "default" {
		t.Errorf("expected default to be preserved, got %q", target.APIKey)
	}
}

func TestLoadEnv_MultipleTargets(t *testing.T) {
	// Given
	type appCfg struct {
		APIKey string `env:"TEST_MULTI_API_KEY"`
	}
	type dbCfg struct {
		DSN string `env:"TEST_MULTI_DSN"`
	}

	t.Setenv("TEST_MULTI_API_KEY", "key123")
	t.Setenv("TEST_MULTI_DSN", "postgres://localhost/db")

	appTarget := appCfg{}
	dbTarget := dbCfg{}

	// When
	err := config.LoadEnv(&appTarget, &dbTarget)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appTarget.APIKey != "key123" {
		t.Errorf("expected 'key123', got %q", appTarget.APIKey)
	}
	if dbTarget.DSN != "postgres://localhost/db" {
		t.Errorf("expected 'postgres://localhost/db', got %q", dbTarget.DSN)
	}
}

func TestLoadEnv_NonPointerReturnsError(t *testing.T) {
	// Given
	type cfg struct {
		APIKey string `env:"TEST_API_KEY"`
	}

	// When
	err := config.LoadEnv(cfg{})

	// Then
	if err == nil {
		t.Fatal("expected error for non-pointer target")
	}
}
