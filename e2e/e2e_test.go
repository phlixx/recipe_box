package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary before running tests
	cmd := exec.Command("go", "build", "-o", "recipe_box_test", "..")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	// Get absolute path to binary
	var err error
	binaryPath, err = filepath.Abs("recipe_box_test")
	if err != nil {
		panic("failed to get binary path: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(binaryPath)

	os.Exit(code)
}

// runCmd executes the CLI with the given args and returns stdout, stderr, and error
func runCmd(t *testing.T, homeDir string, args ...string) (string, string, error) {
	t.Helper()
	return runCmdWithInput(t, homeDir, "", args...)
}

// runCmdWithInput executes the CLI with stdin input
func runCmdWithInput(t *testing.T, homeDir string, input string, args ...string) (string, string, error) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "RECIPE_BOX_HOME="+homeDir, "NO_COLOR=1")

	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// setupTestHome creates a temp directory for test isolation
func setupTestHome(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "recipe_box_e2e_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// extractRecipeID extracts the recipe ID from "Recipe added: Title (ID)" output
// Works with both English and German messages
func extractRecipeID(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Look for lines containing "added" or "hinzugefügt" (the success messages)
		if strings.Contains(line, "added") || strings.Contains(line, "hinzugefügt") {
			// Find the last occurrence of (xxx) pattern which contains the ID
			start := strings.LastIndex(line, "(")
			end := strings.LastIndex(line, ")")
			if start != -1 && end != -1 && start < end {
				return line[start+1 : end]
			}
		}
	}
	return ""
}

// =============================================================================
// Config Tests
// =============================================================================

func TestConfig_SetAndGet(t *testing.T) {
	home := setupTestHome(t)

	// Set a value
	stdout, _, err := runCmd(t, home, "config", "set", "api_key", "test-key-123")
	if err != nil {
		t.Fatalf("config set failed: %v", err)
	}
	if !strings.Contains(stdout, "test-key-123") {
		t.Errorf("expected confirmation message, got: %s", stdout)
	}

	// Get the value back
	stdout, _, err = runCmd(t, home, "config", "get", "api_key")
	if err != nil {
		t.Fatalf("config get failed: %v", err)
	}
	if !strings.Contains(stdout, "test-key-123") {
		t.Errorf("expected 'test-key-123', got: %s", stdout)
	}
}

func TestConfig_GetMissingKey(t *testing.T) {
	home := setupTestHome(t)

	_, stderr, err := runCmd(t, home, "config", "get", "nonexistent")
	if err == nil {
		t.Error("expected error for missing key")
	}
	if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "nicht gefunden") {
		t.Errorf("expected 'not found' error, got: %s", stderr)
	}
}

func TestConfig_SetLanguage(t *testing.T) {
	home := setupTestHome(t)

	// Set language to German
	_, _, err := runCmd(t, home, "config", "set", "language", "de")
	if err != nil {
		t.Fatalf("config set language failed: %v", err)
	}

	// Verify it persists
	stdout, _, err := runCmd(t, home, "config", "get", "language")
	if err != nil {
		t.Fatalf("config get language failed: %v", err)
	}
	if !strings.Contains(stdout, "de") {
		t.Errorf("expected 'de', got: %s", stdout)
	}
}

// =============================================================================
// Recipe Tests
// =============================================================================

func TestRecipe_ListEmpty(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if !strings.Contains(stdout, "No recipes") && !strings.Contains(stdout, "Keine Rezepte") {
		t.Errorf("expected empty list message, got: %s", stdout)
	}
}

func TestRecipe_ViewNotFound(t *testing.T) {
	home := setupTestHome(t)

	_, stderr, err := runCmd(t, home, "recipe", "view", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent recipe")
	}
	if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "nicht gefunden") {
		t.Errorf("expected 'not found' error, got: %s", stderr)
	}
}

func TestRecipe_DeleteNotFound(t *testing.T) {
	home := setupTestHome(t)

	_, stderr, err := runCmd(t, home, "recipe", "delete", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent recipe")
	}
	if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "nicht gefunden") {
		t.Errorf("expected 'not found' error, got: %s", stderr)
	}
}

// TestRecipe_Add tests adding a recipe interactively via stdin
func TestRecipe_Add(t *testing.T) {
	home := setupTestHome(t)

	// Input for recipe add:
	// Title, Description, Servings, PrepTime, CookTime,
	// Ingredients (empty line to finish), Steps (empty line to finish), Tags
	input := `My Test Recipe
A delicious test recipe
4
15
30
200 g flour
100 ml milk
2 eggs

Mix ingredients
Bake at 180C for 20 minutes

quick, easy
`

	stdout, _, err := runCmdWithInput(t, home, input, "recipe", "add")
	if err != nil {
		t.Fatalf("recipe add failed: %v", err)
	}
	if !strings.Contains(stdout, "My Test Recipe") {
		t.Errorf("expected recipe title in output, got: %s", stdout)
	}
	if !strings.Contains(stdout, "added") && !strings.Contains(stdout, "hinzugefügt") {
		t.Errorf("expected 'added' confirmation, got: %s", stdout)
	}

	// Verify recipe was saved by listing
	stdout, _, err = runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if !strings.Contains(stdout, "My Test Recipe") {
		t.Errorf("expected recipe in list, got: %s", stdout)
	}
}

// TestRecipe_AddMinimal tests adding a recipe with minimal info
func TestRecipe_AddMinimal(t *testing.T) {
	home := setupTestHome(t)

	// Minimal input - just title and one ingredient/step
	input := `Minimal Recipe



10
5
salt


Season food


`

	stdout, _, err := runCmdWithInput(t, home, input, "recipe", "add")
	if err != nil {
		t.Fatalf("recipe add failed: %v", err)
	}
	if !strings.Contains(stdout, "Minimal Recipe") {
		t.Errorf("expected recipe title in output, got: %s", stdout)
	}

	// Verify it's in the list
	stdout, _, err = runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if !strings.Contains(stdout, "Minimal Recipe") {
		t.Errorf("expected recipe in list, got: %s", stdout)
	}
}

