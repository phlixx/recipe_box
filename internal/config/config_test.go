package config

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestConfig(t *testing.T) (*Config, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "recipe_box_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	configPath := filepath.Join(tmpDir, configFile)
	cfg := &Config{
		path: configPath,
		data: make(map[string]string),
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return cfg, cleanup
}

func TestConfig_SetAndGet(t *testing.T) {
	cfg, cleanup := setupTestConfig(t)
	defer cleanup()

	err := cfg.Set("api_key", "test-key-123")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	value, err := cfg.Get("api_key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value != "test-key-123" {
		t.Errorf("expected 'test-key-123', got '%s'", value)
	}
}

func TestConfig_GetMissingKey(t *testing.T) {
	cfg, cleanup := setupTestConfig(t)
	defer cleanup()

	_, err := cfg.Get("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestConfig_Persistence(t *testing.T) {
	cfg, cleanup := setupTestConfig(t)
	defer cleanup()

	err := cfg.Set("language", "de")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Create new config instance pointing to same file
	cfg2 := &Config{
		path: cfg.path,
		data: make(map[string]string),
	}
	if err := cfg2.load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	value, err := cfg2.Get("language")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value != "de" {
		t.Errorf("expected 'de', got '%s'", value)
	}
}

func TestConfig_OverwriteValue(t *testing.T) {
	cfg, cleanup := setupTestConfig(t)
	defer cleanup()

	cfg.Set("key", "value1")
	cfg.Set("key", "value2")

	value, _ := cfg.Get("key")
	if value != "value2" {
		t.Errorf("expected 'value2', got '%s'", value)
	}
}
