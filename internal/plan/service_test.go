package plan

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestService_Create(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	plan, err := svc.Create(7)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if plan.Days != 7 {
		t.Errorf("expected 7 days, got %d", plan.Days)
	}

	if plan.ID == "" {
		t.Error("expected non-empty ID")
	}

	if plan.StartDate == "" {
		t.Error("expected non-empty start date")
	}

	if len(plan.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(plan.Entries))
	}
}

func TestService_Create_DefaultDays(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	plan, err := svc.Create(0)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if plan.Days != 7 {
		t.Errorf("expected default 7 days, got %d", plan.Days)
	}
}

func TestService_Get(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan first
	created, err := svc.Create(5)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Get should return the plan
	plan, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if plan == nil {
		t.Fatal("expected plan, got nil")
	}

	if plan.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, plan.ID)
	}

	if plan.Days != 5 {
		t.Errorf("expected 5 days, got %d", plan.Days)
	}
}

func TestService_Get_NoPlan(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	plan, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if plan != nil {
		t.Errorf("expected nil, got plan with ID %s", plan.ID)
	}
}

func TestService_Exists(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	if svc.Exists() {
		t.Error("expected Exists() to return false before creating plan")
	}

	_, err := svc.Create(7)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if !svc.Exists() {
		t.Error("expected Exists() to return true after creating plan")
	}
}

func TestService_Clear(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan first
	_, err := svc.Create(7)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Clear should succeed
	if err := svc.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify file was deleted
	planPath := store.Path(mealPlanFile)
	if _, err := os.Stat(planPath); !os.IsNotExist(err) {
		t.Error("expected meal plan file to be deleted")
	}

	// Get should return nil
	plan, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if plan != nil {
		t.Errorf("expected nil after clear, got plan with ID %s", plan.ID)
	}
}

func TestService_Clear_NoPlan(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Clear on empty should not error
	if err := svc.Clear(); err != nil {
		t.Fatalf("Clear on empty failed: %v", err)
	}
}

func TestMealPlan_GetDates(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      3,
	}

	dates := plan.GetDates()

	if len(dates) != 3 {
		t.Fatalf("expected 3 dates, got %d", len(dates))
	}

	expected := []string{"2024-01-15", "2024-01-16", "2024-01-17"}
	for i, date := range dates {
		if date != expected[i] {
			t.Errorf("date[%d]: expected %s, got %s", i, expected[i], date)
		}
	}
}

func TestMealPlan_GetEntryForDate(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-16", RecipeID: "recipe1", Servings: 4, CoversDays: 1},
		},
	}

	// Should find entry
	entry := plan.GetEntryForDate("2024-01-16")
	if entry == nil {
		t.Fatal("expected to find entry")
	}
	if entry.RecipeID != "recipe1" {
		t.Errorf("expected recipe1, got %s", entry.RecipeID)
	}

	// Should return nil for date without entry
	entry = plan.GetEntryForDate("2024-01-15")
	if entry != nil {
		t.Error("expected nil for date without entry")
	}
}

func TestMealPlan_IsLeftoverDay(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-15", RecipeID: "recipe1", Servings: 4, CoversDays: 3},
		},
	}

	tests := []struct {
		date       string
		isLeftover bool
		recipeID   string
	}{
		{"2024-01-15", false, ""}, // main day, not leftover
		{"2024-01-16", true, "recipe1"},
		{"2024-01-17", true, "recipe1"},
		{"2024-01-18", false, ""}, // outside covers_days
	}

	for _, tc := range tests {
		isLeftover, entry := plan.IsLeftoverDay(tc.date)
		if isLeftover != tc.isLeftover {
			t.Errorf("IsLeftoverDay(%s): expected %v, got %v", tc.date, tc.isLeftover, isLeftover)
		}
		if tc.isLeftover && entry.RecipeID != tc.recipeID {
			t.Errorf("IsLeftoverDay(%s): expected recipe %s, got %s", tc.date, tc.recipeID, entry.RecipeID)
		}
	}
}

func TestNewMealPlan(t *testing.T) {
	plan := NewMealPlan(5)

	if plan.Days != 5 {
		t.Errorf("expected 5 days, got %d", plan.Days)
	}

	if plan.ID == "" {
		t.Error("expected non-empty ID")
	}

	// Start date should be today
	today := time.Now().Format("2006-01-02")
	if plan.StartDate != today {
		t.Errorf("expected start date %s, got %s", today, plan.StartDate)
	}

	if len(plan.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(plan.Entries))
	}
}

func TestService_Save(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	plan := &MealPlan{
		ID:        "test123",
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-15", RecipeID: "recipe1", Servings: 4, CoversDays: 2},
		},
	}

	if err := svc.Save(plan); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	planPath := filepath.Join(store.Path(mealPlanFile))
	if _, err := os.Stat(planPath); os.IsNotExist(err) {
		t.Error("expected meal plan file to exist")
	}

	// Load and verify
	loaded, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if loaded.ID != "test123" {
		t.Errorf("expected ID test123, got %s", loaded.ID)
	}

	if len(loaded.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(loaded.Entries))
	}
}