// TestRecipe_AddThenView tests the full add->view flow
func TestRecipe_AddThenView(t *testing.T) {
	home := setupTestHome(t)

	input := `Pasta Carbonara
Classic Italian pasta dish
2
5
15
400 g spaghetti
150 g pancetta
3 egg yolks
50 g parmesan

Boil pasta until al dente
Fry pancetta until crispy
Mix egg yolks with parmesan
Combine all ingredients

italian, pasta
`

	stdout, _, err := runCmdWithInput(t, home, input, "recipe", "add")
	if err != nil {
		t.Fatalf("recipe add failed: %v", err)
	}

	recipeID := extractRecipeID(stdout)
	if recipeID == "" {
		t.Fatalf("could not extract recipe ID from output: %s", stdout)
	}

	// View the recipe
	stdout, _, err = runCmd(t, home, "recipe", "view", recipeID)
	if err != nil {
		t.Fatalf("recipe view failed: %v", err)
	}

	// Verify all fields are present
	if !strings.Contains(stdout, "Pasta Carbonara") {
		t.Errorf("expected title in view, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Classic Italian") {
		t.Errorf("expected description in view, got: %s", stdout)
	}
	if !strings.Contains(stdout, "spaghetti") {
		t.Errorf("expected ingredient in view, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Boil pasta") {
		t.Errorf("expected step in view, got: %s", stdout)
	}
	if !strings.Contains(stdout, "italian") {
		t.Errorf("expected tag in view, got: %s", stdout)
	}
}

// TestRecipe_AddThenDelete tests add->delete flow
func TestRecipe_AddThenDelete(t *testing.T) {
	home := setupTestHome(t)

	input := `Recipe To Delete
Will be deleted
2
5
10
ingredient

step 1


`

	stdout, _, err := runCmdWithInput(t, home, input, "recipe", "add")
	if err != nil {
		t.Fatalf("recipe add failed: %v", err)
	}

	recipeID := extractRecipeID(stdout)
	if recipeID == "" {
		t.Fatalf("could not extract recipe ID from output: %s", stdout)
	}

	// Delete it
	stdout, _, err = runCmd(t, home, "recipe", "delete", recipeID)
	if err != nil {
		t.Fatalf("recipe delete failed: %v", err)
	}
	if !strings.Contains(stdout, "deleted") && !strings.Contains(stdout, "gelöscht") {
		t.Errorf("expected delete confirmation, got: %s", stdout)
	}

	// Verify it's gone
	stdout, _, err = runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if strings.Contains(stdout, "Recipe To Delete") {
		t.Errorf("recipe should not be in list after delete, got: %s", stdout)
	}
}

// TestRecipe_AddWithIngredientParsing tests that ingredients are parsed correctly
func TestRecipe_AddWithIngredientParsing(t *testing.T) {
	home := setupTestHome(t)

	input := `Ingredient Test
Testing ingredient parsing
4
0
0
500 g flour
2 cups sugar
1 tbsp vanilla extract
3 eggs
pinch of salt

Mix well


`

	stdout, _, err := runCmdWithInput(t, home, input, "recipe", "add")
	if err != nil {
		t.Fatalf("recipe add failed: %v", err)
	}

	recipeID := extractRecipeID(stdout)
	if recipeID == "" {
		t.Fatalf("could not extract recipe ID from output: %s", stdout)
	}

	// View to check ingredient parsing
	stdout, _, err = runCmd(t, home, "recipe", "view", recipeID)
	if err != nil {
		t.Fatalf("recipe view failed: %v", err)
	}

	// Check various ingredient formats are shown
	if !strings.Contains(stdout, "flour") {
		t.Errorf("expected 'flour' in ingredients, got: %s", stdout)
	}
	if !strings.Contains(stdout, "sugar") {
		t.Errorf("expected 'sugar' in ingredients, got: %s", stdout)
	}
	if !strings.Contains(stdout, "vanilla") {
		t.Errorf("expected 'vanilla' in ingredients, got: %s", stdout)
	}
	if !strings.Contains(stdout, "eggs") {
		t.Errorf("expected 'eggs' in ingredients, got: %s", stdout)
	}
}

// TestRecipe_FullCRUD tests creating, viewing, listing, and deleting a recipe
// by directly manipulating the storage (since `recipe add` is interactive)
func TestRecipe_FullCRUD(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe directly in storage (simulating what `recipe add` does)
	recipe := map[string]interface{}{
		"id":          "test-123",
		"title":       "Test Pasta",
		"description": "A test recipe",
		"servings":    4,
		"prep_time":   10,
		"cook_time":   20,
		"ingredients": []map[string]interface{}{
			{"name": "pasta", "quantity": 500, "unit": "g", "category": "pantry"},
			{"name": "tomatoes", "quantity": 400, "unit": "g", "category": "produce"},
		},
		"steps":    []string{"Boil pasta", "Add sauce"},
		"tags":     []string{"italian", "quick"},
		"source":   "manual",
		"language": "en",
	}

	// Write recipes.json
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write test recipe: %v", err)
	}

	// Test list
	stdout, _, err := runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if !strings.Contains(stdout, "Test Pasta") {
		t.Errorf("expected 'Test Pasta' in list, got: %s", stdout)
	}
	if !strings.Contains(stdout, "test-123") {
		t.Errorf("expected recipe ID in list, got: %s", stdout)
	}

	// Test view
	stdout, _, err = runCmd(t, home, "recipe", "view", "test-123")
	if err != nil {
		t.Fatalf("recipe view failed: %v", err)
	}
	if !strings.Contains(stdout, "Test Pasta") {
		t.Errorf("expected title in view, got: %s", stdout)
	}
	if !strings.Contains(stdout, "pasta") {
		t.Errorf("expected ingredients in view, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Boil pasta") {
		t.Errorf("expected steps in view, got: %s", stdout)
	}

	// Test delete
	stdout, _, err = runCmd(t, home, "recipe", "delete", "test-123")
	if err != nil {
		t.Fatalf("recipe delete failed: %v", err)
	}
	if !strings.Contains(stdout, "deleted") && !strings.Contains(stdout, "gelöscht") {
		t.Errorf("expected deletion confirmation, got: %s", stdout)
	}

	// Verify it's gone
	_, stderr, err := runCmd(t, home, "recipe", "view", "test-123")
	if err == nil {
		t.Error("expected error viewing deleted recipe")
	}
	if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "nicht gefunden") {
		t.Errorf("expected 'not found' after delete, got: %s", stderr)
	}
}

func TestRecipe_SaveWithoutGenerate(t *testing.T) {
	home := setupTestHome(t)

	_, stderr, err := runCmd(t, home, "recipe", "save")
	if err == nil {
		t.Error("expected error when saving without prior generate")
	}
	if !strings.Contains(stderr, "No generated") && !strings.Contains(stderr, "Kein generiertes") {
		t.Errorf("expected 'no generated recipe' error, got: %s", stderr)
	}
}

// =============================================================================
// Shopping List Tests
// =============================================================================

func TestShopping_ShowEmpty(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show failed: %v", err)
	}
	if !strings.Contains(stdout, "empty") && !strings.Contains(stdout, "leer") {
		t.Errorf("expected empty list message, got: %s", stdout)
	}
}

func TestShopping_ClearEmpty(t *testing.T) {
	home := setupTestHome(t)

	// Should succeed even when empty
	stdout, _, err := runCmd(t, home, "shopping", "clear")
	if err != nil {
		t.Fatalf("shopping clear failed: %v", err)
	}
	if !strings.Contains(stdout, "cleared") && !strings.Contains(stdout, "geleert") {
		t.Errorf("expected cleared message, got: %s", stdout)
	}
}

func TestShopping_GenerateFromRecipe(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe first
	recipe := map[string]interface{}{
		"id":          "pasta-1",
		"title":       "Simple Pasta",
		"description": "Easy pasta",
		"servings":    2,
		"prep_time":   5,
		"cook_time":   15,
		"ingredients": []map[string]interface{}{
			{"name": "spaghetti", "quantity": 250, "unit": "g", "category": "pantry"},
			{"name": "olive oil", "quantity": 2, "unit": "tbsp", "category": "pantry"},
			{"name": "garlic", "quantity": 3, "unit": "cloves", "category": "produce"},
		},
		"steps":    []string{"Cook pasta", "Add garlic and oil"},
		"tags":     []string{"quick"},
		"source":   "manual",
		"language": "en",
	}

	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write test recipe: %v", err)
	}

	// Generate shopping list
	stdout, _, err := runCmd(t, home, "shopping", "generate", "pasta-1")
	if err != nil {
		t.Fatalf("shopping generate failed: %v", err)
	}
	if !strings.Contains(stdout, "3") && !strings.Contains(stdout, "ingredients") {
		t.Errorf("expected success message with ingredient count, got: %s", stdout)
	}

	// Show shopping list
	stdout, _, err = runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show failed: %v", err)
	}
	if !strings.Contains(stdout, "spaghetti") {
		t.Errorf("expected 'spaghetti' in shopping list, got: %s", stdout)
	}
	if !strings.Contains(stdout, "garlic") {
		t.Errorf("expected 'garlic' in shopping list, got: %s", stdout)
	}

	// Clear shopping list
	stdout, _, err = runCmd(t, home, "shopping", "clear")
	if err != nil {
		t.Fatalf("shopping clear failed: %v", err)
	}

	// Verify it's empty
	stdout, _, err = runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show after clear failed: %v", err)
	}
	if !strings.Contains(stdout, "empty") && !strings.Contains(stdout, "leer") {
		t.Errorf("expected empty after clear, got: %s", stdout)
	}
}

func TestShopping_GenerateNonexistentRecipe(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "shopping", "generate", "nonexistent")
	// This should not fail entirely, but should show error for the recipe
	if err != nil {
		t.Fatalf("shopping generate failed entirely: %v", err)
	}
	if !strings.Contains(stdout, "not found") && !strings.Contains(stdout, "nicht gefunden") {
		t.Errorf("expected 'not found' message, got: %s", stdout)
	}
}

func TestShopping_GenerateMultipleRecipes(t *testing.T) {
	home := setupTestHome(t)

	// Create two recipes
	recipes := []map[string]interface{}{
		{
			"id":       "r1",
			"title":    "Recipe One",
			"servings": 2,
			"ingredients": []map[string]interface{}{
				{"name": "flour", "quantity": 200, "unit": "g", "category": "pantry"},
			},
			"steps":    []string{"Step 1"},
			"source":   "manual",
			"language": "en",
		},
		{
			"id":       "r2",
			"title":    "Recipe Two",
			"servings": 4,
			"ingredients": []map[string]interface{}{
				{"name": "sugar", "quantity": 100, "unit": "g", "category": "pantry"},
				{"name": "flour", "quantity": 300, "unit": "g", "category": "pantry"},
			},
			"steps":    []string{"Step 1"},
			"source":   "manual",
			"language": "en",
		},
	}

	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write test recipes: %v", err)
	}

	// Generate from both recipes
	_, _, err := runCmd(t, home, "shopping", "generate", "r1", "r2")
	if err != nil {
		t.Fatalf("shopping generate multiple failed: %v", err)
	}

	// Check shopping list has merged ingredients
	stdout, _, err := runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show failed: %v", err)
	}
	if !strings.Contains(stdout, "flour") {
		t.Errorf("expected 'flour' in list, got: %s", stdout)
	}
	if !strings.Contains(stdout, "sugar") {
		t.Errorf("expected 'sugar' in list, got: %s", stdout)
	}
	// Flour should be merged (200 + 300 = 500g)
	if !strings.Contains(stdout, "500") {
		t.Errorf("expected merged flour quantity (500g), got: %s", stdout)
	}
}

