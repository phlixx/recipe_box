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

func TestRecipe_Scale(t *testing.T) {
	tests := []struct {
		name           string
		original       *Recipe
		targetServings int
		wantServings   int
		wantQuantities []float64
	}{
		{
			name: "scale up from 4 to 8 servings",
			original: &Recipe{
				Title:    "Test Recipe",
				Servings: 4,
				Ingredients: []Ingredient{
					{Name: "flour", Quantity: 2, Unit: "cups"},
					{Name: "sugar", Quantity: 1, Unit: "cup"},
				},
				Steps: []string{"Mix", "Bake"},
				Tags:  []string{"dessert"},
			},
			targetServings: 8,
			wantServings:   8,
			wantQuantities: []float64{4, 2},
		},
		{
			name: "scale down from 4 to 2 servings",
			original: &Recipe{
				Title:    "Test Recipe",
				Servings: 4,
				Ingredients: []Ingredient{
					{Name: "flour", Quantity: 2, Unit: "cups"},
					{Name: "sugar", Quantity: 1, Unit: "cup"},
				},
			},
			targetServings: 2,
			wantServings:   2,
			wantQuantities: []float64{1, 0.5},
		},
		{
			name: "same servings returns copy",
			original: &Recipe{
				Title:    "Test Recipe",
				Servings: 4,
				Ingredients: []Ingredient{
					{Name: "flour", Quantity: 2, Unit: "cups"},
				},
			},
			targetServings: 4,
			wantServings:   4,
			wantQuantities: []float64{2},
		},
		{
			name: "zero target returns original",
			original: &Recipe{
				Title:    "Test Recipe",
				Servings: 4,
				Ingredients: []Ingredient{
					{Name: "flour", Quantity: 2, Unit: "cups"},
				},
			},
			targetServings: 0,
			wantServings:   4,
			wantQuantities: []float64{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaled := tt.original.Scale(tt.targetServings)

			if scaled.Servings != tt.wantServings {
				t.Errorf("Servings = %d, want %d", scaled.Servings, tt.wantServings)
			}

			for i, want := range tt.wantQuantities {
				if scaled.Ingredients[i].Quantity != want {
					t.Errorf("Ingredient[%d].Quantity = %v, want %v", i, scaled.Ingredients[i].Quantity, want)
				}
			}

			// Verify original is not modified
			if tt.original.Servings != 4 {
				t.Error("original recipe was modified")
			}
		})
	}
}

func TestRecipe_Scale_PreservesOtherFields(t *testing.T) {
	original := &Recipe{
		ID:          "abc123",
		Title:       "Test Recipe",
		Description: "A test",
		Servings:    4,
		PrepTime:    10,
		CookTime:    20,
		Ingredients: []Ingredient{
			{Name: "flour", Quantity: 2, Unit: "cups", Category: "pantry"},
		},
		Steps:    []string{"Step 1", "Step 2"},
		Tags:     []string{"test", "dessert"},
		Source:   SourceManual,
		Language: LangEN,
		AsciiArt: "art",
	}

	scaled := original.Scale(8)

	if scaled.ID != original.ID {
		t.Error("ID not preserved")
	}
	if scaled.Title != original.Title {
		t.Error("Title not preserved")
	}
	if scaled.Description != original.Description {
		t.Error("Description not preserved")
	}
	if scaled.PrepTime != original.PrepTime {
		t.Error("PrepTime not preserved")
	}
	if scaled.CookTime != original.CookTime {
		t.Error("CookTime not preserved")
	}
	if scaled.Source != original.Source {
		t.Error("Source not preserved")
	}
	if scaled.Language != original.Language {
		t.Error("Language not preserved")
	}
	if scaled.AsciiArt != original.AsciiArt {
		t.Error("AsciiArt not preserved")
	}
	if len(scaled.Steps) != len(original.Steps) {
		t.Error("Steps length mismatch")
	}
	if len(scaled.Tags) != len(original.Tags) {
		t.Error("Tags length mismatch")
	}
}
