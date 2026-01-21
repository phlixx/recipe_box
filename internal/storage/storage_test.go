package storage

import (
	"os"
	"testing"
)

type testData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func setupTestStorage(t *testing.T) (*Storage, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "recipe_box_storage_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	storage := &Storage{baseDir: tmpDir}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return storage, cleanup
}

func TestStorage_SaveAndLoad(t *testing.T) {
	s, cleanup := setupTestStorage(t)
	defer cleanup()

	data := testData{Name: "test", Value: 42}

	err := s.Save("test.json", data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loaded testData
	err = s.Load("test.json", &loaded)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != "test" || loaded.Value != 42 {
		t.Errorf("expected {test, 42}, got {%s, %d}", loaded.Name, loaded.Value)
	}
}

func TestStorage_Exists(t *testing.T) {
	s, cleanup := setupTestStorage(t)
	defer cleanup()

	if s.Exists("nonexistent.json") {
		t.Error("Exists returned true for nonexistent file")
	}

	s.Save("exists.json", testData{})

	if !s.Exists("exists.json") {
		t.Error("Exists returned false for existing file")
	}
}

func TestStorage_Delete(t *testing.T) {
	s, cleanup := setupTestStorage(t)
	defer cleanup()

	s.Save("todelete.json", testData{})

	if !s.Exists("todelete.json") {
		t.Fatal("file should exist before delete")
	}

	err := s.Delete("todelete.json")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if s.Exists("todelete.json") {
		t.Error("file should not exist after delete")
	}
}

func TestStorage_LoadNonexistent(t *testing.T) {
	s, cleanup := setupTestStorage(t)
	defer cleanup()

	var data testData
	err := s.Load("nonexistent.json", &data)
	if err == nil {
		t.Error("expected error loading nonexistent file")
	}
}

func TestStorage_SaveSlice(t *testing.T) {
	s, cleanup := setupTestStorage(t)
	defer cleanup()

	items := []testData{
		{Name: "first", Value: 1},
		{Name: "second", Value: 2},
	}

	err := s.Save("slice.json", items)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loaded []testData
	err = s.Load("slice.json", &loaded)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("expected 2 items, got %d", len(loaded))
	}
}