// =============================================================================
// Localization Tests
// =============================================================================

func TestLocalization_GermanMessages(t *testing.T) {
	home := setupTestHome(t)

	// Set language to German
	_, _, err := runCmd(t, home, "config", "set", "language", "de")
	if err != nil {
		t.Fatalf("config set language failed: %v", err)
	}

	// Check that messages are in German
	stdout, _, err := runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if !strings.Contains(stdout, "Keine Rezepte") {
		t.Errorf("expected German message 'Keine Rezepte', got: %s", stdout)
	}

	// Shopping list should also be in German
	stdout, _, err = runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show failed: %v", err)
	}
	if !strings.Contains(stdout, "leer") {
		t.Errorf("expected German 'leer', got: %s", stdout)
	}
}

// =============================================================================
// Help and Usage Tests
// =============================================================================

func TestHelp_RootCommand(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "--help")
	if err != nil {
		t.Fatalf("help failed: %v", err)
	}
	if !strings.Contains(stdout, "recipe") {
		t.Errorf("expected 'recipe' in help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "shopping") {
		t.Errorf("expected 'shopping' in help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "config") {
		t.Errorf("expected 'config' in help, got: %s", stdout)
	}
}

func TestHelp_RecipeCommand(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "recipe", "--help")
	if err != nil {
		t.Fatalf("recipe help failed: %v", err)
	}
	if !strings.Contains(stdout, "list") {
		t.Errorf("expected 'list' in recipe help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "view") {
		t.Errorf("expected 'view' in recipe help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "generate") {
		t.Errorf("expected 'generate' in recipe help, got: %s", stdout)
	}
}

func TestHelp_ShoppingCommand(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "shopping", "--help")
	if err != nil {
		t.Fatalf("shopping help failed: %v", err)
	}
	if !strings.Contains(stdout, "generate") {
		t.Errorf("expected 'generate' in shopping help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "show") {
		t.Errorf("expected 'show' in shopping help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "clear") {
		t.Errorf("expected 'clear' in shopping help, got: %s", stdout)
	}
}

// =============================================================================
// Interactive Mode (REPL) Tests
// =============================================================================

// interactiveSession manages a PTY session for testing interactive mode
type interactiveSession struct {
	t      *testing.T
	cmd    *exec.Cmd
	ptmx   *os.File
	output bytes.Buffer
}

// startInteractive starts the binary in interactive mode with a PTY
func startInteractive(t *testing.T, homeDir string) *interactiveSession {
	t.Helper()

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(),
		"RECIPE_BOX_HOME="+homeDir,
		"NO_COLOR=1",
		"TERM=xterm-256color",
	)

	// Start with PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("failed to start PTY: %v", err)
	}

	// Set terminal size
	pty.Setsize(ptmx, &pty.Winsize{Rows: 24, Cols: 80})

	s := &interactiveSession{
		t:    t,
		cmd:  cmd,
		ptmx: ptmx,
	}

	// Start reading output in background
	go s.readLoop()

	// Wait for initial prompt
	time.Sleep(500 * time.Millisecond)

	return s
}

func (s *interactiveSession) readLoop() {
	buf := make([]byte, 1024)
	for {
		n, err := s.ptmx.Read(buf)
		if n > 0 {
			s.output.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func (s *interactiveSession) send(input string) {
	// Type each character with small delay to let go-prompt process
	for _, c := range input {
		s.ptmx.Write([]byte(string(c)))
		time.Sleep(10 * time.Millisecond)
	}
	// Send Enter key
	time.Sleep(50 * time.Millisecond)
	s.ptmx.Write([]byte("\n"))
}

func (s *interactiveSession) waitAndRead(d time.Duration) string {
	time.Sleep(d)
	out := s.output.String()
	s.output.Reset()
	return out
}

func (s *interactiveSession) getFullOutput() string {
	return s.output.String()
}

func (s *interactiveSession) close() {
	s.send("exit")
	time.Sleep(200 * time.Millisecond)
	s.ptmx.Close()
	s.cmd.Process.Kill()
	s.cmd.Wait()
}

func TestInteractive_StartsAndShowsPrompt(t *testing.T) {
	home := setupTestHome(t)
	session := startInteractive(t, home)
	defer session.close()

	// Check initial output has welcome message
	output := session.getFullOutput()
	if !strings.Contains(output, "Recipe Box") {
		t.Errorf("expected 'Recipe Box' in welcome, got: %s", output)
	}

	// Send help command
	session.send("help")
	output = session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "recipe") {
		t.Errorf("expected 'recipe' in help output, got: %s", output)
	}
}

func TestInteractive_RecipeList(t *testing.T) {
	home := setupTestHome(t)
	session := startInteractive(t, home)
	defer session.close()

	// Run recipe list
	session.send("recipe list")
	output := session.waitAndRead(500 * time.Millisecond)

	// Should show empty list message
	if !strings.Contains(output, "No recipes") && !strings.Contains(output, "Keine Rezepte") {
		t.Errorf("expected empty list message, got: %s", output)
	}
}

func TestInteractive_MultipleCommands(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe first
	recipe := map[string]interface{}{
		"id":       "int-test-1",
		"title":    "Interactive Test Recipe",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "test ingredient", "quantity": 1, "unit": "cup", "category": "pantry"},
		},
		"steps":    []string{"Test step"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	session := startInteractive(t, home)
	defer session.close()

	// List recipes
	session.send("recipe list")
	output := session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "Interactive Test Recipe") {
		t.Errorf("expected recipe in list, got: %s", output)
	}

	// View recipe
	session.send("recipe view int-test-1")
	output = session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "Interactive Test Recipe") {
		t.Errorf("expected recipe title in view, got: %s", output)
	}

	// Generate shopping list
	session.send("shopping generate int-test-1")
	output = session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "Added") && !strings.Contains(output, "hinzugefügt") {
		t.Errorf("expected shopping add confirmation, got: %s", output)
	}

	// Show shopping list
	session.send("shopping show")
	output = session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "test ingredient") {
		t.Errorf("expected ingredient in shopping list, got: %s", output)
	}
}

func TestInteractive_ExitCommand(t *testing.T) {
	home := setupTestHome(t)

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "RECIPE_BOX_HOME="+home, "NO_COLOR=1", "TERM=xterm-256color")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("failed to start PTY: %v", err)
	}
	defer func() {
		ptmx.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	pty.Setsize(ptmx, &pty.Winsize{Rows: 24, Cols: 80})

	// Read output in background
	var output bytes.Buffer
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	// Wait for startup
	time.Sleep(500 * time.Millisecond)

	// Type exit command character by character
	for _, c := range "exit" {
		ptmx.Write([]byte(string(c)))
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	ptmx.Write([]byte("\n"))

	// Wait and check output for "Goodbye!"
	time.Sleep(500 * time.Millisecond)
	if !strings.Contains(output.String(), "Goodbye") {
		t.Errorf("expected 'Goodbye' in output, got: %s", output.String())
	}
}

func TestInteractive_QuitCommand(t *testing.T) {
	home := setupTestHome(t)

	cmd := exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), "RECIPE_BOX_HOME="+home, "NO_COLOR=1", "TERM=xterm-256color")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("failed to start PTY: %v", err)
	}
	defer func() {
		ptmx.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	pty.Setsize(ptmx, &pty.Winsize{Rows: 24, Cols: 80})

	// Read output in background
	var output bytes.Buffer
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	// Wait for startup
	time.Sleep(500 * time.Millisecond)

	// Type quit command character by character
	for _, c := range "quit" {
		ptmx.Write([]byte(string(c)))
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	ptmx.Write([]byte("\n"))

	// Wait and check output for "Goodbye!"
	time.Sleep(500 * time.Millisecond)
	if !strings.Contains(output.String(), "Goodbye") {
		t.Errorf("expected 'Goodbye' in output, got: %s", output.String())
	}
}

func TestInteractive_InvalidCommand(t *testing.T) {
	home := setupTestHome(t)
	session := startInteractive(t, home)
	defer session.close()

	// Send invalid command
	session.send("notacommand")
	output := session.waitAndRead(500 * time.Millisecond)

	// Should show error (go-prompt echoes command, Cobra shows error)
	if !strings.Contains(output, "unknown") && !strings.Contains(output, "Error") && !strings.Contains(output, "notacommand") {
		t.Errorf("expected error or echo for unknown command, got: %s", output)
	}
}

func TestInteractive_ConfigPersists(t *testing.T) {
	home := setupTestHome(t)
	session := startInteractive(t, home)
	defer session.close()

	// Set config
	session.send("config set language de")
	output := session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "de") {
		t.Errorf("expected config confirmation, got: %s", output)
	}

	// Verify with get
	session.send("config get language")
	output = session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "de") {
		t.Errorf("expected 'de' in config get, got: %s", output)
	}
}

