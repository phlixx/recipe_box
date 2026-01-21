package shopping

import (
	"os"
	"path/filepath"
	"testing"

	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
)

func setupTestStorage(t *testing.T) *storage.Storage {
	t.Helper()

	// Create temp dir for testing
	tmpDir := t.TempDir()

	// Override home dir for storage
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", origHome)
	})

	store, err := storage.New()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	return store
}

func TestService_GenerateFromRecipe(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	r := &recipe.Recipe{
		ID:    "test123",
		Title: "Test Recipe",
		Ingredients: []recipe.Ingredient{
			{Name: "Flour", Quantity: 200, Unit: "g", Category: "pantry"},
			{Name: "Eggs", Quantity: 2, Unit: "", Category: "dairy"},
		},
	}

	err := svc.GenerateFromRecipe(r)
	if err != nil {
		t.Fatalf("GenerateFromRecipe failed: %v", err)
	}

	list, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(list.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(list.Items))
	}
}

func TestService_GenerateFromRecipe_MergesSameIngredients(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	r1 := &recipe.Recipe{
		ID:    "recipe1",
		Title: "Recipe 1",
		Ingredients: []recipe.Ingredient{
			{Name: "Flour", Quantity: 200, Unit: "g", Category: "pantry"},
		},
	}

	r2 := &recipe.Recipe{
		ID:    "recipe2",
		Title: "Recipe 2",
		Ingredients: []recipe.Ingredient{
			{Name: "Flour", Quantity: 100, Unit: "g", Category: "pantry"},
		},
	}

	if err := svc.GenerateFromRecipe(r1); err != nil {
		t.Fatalf("GenerateFromRecipe r1 failed: %v", err)
	}

	if err := svc.GenerateFromRecipe(r2); err != nil {
		t.Fatalf("GenerateFromRecipe r2 failed: %v", err)
	}

	list, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(list.Items) != 1 {
		t.Errorf("expected 1 merged item, got %d", len(list.Items))
	}

	if list.Items[0].Quantity != 300 {
		t.Errorf("expected merged quantity 300, got %g", list.Items[0].Quantity)
	}
}

func TestService_Get_EmptyList(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	list, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(list.Items) != 0 {
		t.Errorf("expected empty list, got %d items", len(list.Items))
	}
}

func TestService_Clear(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	r := &recipe.Recipe{
		ID:    "test123",
		Title: "Test Recipe",
		Ingredients: []recipe.Ingredient{
			{Name: "Flour", Quantity: 200, Unit: "g", Category: "pantry"},
		},
	}

	if err := svc.GenerateFromRecipe(r); err != nil {
		t.Fatalf("GenerateFromRecipe failed: %v", err)
	}

	if err := svc.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify file was deleted
	listPath := filepath.Join(store.Path(shoppingListFile))
	if _, err := os.Stat(listPath); !os.IsNotExist(err) {
		t.Error("expected shopping list file to be deleted")
	}

	// Get should return empty list
	list, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(list.Items) != 0 {
		t.Errorf("expected empty list after clear, got %d items", len(list.Items))
	}
}

func TestService_Clear_EmptyList(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Clear on empty list should not error
	if err := svc.Clear(); err != nil {
		t.Fatalf("Clear on empty list failed: %v", err)
	}
}

func TestList_GroupByCategory(t *testing.T) {
	list := &List{
		Items: []Item{
			{Name: "Flour", Category: "pantry"},
			{Name: "Sugar", Category: "pantry"},
			{Name: "Eggs", Category: "dairy"},
			{Name: "Mystery Item", Category: ""},
		},
	}

	groups := list.GroupByCategory()

	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}

	if len(groups["pantry"]) != 2 {
		t.Errorf("expected 2 pantry items, got %d", len(groups["pantry"]))
	}

	if len(groups["dairy"]) != 1 {
		t.Errorf("expected 1 dairy item, got %d", len(groups["dairy"]))
	}

	if len(groups["other"]) != 1 {
		t.Errorf("expected 1 other item, got %d", len(groups["other"]))
	}
}

func TestMergeItem_NewItem(t *testing.T) {
	items := []Item{}
	newItem := Item{Name: "Flour", Quantity: 200, Unit: "g"}

	result := mergeItem(items, newItem)

	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
}

func TestMergeItem_MergesQuantity(t *testing.T) {
	items := []Item{
		{Name: "Flour", Quantity: 200, Unit: "g"},
	}
	newItem := Item{Name: "Flour", Quantity: 100, Unit: "g"}

	result := mergeItem(items, newItem)

	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}

	if result[0].Quantity != 300 {
		t.Errorf("expected quantity 300, got %g", result[0].Quantity)
	}
}

func TestMergeItem_DifferentUnits(t *testing.T) {
	items := []Item{
		{Name: "Flour", Quantity: 200, Unit: "g"},
	}
	newItem := Item{Name: "Flour", Quantity: 1, Unit: "cup"}

	result := mergeItem(items, newItem)

	if len(result) != 2 {
		t.Errorf("expected 2 items for different units, got %d", len(result))
	}
}
