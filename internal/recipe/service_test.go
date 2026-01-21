package recipe

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestService_AddAndList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "recipe_box_service_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := newTestService(tmpDir)

	r := &Recipe{
		Title:    "Test Recipe",
		Servings: 4,
		Ingredients: []Ingredient{
			{Name: "flour", Quantity: 2, Unit: "cups"},
		},
		Steps: []string{"Mix ingredients", "Bake"},
	}

	err = svc.Add(r)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if r.ID == "" {
		t.Error("expected ID to be generated")
	}

	recipes, err := svc.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(recipes) != 1 {
		t.Errorf("expected 1 recipe, got %d", len(recipes))
	}

	if recipes[0].Title != "Test Recipe" {
		t.Errorf("expected title 'Test Recipe', got '%s'", recipes[0].Title)
	}
}

func TestService_Get(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "recipe_box_service_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := newTestService(tmpDir)

	r := &Recipe{Title: "Get Test"}
	svc.Add(r)

	found, err := svc.Get(r.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if found.Title != "Get Test" {
		t.Errorf("expected title 'Get Test', got '%s'", found.Title)
	}
}

func TestService_GetNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "recipe_box_service_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := newTestService(tmpDir)

	_, err = svc.Get("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestService_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "recipe_box_service_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := newTestService(tmpDir)

	r := &Recipe{Title: "To Delete"}
	svc.Add(r)

	err = svc.Delete(r.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = svc.Get(r.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestService_DeleteNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "recipe_box_service_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := newTestService(tmpDir)

	err = svc.Delete("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestService_ListEmpty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "recipe_box_service_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := newTestService(tmpDir)

	recipes, err := svc.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(recipes) != 0 {
		t.Errorf("expected empty list, got %d items", len(recipes))
	}
}

func newTestService(baseDir string) *Service {
	store := &testStorage{baseDir: baseDir}
	return &Service{storage: store}
}

// testStorage implements Storage interface for testing
type testStorage struct {
	baseDir string
}

func (s *testStorage) Save(filename string, data any) error {
	path := filepath.Join(s.baseDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (s *testStorage) Load(filename string, dest any) error {
	path := filepath.Join(s.baseDir, filename)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(dest)
}

func (s *testStorage) Exists(filename string) bool {
	path := filepath.Join(s.baseDir, filename)
	_, err := os.Stat(path)
	return err == nil
}

func (s *testStorage) Delete(filename string) error {
	path := filepath.Join(s.baseDir, filename)
	return os.Remove(path)
}