// =============================================================================
// Mock API Server for Testing
// =============================================================================

// mockAPIResponse returns a Claude API response with a recipe JSON
func mockAPIResponse(recipeJSON string) string {
	return fmt.Sprintf(`{
		"content": [
			{
				"type": "text",
				"text": %s
			}
		]
	}`, recipeJSON)
}

// startMockAPIServer starts a mock server that returns canned recipe responses
func startMockAPIServer(t *testing.T) *httptest.Server {
	t.Helper()

	// Sample recipe response that Claude would return
	recipeJSON := `"{\"title\":\"Mock Pasta Carbonara\",\"description\":\"A creamy Italian classic\",\"servings\":4,\"prep_time\":10,\"cook_time\":20,\"ingredients\":[{\"name\":\"spaghetti\",\"quantity\":400,\"unit\":\"g\",\"category\":\"pantry\"},{\"name\":\"pancetta\",\"quantity\":200,\"unit\":\"g\",\"category\":\"meat\"},{\"name\":\"eggs\",\"quantity\":3,\"unit\":\"pcs\",\"category\":\"dairy\"},{\"name\":\"parmesan\",\"quantity\":100,\"unit\":\"g\",\"category\":\"dairy\"}],\"steps\":[\"Boil pasta until al dente\",\"Fry pancetta until crispy\",\"Mix eggs with parmesan\",\"Combine all ingredients\"],\"tags\":[\"italian\",\"pasta\",\"quick\"],\"ascii_art\":\"  ____\\n /    \\\\\\n| PASTA |\\n \\\\____/\"}"`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify it's a POST request with correct headers
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "missing API key"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockAPIResponse(recipeJSON)))
	}))

	return server
}

// runCmdWithEnv executes the CLI with additional environment variables
func runCmdWithEnv(t *testing.T, homeDir string, env map[string]string, args ...string) (string, string, error) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)

	// Start with base env
	cmdEnv := append(os.Environ(), "RECIPE_BOX_HOME="+homeDir, "NO_COLOR=1")

	// Add custom env vars
	for k, v := range env {
		cmdEnv = append(cmdEnv, k+"="+v)
	}
	cmd.Env = cmdEnv

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// =============================================================================
// Recipe Generate Tests (with Mock API)
// =============================================================================

func TestRecipe_GenerateWithMockAPI(t *testing.T) {
	home := setupTestHome(t)

	// Start mock API server
	server := startMockAPIServer(t)
	defer server.Close()

	// Set API key in config
	_, _, err := runCmd(t, home, "config", "set", "api_key", "test-key-123")
	if err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	// Run generate with mock API URL
	env := map[string]string{"ANTHROPIC_API_URL": server.URL}
	stdout, stderr, err := runCmdWithEnv(t, home, env, "recipe", "generate", "--prompt", "pasta dish")
	if err != nil {
		t.Fatalf("recipe generate failed: %v\nstderr: %s", err, stderr)
	}

	// Verify recipe was generated
	if !strings.Contains(stdout, "Mock Pasta Carbonara") {
		t.Errorf("expected recipe title in output, got: %s", stdout)
	}
	if !strings.Contains(stdout, "spaghetti") {
		t.Errorf("expected ingredient in output, got: %s", stdout)
	}
}

func TestRecipe_GenerateAndSaveWithMockAPI(t *testing.T) {
	home := setupTestHome(t)

	// Start mock API server
	server := startMockAPIServer(t)
	defer server.Close()

	// Set API key
	_, _, err := runCmd(t, home, "config", "set", "api_key", "test-key-123")
	if err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	// Generate recipe
	env := map[string]string{"ANTHROPIC_API_URL": server.URL}
	_, stderr, err := runCmdWithEnv(t, home, env, "recipe", "generate", "--prompt", "pasta")
	if err != nil {
		t.Fatalf("recipe generate failed: %v\nstderr: %s", err, stderr)
	}

	// Save the generated recipe
	stdout, stderr, err := runCmd(t, home, "recipe", "save")
	if err != nil {
		t.Fatalf("recipe save failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "saved") && !strings.Contains(stdout, "gespeichert") {
		t.Errorf("expected save confirmation, got: %s", stdout)
	}

	// Verify it appears in list
	stdout, _, err = runCmd(t, home, "recipe", "list")
	if err != nil {
		t.Fatalf("recipe list failed: %v", err)
	}
	if !strings.Contains(stdout, "Mock Pasta Carbonara") {
		t.Errorf("expected saved recipe in list, got: %s", stdout)
	}
}

func TestRecipe_GenerateWithoutAPIKey(t *testing.T) {
	home := setupTestHome(t)

	// Try to generate without setting API key
	_, stderr, err := runCmd(t, home, "recipe", "generate", "--prompt", "pasta")
	if err == nil {
		t.Error("expected error when generating without API key")
	}
	if !strings.Contains(stderr, "API key") && !strings.Contains(stderr, "api_key") {
		t.Errorf("expected API key error, got: %s", stderr)
	}
}

func TestRecipe_GenerateWithFlags(t *testing.T) {
	home := setupTestHome(t)

	// Start mock API server
	server := startMockAPIServer(t)
	defer server.Close()

	// Set API key
	_, _, err := runCmd(t, home, "config", "set", "api_key", "test-key-123")
	if err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	// Generate with various flags
	env := map[string]string{"ANTHROPIC_API_URL": server.URL}
	stdout, stderr, err := runCmdWithEnv(t, home, env, "recipe", "generate",
		"--prompt", "italian dinner",
		"--cuisine", "italian",
		"--vegetarian",
		"--quick")
	if err != nil {
		t.Fatalf("recipe generate with flags failed: %v\nstderr: %s", err, stderr)
	}

	// Verify recipe was generated (mock returns same recipe regardless of flags)
	if !strings.Contains(stdout, "Mock Pasta Carbonara") {
		t.Errorf("expected recipe in output, got: %s", stdout)
	}
}

// =============================================================================
// Servings Scaling Tests
// =============================================================================

func TestRecipe_ViewWithServingsScaling(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe with known quantities
	recipe := map[string]interface{}{
		"id":          "scale-test",
		"title":       "Scalable Recipe",
		"description": "Test recipe for scaling",
		"servings":    2,
		"prep_time":   10,
		"cook_time":   20,
		"ingredients": []map[string]interface{}{
			{"name": "flour", "quantity": 200, "unit": "g", "category": "pantry"},
			{"name": "sugar", "quantity": 100, "unit": "g", "category": "pantry"},
			{"name": "eggs", "quantity": 2, "unit": "pcs", "category": "dairy"},
		},
		"steps":    []string{"Mix ingredients", "Bake"},
		"tags":     []string{"baking"},
		"source":   "manual",
		"language": "en",
	}

	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write test recipe: %v", err)
	}

	// View with scaled servings (2 -> 4 = 2x)
	stdout, _, err := runCmd(t, home, "recipe", "view", "scale-test", "--servings", "4")
	if err != nil {
		t.Fatalf("recipe view with servings failed: %v", err)
	}

	// Check scaling indicator is shown
	if !strings.Contains(stdout, "4") {
		t.Errorf("expected scaled servings (4) in output, got: %s", stdout)
	}

	// Flour should be 400g (200 * 2)
	if !strings.Contains(stdout, "400") {
		t.Errorf("expected scaled flour (400g) in output, got: %s", stdout)
	}

	// Sugar should be 200g (100 * 2)
	if !strings.Contains(stdout, "200") && !strings.Contains(stdout, "sugar") {
		t.Errorf("expected scaled sugar in output, got: %s", stdout)
	}

	// Eggs should be 4 (2 * 2)
	if !strings.Contains(stdout, "4") && !strings.Contains(stdout, "eggs") {
		t.Errorf("expected scaled eggs in output, got: %s", stdout)
	}
}

