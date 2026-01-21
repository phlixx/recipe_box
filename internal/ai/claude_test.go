package ai

import (
	"strings"
	"testing"

	"recipe_box/internal/recipe"
)

func TestNewClient_NoAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err != ErrNoAPIKey {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}

func TestNewClient_WithAPIKey(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.apiKey != "test-key" {
		t.Errorf("expected apiKey 'test-key', got '%s'", client.apiKey)
	}
}

func TestBuildPrompt_Basic(t *testing.T) {
	opts := GenerateOptions{
		Prompt: "pasta dish",
	}
	prompt := buildPrompt(opts)

	if !strings.Contains(prompt, "pasta dish") {
		t.Error("expected prompt to contain user request")
	}
	if !strings.Contains(prompt, "Return ONLY valid JSON") {
		t.Error("expected prompt to contain JSON instructions")
	}
}

func TestBuildPrompt_WithConstraints(t *testing.T) {
	opts := GenerateOptions{
		Prompt:      "dinner",
		Cuisine:     "italian",
		Vegetarian:  true,
		Quick:       true,
		Ingredients: []string{"tomato", "basil"},
	}
	prompt := buildPrompt(opts)

	if !strings.Contains(prompt, "italian") {
		t.Error("expected prompt to contain cuisine")
	}
	if !strings.Contains(prompt, "vegetarian") {
		t.Error("expected prompt to contain vegetarian constraint")
	}
	if !strings.Contains(prompt, "under 30 minutes") {
		t.Error("expected prompt to contain quick constraint")
	}
	if !strings.Contains(prompt, "tomato") || !strings.Contains(prompt, "basil") {
		t.Error("expected prompt to contain ingredients")
	}
}

func TestBuildPrompt_DefaultRequest(t *testing.T) {
	opts := GenerateOptions{}
	prompt := buildPrompt(opts)

	if !strings.Contains(prompt, "a delicious recipe") {
		t.Error("expected default request when no prompt provided")
	}
}

func TestParseRecipeFromResponse_ValidJSON(t *testing.T) {
	jsonResponse := `{
		"title": "Test Recipe",
		"description": "A test description",
		"servings": 4,
		"prep_time": 10,
		"cook_time": 20,
		"ingredients": [
			{"name": "flour", "quantity": 2, "unit": "cups", "category": "pantry"}
		],
		"steps": ["Step 1", "Step 2"],
		"tags": ["quick", "easy"]
	}`

	r, err := parseRecipeFromResponse(jsonResponse, "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Title != "Test Recipe" {
		t.Errorf("expected title 'Test Recipe', got '%s'", r.Title)
	}
	if r.Servings != 4 {
		t.Errorf("expected servings 4, got %d", r.Servings)
	}
	if len(r.Ingredients) != 1 {
		t.Errorf("expected 1 ingredient, got %d", len(r.Ingredients))
	}
	if len(r.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(r.Steps))
	}
	if r.Source != recipe.SourceAI {
		t.Errorf("expected source 'ai', got '%s'", r.Source)
	}
	if r.Language != recipe.LangEN {
		t.Errorf("expected language 'en', got '%s'", r.Language)
	}
}

func TestParseRecipeFromResponse_MarkdownWrapped(t *testing.T) {
	jsonResponse := "```json\n" + `{
		"title": "Wrapped Recipe",
		"description": "Test",
		"servings": 2,
		"prep_time": 5,
		"cook_time": 10,
		"ingredients": [],
		"steps": ["Do it"],
		"tags": []
	}` + "\n```"

	r, err := parseRecipeFromResponse(jsonResponse, "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Title != "Wrapped Recipe" {
		t.Errorf("expected title 'Wrapped Recipe', got '%s'", r.Title)
	}
}

func TestParseRecipeFromResponse_InvalidJSON(t *testing.T) {
	_, err := parseRecipeFromResponse("not json at all", "en")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseRecipeFromResponse_GermanLanguage(t *testing.T) {
	jsonResponse := `{
		"title": "Test Rezept",
		"description": "Ein Test",
		"servings": 4,
		"prep_time": 10,
		"cook_time": 20,
		"ingredients": [],
		"steps": ["Schritt 1"],
		"tags": []
	}`

	r, err := parseRecipeFromResponse(jsonResponse, "de")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Language != recipe.LangDE {
		t.Errorf("expected language 'de', got '%s'", r.Language)
	}
}

func TestBuildPrompt_GermanLanguage(t *testing.T) {
	opts := GenerateOptions{
		Prompt:   "pasta dish",
		Language: "de",
	}
	prompt := buildPrompt(opts)

	if !strings.Contains(prompt, "German") {
		t.Error("expected prompt to contain German language instruction")
	}
}

func TestBuildPrompt_EnglishLanguage(t *testing.T) {
	opts := GenerateOptions{
		Prompt:   "pasta dish",
		Language: "en",
	}
	prompt := buildPrompt(opts)

	if !strings.Contains(prompt, "English") {
		t.Error("expected prompt to contain English language instruction")
	}
}