func TestShopping_GenerateWithServingsScaling(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "shop-scale-test",
		"title":    "Shopping Scale Test",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "pasta", "quantity": 250, "unit": "g", "category": "pantry"},
			{"name": "tomatoes", "quantity": 400, "unit": "g", "category": "produce"},
		},
		"steps":    []string{"Cook"},
		"source":   "manual",
		"language": "en",
	}

	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write test recipe: %v", err)
	}

	// Generate shopping list with scaled servings (2 -> 6 = 3x)
	stdout, _, err := runCmd(t, home, "shopping", "generate", "shop-scale-test", "--servings", "6")
	if err != nil {
		t.Fatalf("shopping generate with servings failed: %v", err)
	}

	// Should show success
	if !strings.Contains(stdout, "Added") && !strings.Contains(stdout, "hinzugefügt") {
		t.Errorf("expected success message, got: %s", stdout)
	}

	// Show shopping list
	stdout, _, err = runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show failed: %v", err)
	}

	// Pasta should be 750g (250 * 3)
	if !strings.Contains(stdout, "750") {
		t.Errorf("expected scaled pasta (750g) in output, got: %s", stdout)
	}

	// Tomatoes should be 1200g (400 * 3)
	if !strings.Contains(stdout, "1200") {
		t.Errorf("expected scaled tomatoes (1200g) in output, got: %s", stdout)
	}
}

// =============================================================================
// Interactive Mode: Recipe Add Test
// =============================================================================

func TestInteractive_RecipeAdd(t *testing.T) {
	home := setupTestHome(t)
	session := startInteractive(t, home)
	defer session.close()

	// Start recipe add
	session.send("recipe add")
	time.Sleep(300 * time.Millisecond)

	// Enter recipe details one by one
	// Title
	session.send("Interactive Test Cake")
	time.Sleep(100 * time.Millisecond)

	// Description
	session.send("A cake made in interactive mode")
	time.Sleep(100 * time.Millisecond)

	// Servings
	session.send("8")
	time.Sleep(100 * time.Millisecond)

	// Prep time
	session.send("15")
	time.Sleep(100 * time.Millisecond)

	// Cook time
	session.send("45")
	time.Sleep(100 * time.Millisecond)

	// Ingredients
	session.send("300 g flour")
	time.Sleep(100 * time.Millisecond)
	session.send("200 g sugar")
	time.Sleep(100 * time.Millisecond)
	session.send("4 eggs")
	time.Sleep(100 * time.Millisecond)
	session.send("") // Empty line to finish ingredients
	time.Sleep(100 * time.Millisecond)

	// Steps
	session.send("Mix dry ingredients")
	time.Sleep(100 * time.Millisecond)
	session.send("Add eggs")
	time.Sleep(100 * time.Millisecond)
	session.send("Bake at 180C")
	time.Sleep(100 * time.Millisecond)
	session.send("") // Empty line to finish steps
	time.Sleep(100 * time.Millisecond)

	// Tags
	session.send("dessert, baking")
	output := session.waitAndRead(1 * time.Second)

	// Verify success message
	if !strings.Contains(output, "added") && !strings.Contains(output, "hinzugefügt") {
		t.Errorf("expected 'added' confirmation, got: %s", output)
	}

	// Verify recipe is in list
	session.send("recipe list")
	output = session.waitAndRead(500 * time.Millisecond)

	if !strings.Contains(output, "Interactive Test Cake") {
		t.Errorf("expected recipe in list, got: %s", output)
	}
}

// =============================================================================
// Meal Plan Tests
// =============================================================================

func TestPlan_CreateAndShow(t *testing.T) {
	home := setupTestHome(t)

	// Create a plan
	stdout, _, err := runCmd(t, home, "plan", "create", "--days", "5")
	if err != nil {
		t.Fatalf("plan create failed: %v", err)
	}
	if !strings.Contains(stdout, "5") {
		t.Errorf("expected '5 days' in output, got: %s", stdout)
	}

	// Show the plan
	stdout, _, err = runCmd(t, home, "plan", "show")
	if err != nil {
		t.Fatalf("plan show failed: %v", err)
	}
	// Should show 5 days with day names
	if !strings.Contains(stdout, "day") && !strings.Contains(stdout, "Day") && !strings.Contains(stdout, "Tage") {
		t.Errorf("expected days in output, got: %s", stdout)
	}
}

func TestPlan_ShowNoPlan(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "plan", "show")
	if err != nil {
		t.Fatalf("plan show failed: %v", err)
	}
	if !strings.Contains(stdout, "No meal plan") && !strings.Contains(stdout, "Kein Essensplan") {
		t.Errorf("expected 'no plan' message, got: %s", stdout)
	}
}

func TestPlan_Clear(t *testing.T) {
	home := setupTestHome(t)

	// Create a plan first
	_, _, err := runCmd(t, home, "plan", "create")
	if err != nil {
		t.Fatalf("plan create failed: %v", err)
	}

	// Clear the plan
	stdout, _, err := runCmd(t, home, "plan", "clear")
	if err != nil {
		t.Fatalf("plan clear failed: %v", err)
	}
	if !strings.Contains(stdout, "cleared") && !strings.Contains(stdout, "gelöscht") {
		t.Errorf("expected 'cleared' message, got: %s", stdout)
	}

	// Verify it's gone
	stdout, _, err = runCmd(t, home, "plan", "show")
	if err != nil {
		t.Fatalf("plan show failed: %v", err)
	}
	if !strings.Contains(stdout, "No meal plan") && !strings.Contains(stdout, "Kein Essensplan") {
		t.Errorf("expected 'no plan' message after clear, got: %s", stdout)
	}
}

func TestPlan_AddRecipe(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe first
	recipe := map[string]interface{}{
		"id":       "plan-recipe-1",
		"title":    "Test Meal",
		"servings": 4,
		"ingredients": []map[string]interface{}{
			{"name": "chicken", "quantity": 500, "unit": "g", "category": "meat"},
		},
		"steps":    []string{"Cook chicken"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write recipe: %v", err)
	}

	// Create a plan
	_, _, err := runCmd(t, home, "plan", "create", "--days", "7")
	if err != nil {
		t.Fatalf("plan create failed: %v", err)
	}

	// Add recipe to monday
	stdout, _, err := runCmd(t, home, "plan", "add", "monday", "plan-recipe-1")
	if err != nil {
		t.Fatalf("plan add failed: %v", err)
	}
	if !strings.Contains(stdout, "Test Meal") && !strings.Contains(stdout, "Added") && !strings.Contains(stdout, "hinzugefügt") {
		t.Errorf("expected success message with recipe title, got: %s", stdout)
	}

	// Verify it shows in plan
	stdout, _, err = runCmd(t, home, "plan", "show")
	if err != nil {
		t.Fatalf("plan show failed: %v", err)
	}
	if !strings.Contains(stdout, "Test Meal") {
		t.Errorf("expected recipe in plan show, got: %s", stdout)
	}
}

func TestPlan_AddRecipeWithServings(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "plan-recipe-2",
		"title":    "Scalable Dish",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "pasta", "quantity": 200, "unit": "g", "category": "pantry"},
		},
		"steps":    []string{"Cook"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan
	runCmd(t, home, "plan", "create", "--days", "7")

	// Add recipe with custom servings
	stdout, _, err := runCmd(t, home, "plan", "add", "tuesday", "plan-recipe-2", "--servings", "6")
	if err != nil {
		t.Fatalf("plan add with servings failed: %v", err)
	}
	if !strings.Contains(stdout, "Scalable Dish") {
		t.Errorf("expected recipe name in output, got: %s", stdout)
	}

	// Verify servings show in plan
	stdout, _, _ = runCmd(t, home, "plan", "show")
	if !strings.Contains(stdout, "6") {
		t.Errorf("expected '6 servings' in plan show, got: %s", stdout)
	}
}

func TestPlan_AddRecipeWithDays(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "plan-recipe-3",
		"title":    "Leftover Meal",
		"servings": 6,
		"ingredients": []map[string]interface{}{
			{"name": "rice", "quantity": 500, "unit": "g", "category": "pantry"},
		},
		"steps":    []string{"Cook rice"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan
	runCmd(t, home, "plan", "create", "--days", "7")

	// Add recipe that covers 3 days
	stdout, _, err := runCmd(t, home, "plan", "add", "wednesday", "plan-recipe-3", "--days", "3")
	if err != nil {
		t.Fatalf("plan add with days failed: %v", err)
	}
	if !strings.Contains(stdout, "Leftover Meal") {
		t.Errorf("expected recipe name in output, got: %s", stdout)
	}

	// Verify leftover indicators show in plan
	stdout, _, _ = runCmd(t, home, "plan", "show")
	// Should show "covers 3 days" or leftover indicators
	if !strings.Contains(stdout, "3") && !strings.Contains(stdout, "leftover") && !strings.Contains(stdout, "Reste") {
		t.Errorf("expected leftover/covers days info in plan show, got: %s", stdout)
	}
}

func TestPlan_RemoveRecipe(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "plan-recipe-4",
		"title":    "To Be Removed",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "salad", "quantity": 200, "unit": "g", "category": "produce"},
		},
		"steps":    []string{"Mix"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan and add recipe
	runCmd(t, home, "plan", "create", "--days", "7")
	runCmd(t, home, "plan", "add", "friday", "plan-recipe-4")

	// Verify it's there
	stdout, _, _ := runCmd(t, home, "plan", "show")
	if !strings.Contains(stdout, "To Be Removed") {
		t.Fatalf("recipe should be in plan before removal, got: %s", stdout)
	}

	// Remove it
	stdout, _, err := runCmd(t, home, "plan", "remove", "friday")
	if err != nil {
		t.Fatalf("plan remove failed: %v", err)
	}
	if !strings.Contains(stdout, "Removed") && !strings.Contains(stdout, "entfernt") {
		t.Errorf("expected 'removed' confirmation, got: %s", stdout)
	}

	// Verify it's gone
	stdout, _, _ = runCmd(t, home, "plan", "show")
	if strings.Contains(stdout, "To Be Removed") {
		t.Errorf("recipe should not be in plan after removal, got: %s", stdout)
	}
}

func TestPlan_AddWithGermanDayName(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "german-plan-recipe",
		"title":    "German Day Test",
		"servings": 4,
		"ingredients": []map[string]interface{}{
			{"name": "wurst", "quantity": 200, "unit": "g", "category": "meat"},
		},
		"steps":    []string{"Grill"},
		"source":   "manual",
		"language": "de",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan
	runCmd(t, home, "plan", "create", "--days", "7")

	// Add recipe using German day name
	stdout, _, err := runCmd(t, home, "plan", "add", "montag", "german-plan-recipe")
	if err != nil {
		t.Fatalf("plan add with German day name failed: %v", err)
	}
	if !strings.Contains(stdout, "German Day Test") {
		t.Errorf("expected recipe name in output, got: %s", stdout)
	}

	// Verify it shows in plan
	stdout, _, _ = runCmd(t, home, "plan", "show")
	if !strings.Contains(stdout, "German Day Test") {
		t.Errorf("expected recipe in plan after adding with German day name, got: %s", stdout)
	}
}

func TestPlan_AddWithShortDayName(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "short-day-recipe",
		"title":    "Short Day Test",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "bread", "quantity": 2, "unit": "slices", "category": "pantry"},
		},
		"steps":    []string{"Toast"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan
	runCmd(t, home, "plan", "create", "--days", "7")

	// Add recipe using short day name
	stdout, _, err := runCmd(t, home, "plan", "add", "tue", "short-day-recipe")
	if err != nil {
		t.Fatalf("plan add with short day name failed: %v", err)
	}
	if !strings.Contains(stdout, "Short Day Test") {
		t.Errorf("expected recipe name in output, got: %s", stdout)
	}
}

func TestPlan_AddWithISODate(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "iso-date-recipe",
		"title":    "ISO Date Test",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "beans", "quantity": 400, "unit": "g", "category": "pantry"},
		},
		"steps":    []string{"Heat"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan and get the start date from show output
	runCmd(t, home, "plan", "create", "--days", "7")
	stdout, _, _ := runCmd(t, home, "plan", "show")

	// Extract a date from the output (format: YYYY-MM-DD)
	// The plan show displays dates like "Monday (2024-01-15)"
	dateStart := strings.Index(stdout, "20") // Look for year starting with 20
	if dateStart == -1 {
		t.Skip("Could not extract date from plan show output")
	}
	// Extract 10 characters for YYYY-MM-DD
	if dateStart+10 > len(stdout) {
		t.Skip("Could not extract full date from plan show output")
	}
	isoDate := stdout[dateStart : dateStart+10]

	// Add recipe using ISO date
	stdout, stderr, err := runCmd(t, home, "plan", "add", isoDate, "iso-date-recipe")
	if err != nil {
		t.Fatalf("plan add with ISO date failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "ISO Date Test") {
		t.Errorf("expected recipe name in output, got: %s", stdout)
	}
}

func TestPlan_AddNoPlan(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "no-plan-recipe",
		"title":    "No Plan Test",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "egg", "quantity": 1, "unit": "pcs", "category": "dairy"},
		},
		"steps":    []string{"Boil"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Try to add without creating a plan first
	stdout, _, err := runCmd(t, home, "plan", "add", "monday", "no-plan-recipe")
	if err != nil {
		// Error is expected
	}
	if !strings.Contains(stdout, "No meal plan") && !strings.Contains(stdout, "Kein Essensplan") {
		t.Errorf("expected 'no plan' message, got: %s", stdout)
	}
}

func TestPlan_AddRecipeNotFound(t *testing.T) {
	home := setupTestHome(t)

	// Create a plan but no recipe
	runCmd(t, home, "plan", "create", "--days", "7")

	// Try to add non-existent recipe
	_, stderr, err := runCmd(t, home, "plan", "add", "monday", "nonexistent-recipe")
	if err == nil {
		t.Error("expected error for non-existent recipe")
	}
	if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "nicht gefunden") {
		t.Errorf("expected 'not found' error, got: %s", stderr)
	}
}

func TestPlan_AddInvalidDay(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "invalid-day-recipe",
		"title":    "Invalid Day Test",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "tomato", "quantity": 1, "unit": "pcs", "category": "produce"},
		},
		"steps":    []string{"Slice"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan
	runCmd(t, home, "plan", "create", "--days", "7")

	// Try to add with invalid day name
	_, stderr, err := runCmd(t, home, "plan", "add", "invalidday", "invalid-day-recipe")
	if err == nil {
		t.Error("expected error for invalid day name")
	}
	if !strings.Contains(stderr, "not in") && !strings.Contains(stderr, "nicht im") {
		t.Errorf("expected 'not in plan' error, got: %s", stderr)
	}
}

func TestPlan_RemoveEmptyDay(t *testing.T) {
	home := setupTestHome(t)

	// Create a plan but don't add any recipes
	runCmd(t, home, "plan", "create", "--days", "7")

	// Try to remove from empty day (should not error, just report nothing to remove)
	stdout, _, err := runCmd(t, home, "plan", "remove", "monday")
	if err != nil {
		t.Fatalf("plan remove empty day failed: %v", err)
	}
	// Should indicate day is empty
	if !strings.Contains(stdout, "empty") && !strings.Contains(stdout, "leer") && !strings.Contains(stdout, "Monday") {
		t.Errorf("expected empty day indication, got: %s", stdout)
	}
}

func TestPlan_ReplaceExistingEntry(t *testing.T) {
	home := setupTestHome(t)

	// Create two recipes
	recipes := []map[string]interface{}{
		{
			"id":       "replace-recipe-1",
			"title":    "First Recipe",
			"servings": 2,
			"ingredients": []map[string]interface{}{
				{"name": "apple", "quantity": 1, "unit": "pcs", "category": "produce"},
			},
			"steps":    []string{"Eat"},
			"source":   "manual",
			"language": "en",
		},
		{
			"id":       "replace-recipe-2",
			"title":    "Second Recipe",
			"servings": 4,
			"ingredients": []map[string]interface{}{
				{"name": "banana", "quantity": 2, "unit": "pcs", "category": "produce"},
			},
			"steps":    []string{"Peel and eat"},
			"source":   "manual",
			"language": "en",
		},
	}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a plan and add first recipe
	runCmd(t, home, "plan", "create", "--days", "7")
	runCmd(t, home, "plan", "add", "thursday", "replace-recipe-1")

	// Verify first recipe is there
	stdout, _, _ := runCmd(t, home, "plan", "show")
	if !strings.Contains(stdout, "First Recipe") {
		t.Fatalf("first recipe should be in plan, got: %s", stdout)
	}

	// Add second recipe to same day (should replace)
	runCmd(t, home, "plan", "add", "thursday", "replace-recipe-2")

	// Verify second recipe replaced first
	stdout, _, _ = runCmd(t, home, "plan", "show")
	if strings.Contains(stdout, "First Recipe") {
		t.Errorf("first recipe should have been replaced, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Second Recipe") {
		t.Errorf("second recipe should be in plan, got: %s", stdout)
	}
}

func TestPlan_HelpCommands(t *testing.T) {
	home := setupTestHome(t)

	// Test plan help
	stdout, _, err := runCmd(t, home, "plan", "--help")
	if err != nil {
		t.Fatalf("plan help failed: %v", err)
	}
	if !strings.Contains(stdout, "add") {
		t.Errorf("expected 'add' in plan help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "remove") {
		t.Errorf("expected 'remove' in plan help, got: %s", stdout)
	}

	// Test plan add help
	stdout, _, err = runCmd(t, home, "plan", "add", "--help")
	if err != nil {
		t.Fatalf("plan add help failed: %v", err)
	}
	if !strings.Contains(stdout, "servings") {
		t.Errorf("expected 'servings' flag in help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "days") {
		t.Errorf("expected 'days' flag in help, got: %s", stdout)
	}

	// Test plan remove help
	stdout, _, err = runCmd(t, home, "plan", "remove", "--help")
	if err != nil {
		t.Fatalf("plan remove help failed: %v", err)
	}
	if !strings.Contains(stdout, "day") || !strings.Contains(stdout, "remove") {
		t.Errorf("expected usage info in help, got: %s", stdout)
	}
}

func TestInteractive_PlanCommands(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "interactive-plan-recipe",
		"title":    "Interactive Plan Test",
		"servings": 4,
		"ingredients": []map[string]interface{}{
			{"name": "flour", "quantity": 250, "unit": "g", "category": "pantry"},
		},
		"steps":    []string{"Mix"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	session := startInteractive(t, home)
	defer session.close()

	// Create plan
	session.send("plan create --days 7")
	output := session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "7") && !strings.Contains(output, "created") && !strings.Contains(output, "erstellt") {
		t.Errorf("expected plan created message, got: %s", output)
	}

	// Add recipe
	session.send("plan add monday interactive-plan-recipe --servings 6")
	output = session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "Interactive Plan Test") {
		t.Errorf("expected recipe name in add output, got: %s", output)
	}

	// Show plan
	session.send("plan show")
	output = session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "Interactive Plan Test") {
		t.Errorf("expected recipe in plan show, got: %s", output)
	}

	// Remove recipe
	session.send("plan remove monday")
	output = session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "Removed") && !strings.Contains(output, "entfernt") {
		t.Errorf("expected remove confirmation, got: %s", output)
	}

	// Clear plan
	session.send("plan clear")
	output = session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "cleared") && !strings.Contains(output, "gelöscht") {
		t.Errorf("expected clear confirmation, got: %s", output)
	}
}

func TestInteractive_ShoppingFromPlan(t *testing.T) {
	home := setupTestHome(t)

	// Create a recipe
	recipe := map[string]interface{}{
		"id":       "interactive-shop-plan",
		"title":    "Interactive Shopping Test",
		"servings": 2,
		"ingredients": []map[string]interface{}{
			{"name": "test ingredient", "quantity": 100, "unit": "g", "category": "pantry"},
		},
		"steps":    []string{"Cook"},
		"source":   "manual",
		"language": "en",
	}
	recipes := []interface{}{recipe}
	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	session := startInteractive(t, home)
	defer session.close()

	// Create plan and add recipe
	session.send("plan create --days 7")
	session.waitAndRead(500 * time.Millisecond)

	session.send("plan add monday interactive-shop-plan --servings 4")
	session.waitAndRead(500 * time.Millisecond)

	// Clear shopping list first
	session.send("shopping clear")
	session.waitAndRead(500 * time.Millisecond)

	// Generate shopping list from plan
	session.send("shopping generate --plan")
	output := session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "Interactive Shopping Test") {
		t.Errorf("expected recipe name in generate output, got: %s", output)
	}
	if !strings.Contains(output, "Added") && !strings.Contains(output, "hinzugefügt") {
		t.Errorf("expected 'Added' message, got: %s", output)
	}

	// Show shopping list
	session.send("shopping show")
	output = session.waitAndRead(500 * time.Millisecond)
	if !strings.Contains(output, "test ingredient") {
		t.Errorf("expected 'test ingredient' in shopping list, got: %s", output)
	}
	// Check scaled quantity (100g * 2 = 200g for 4 servings from 2)
	if !strings.Contains(output, "200") {
		t.Errorf("expected scaled quantity (200g) in shopping list, got: %s", output)
	}
}

// =============================================================================
// Shopping from Meal Plan Tests
// =============================================================================

func TestShopping_GenerateFromPlan(t *testing.T) {
	home := setupTestHome(t)

	// Create two recipes
	recipes := []map[string]interface{}{
		{
			"id":       "plan-shop-1",
			"title":    "Pasta Dish",
			"servings": 2,
			"ingredients": []map[string]interface{}{
				{"name": "pasta", "quantity": 200, "unit": "g", "category": "pantry"},
				{"name": "tomatoes", "quantity": 400, "unit": "g", "category": "produce"},
			},
			"steps":    []string{"Cook pasta", "Add sauce"},
			"source":   "manual",
			"language": "en",
		},
		{
			"id":       "plan-shop-2",
			"title":    "Rice Bowl",
			"servings": 4,
			"ingredients": []map[string]interface{}{
				{"name": "rice", "quantity": 300, "unit": "g", "category": "pantry"},
				{"name": "chicken", "quantity": 500, "unit": "g", "category": "meat"},
			},
			"steps":    []string{"Cook rice", "Add chicken"},
			"source":   "manual",
			"language": "en",
		},
	}

	data, _ := json.MarshalIndent(recipes, "", "  ")
	if err := os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644); err != nil {
		t.Fatalf("failed to write test recipes: %v", err)
	}

	// Create a meal plan and add recipes
	_, _, err := runCmd(t, home, "plan", "create", "--days", "7")
	if err != nil {
		t.Fatalf("plan create failed: %v", err)
	}

	_, _, err = runCmd(t, home, "plan", "add", "monday", "plan-shop-1", "--servings", "4")
	if err != nil {
		t.Fatalf("plan add monday failed: %v", err)
	}

	_, _, err = runCmd(t, home, "plan", "add", "wednesday", "plan-shop-2")
	if err != nil {
		t.Fatalf("plan add wednesday failed: %v", err)
	}

	// Generate shopping list from plan
	stdout, _, err := runCmd(t, home, "shopping", "generate", "--plan")
	if err != nil {
		t.Fatalf("shopping generate --plan failed: %v", err)
	}

	// Should show success messages for both recipes
	if !strings.Contains(stdout, "Pasta Dish") {
		t.Errorf("expected 'Pasta Dish' in output, got: %s", stdout)
	}
	if !strings.Contains(stdout, "Rice Bowl") {
		t.Errorf("expected 'Rice Bowl' in output, got: %s", stdout)
	}

	// Show shopping list
	stdout, _, err = runCmd(t, home, "shopping", "show")
	if err != nil {
		t.Fatalf("shopping show failed: %v", err)
	}

	// Check that all ingredients are present
	if !strings.Contains(stdout, "pasta") {
		t.Errorf("expected 'pasta' in shopping list, got: %s", stdout)
	}
	if !strings.Contains(stdout, "tomatoes") {
		t.Errorf("expected 'tomatoes' in shopping list, got: %s", stdout)
	}
	if !strings.Contains(stdout, "rice") {
		t.Errorf("expected 'rice' in shopping list, got: %s", stdout)
	}
	if !strings.Contains(stdout, "chicken") {
		t.Errorf("expected 'chicken' in shopping list, got: %s", stdout)
	}

	// Check that pasta was scaled (200g * 2 = 400g for 4 servings from 2)
	if !strings.Contains(stdout, "400") {
		t.Errorf("expected scaled pasta (400g) in shopping list, got: %s", stdout)
	}
}

func TestShopping_GenerateFromPlanWithDays(t *testing.T) {
	home := setupTestHome(t)

	// Create two recipes
	recipes := []map[string]interface{}{
		{
			"id":       "days-shop-1",
			"title":    "Day One Meal",
			"servings": 2,
			"ingredients": []map[string]interface{}{
				{"name": "ingredient1", "quantity": 100, "unit": "g", "category": "pantry"},
			},
			"steps":    []string{"Cook"},
			"source":   "manual",
			"language": "en",
		},
		{
			"id":       "days-shop-2",
			"title":    "Day Six Meal",
			"servings": 2,
			"ingredients": []map[string]interface{}{
				{"name": "ingredient2", "quantity": 200, "unit": "g", "category": "pantry"},
			},
			"steps":    []string{"Cook"},
			"source":   "manual",
			"language": "en",
		},
	}

	data, _ := json.MarshalIndent(recipes, "", "  ")
	os.WriteFile(filepath.Join(home, "recipes.json"), data, 0644)

	// Create a meal plan
	runCmd(t, home, "plan", "create", "--days", "7")

	// Get the plan dates by reading the plan file
	today := time.Now().Format("2006-01-02")
	day6 := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Add recipes to day 1 (today) and day 6 using ISO dates
	runCmd(t, home, "plan", "add", today, "days-shop-1")
	runCmd(t, home, "plan", "add", day6, "days-shop-2")

	// Generate shopping list from only first 3 days
	stdout, _, err := runCmd(t, home, "shopping", "generate", "--plan", "--days", "3")
	if err != nil {
		t.Fatalf("shopping generate --plan --days failed: %v", err)
	}

	// Should only show Day One's recipe (today is day 1)
	if !strings.Contains(stdout, "Day One Meal") {
		t.Errorf("expected 'Day One Meal' in output, got: %s", stdout)
	}

	// Day Six's recipe should NOT be included
	if strings.Contains(stdout, "Day Six Meal") {
		t.Errorf("did not expect 'Day Six Meal' in output for --days 3, got: %s", stdout)
	}

	// Check shopping list
	stdout, _, _ = runCmd(t, home, "shopping", "show")
	if !strings.Contains(stdout, "ingredient1") {
		t.Errorf("expected 'ingredient1' in shopping list, got: %s", stdout)
	}
	if strings.Contains(stdout, "ingredient2") {
		t.Errorf("did not expect 'ingredient2' in shopping list for --days 3, got: %s", stdout)
	}
}

func TestShopping_GenerateFromPlanNoPlan(t *testing.T) {
	home := setupTestHome(t)

	// Try to generate from plan when no plan exists
	stdout, _, err := runCmd(t, home, "shopping", "generate", "--plan")
	if err != nil {
		t.Fatalf("shopping generate --plan should not error, got: %v", err)
	}

	// Should show "no plan" message
	if !strings.Contains(stdout, "No meal plan") && !strings.Contains(stdout, "Kein Essensplan") {
		t.Errorf("expected 'no plan' message, got: %s", stdout)
	}
}

func TestShopping_GenerateFromPlanEmptyPlan(t *testing.T) {
	home := setupTestHome(t)

	// Create an empty plan
	_, _, err := runCmd(t, home, "plan", "create", "--days", "7")
	if err != nil {
		t.Fatalf("plan create failed: %v", err)
	}

	// Try to generate from empty plan
	stdout, _, err := runCmd(t, home, "shopping", "generate", "--plan")
	if err != nil {
		t.Fatalf("shopping generate --plan should not error, got: %v", err)
	}

	// Should show "empty plan" message
	if !strings.Contains(stdout, "empty") && !strings.Contains(stdout, "leer") {
		t.Errorf("expected 'empty' message, got: %s", stdout)
	}
}

func TestShopping_GenerateNoArgsError(t *testing.T) {
	home := setupTestHome(t)

	// Try to generate without args and without --plan
	_, stderr, err := runCmd(t, home, "shopping", "generate")
	if err == nil {
		t.Errorf("expected error when no args and no --plan flag")
	}

	// Should show helpful error message
	combined := stderr
	if !strings.Contains(combined, "recipe ID") && !strings.Contains(combined, "--plan") && !strings.Contains(combined, "Rezept-ID") {
		t.Errorf("expected helpful error message about recipe ID or --plan, got: %s", combined)
	}
}

// =============================================================================
// Plan Generate Tests (AI Meal Planning)
// =============================================================================

// startMockMealPlanAPIServer starts a mock server for meal plan generation
func startMockMealPlanAPIServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("x-api-key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "missing API key"}`))
			return
		}

		// Parse request to check if it's for meal plan or recipe
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		messages := reqBody["messages"].([]interface{})
		firstMsg := messages[0].(map[string]interface{})
		content := firstMsg["content"].(string)

		var responseJSON string

		if strings.Contains(content, "meal plan") || strings.Contains(content, "balanced") {
			// Return meal plan suggestions
			mealPlanJSON := `"{\"suggestions\":[{\"day\":\"Thursday\",\"date\":\"2026-01-23\",\"meal_name\":\"Pasta Primavera\",\"description\":\"Fresh vegetables with pasta\",\"servings\":4,\"prep_time\":15,\"cook_time\":20},{\"day\":\"Friday\",\"date\":\"2026-01-24\",\"meal_name\":\"Grilled Salmon\",\"description\":\"Salmon with lemon and herbs\",\"servings\":4,\"prep_time\":10,\"cook_time\":15},{\"day\":\"Saturday\",\"date\":\"2026-01-25\",\"meal_name\":\"Vegetable Stir Fry\",\"description\":\"Asian-style vegetables\",\"servings\":4,\"prep_time\":15,\"cook_time\":10}]}"`
			responseJSON = fmt.Sprintf(`{"content":[{"type":"text","text":%s}]}`, mealPlanJSON)
		} else {
			// Return recipe for follow-up recipe generation
			recipeJSON := `"{\"title\":\"Generated Recipe\",\"description\":\"A delicious dish\",\"servings\":4,\"prep_time\":15,\"cook_time\":20,\"ingredients\":[{\"name\":\"ingredient\",\"quantity\":200,\"unit\":\"g\",\"category\":\"pantry\"}],\"steps\":[\"Step 1\",\"Step 2\"],\"tags\":[\"test\"],\"ascii_art\":\"  ___\\n /   \\\\\\n| O O |\\n \\\\___/\"}"`
			responseJSON = fmt.Sprintf(`{"content":[{"type":"text","text":%s}]}`, recipeJSON)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
}

func TestPlan_GenerateHelp(t *testing.T) {
	home := setupTestHome(t)

	stdout, _, err := runCmd(t, home, "plan", "generate", "--help")
	if err != nil {
		t.Fatalf("plan generate help failed: %v", err)
	}

	// Check for expected flags
	if !strings.Contains(stdout, "days") {
		t.Errorf("expected '--days' in help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "vegetarian") {
		t.Errorf("expected '--vegetarian' in help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "quick") {
		t.Errorf("expected '--quick' in help, got: %s", stdout)
	}
	if !strings.Contains(stdout, "cuisine") {
		t.Errorf("expected '--cuisine' in help, got: %s", stdout)
	}
}

func TestPlan_GenerateWithoutAPIKey(t *testing.T) {
	home := setupTestHome(t)

	// Try to generate without API key
	_, stderr, err := runCmd(t, home, "plan", "generate")
	if err == nil {
		t.Error("expected error when generating without API key")
	}
	if !strings.Contains(stderr, "API key") && !strings.Contains(stderr, "api_key") && !strings.Contains(stderr, "API-Schlüssel") {
		t.Errorf("expected API key error, got: %s", stderr)
	}
}

func TestPlan_GenerateWithMockAPI(t *testing.T) {
	home := setupTestHome(t)

	// Start mock API server
	server := startMockMealPlanAPIServer(t)
	defer server.Close()

	// Set API key
	_, _, err := runCmd(t, home, "config", "set", "api_key", "test-key-123")
	if err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	// Run plan generate with mock API and auto-answer "no" to prompts
	env := map[string]string{"ANTHROPIC_API_URL": server.URL}

	// Since interactive prompts are used, we need to provide input
	// For testing, we'll just check that the command starts and generates output
	cmd := exec.Command(binaryPath, "plan", "generate", "--days", "3")
	cmd.Env = append(os.Environ(),
		"RECIPE_BOX_HOME="+home,
		"NO_COLOR=1",
		"ANTHROPIC_API_URL="+server.URL,
	)

	// Provide "n" twice (decline approval and decline recipe creation)
	cmd.Stdin = strings.NewReader("n\n")

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		t.Fatalf("plan generate failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()

	// Should show generating message
	if !strings.Contains(output, "Generating") && !strings.Contains(output, "generiert") {
		t.Errorf("expected 'Generating' message, got: %s", output)
	}

	// Should show meal suggestions
	if !strings.Contains(output, "Pasta Primavera") && !strings.Contains(output, "meal") {
		t.Errorf("expected meal suggestions in output, got: %s", output)
	}

	// Verify env var was used properly
	_ = env
}

func TestPlan_GenerateWithFlagsNoInteraction(t *testing.T) {
	home := setupTestHome(t)

	// Start mock API server
	server := startMockMealPlanAPIServer(t)
	defer server.Close()

	// Set API key
	runCmd(t, home, "config", "set", "api_key", "test-key-123")

	// Run with various flags and decline prompts
	cmd := exec.Command(binaryPath, "plan", "generate",
		"--days", "3",
		"--vegetarian",
		"--quick",
		"--cuisine", "italian",
	)
	cmd.Env = append(os.Environ(),
		"RECIPE_BOX_HOME="+home,
		"NO_COLOR=1",
		"ANTHROPIC_API_URL="+server.URL,
	)
	cmd.Stdin = strings.NewReader("n\n")

	var stdout strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()
	if err != nil {
		t.Fatalf("plan generate with flags failed: %v\noutput: %s", err, stdout.String())
	}

	output := stdout.String()

	// Should show meal suggestions from mock
	if !strings.Contains(output, "Suggested") && !strings.Contains(output, "Vorgeschlagen") {
		t.Errorf("expected suggestion header in output, got: %s", output)
	}
}
