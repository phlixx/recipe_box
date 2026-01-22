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

func TestService_AddEntry(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan with a specific start date
	plan := &MealPlan{
		ID:        "test123",
		StartDate: "2024-01-15",
		Days:      7,
		Entries:   []PlanEntry{},
	}
	if err := svc.Save(plan); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Add an entry
	err := svc.AddEntry("2024-01-16", "recipe1", 4, 2)
	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	// Verify entry was added
	loaded, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}

	entry := loaded.Entries[0]
	if entry.Date != "2024-01-16" {
		t.Errorf("expected date 2024-01-16, got %s", entry.Date)
	}
	if entry.RecipeID != "recipe1" {
		t.Errorf("expected recipe1, got %s", entry.RecipeID)
	}
	if entry.Servings != 4 {
		t.Errorf("expected 4 servings, got %d", entry.Servings)
	}
	if entry.CoversDays != 2 {
		t.Errorf("expected 2 covers_days, got %d", entry.CoversDays)
	}
}

func TestService_AddEntry_ReplaceExisting(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan with an existing entry
	plan := &MealPlan{
		ID:        "test123",
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-16", RecipeID: "recipe1", Servings: 2, CoversDays: 1},
		},
	}
	if err := svc.Save(plan); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Replace the entry
	err := svc.AddEntry("2024-01-16", "recipe2", 6, 3)
	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	// Verify entry was replaced
	loaded, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry (replaced), got %d", len(loaded.Entries))
	}

	entry := loaded.Entries[0]
	if entry.RecipeID != "recipe2" {
		t.Errorf("expected recipe2 (replaced), got %s", entry.RecipeID)
	}
	if entry.Servings != 6 {
		t.Errorf("expected 6 servings (replaced), got %d", entry.Servings)
	}
}

func TestService_AddEntry_NoPlan(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	err := svc.AddEntry("2024-01-15", "recipe1", 4, 1)
	if err != ErrNoPlan {
		t.Errorf("expected ErrNoPlan, got %v", err)
	}
}

func TestService_AddEntry_DateNotInPlan(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan
	plan := &MealPlan{
		ID:        "test123",
		StartDate: "2024-01-15",
		Days:      3,
		Entries:   []PlanEntry{},
	}
	if err := svc.Save(plan); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Try to add entry outside plan range
	err := svc.AddEntry("2024-01-25", "recipe1", 4, 1)
	if err != ErrDateNotInPlan {
		t.Errorf("expected ErrDateNotInPlan, got %v", err)
	}
}

func TestService_RemoveEntry(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan with an entry
	plan := &MealPlan{
		ID:        "test123",
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-16", RecipeID: "recipe1", Servings: 4, CoversDays: 1},
		},
	}
	if err := svc.Save(plan); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Remove the entry
	removed, err := svc.RemoveEntry("2024-01-16")
	if err != nil {
		t.Fatalf("RemoveEntry failed: %v", err)
	}
	if !removed {
		t.Error("expected entry to be removed")
	}

	// Verify entry was removed
	loaded, err := svc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(loaded.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(loaded.Entries))
	}
}

func TestService_RemoveEntry_NotFound(t *testing.T) {
	store := setupTestStorage(t)
	svc := NewService(store)

	// Create a plan without the target entry
	plan := &MealPlan{
		ID:        "test123",
		StartDate: "2024-01-15",
		Days:      7,
		Entries:   []PlanEntry{},
	}
	if err := svc.Save(plan); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Try to remove non-existent entry
	removed, err := svc.RemoveEntry("2024-01-16")
	if err != nil {
		t.Fatalf("RemoveEntry failed: %v", err)
	}
	if removed {
		t.Error("expected entry not to be found")
	}
}

func TestMealPlan_ParseDayToDate(t *testing.T) {
	// Create a plan starting on Monday 2024-01-15
	plan := &MealPlan{
		StartDate: "2024-01-15", // Monday
		Days:      7,
	}

	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		// English day names
		{"monday", "2024-01-15", true},
		{"tuesday", "2024-01-16", true},
		{"wednesday", "2024-01-17", true},
		{"thursday", "2024-01-18", true},
		{"friday", "2024-01-19", true},
		{"saturday", "2024-01-20", true},
		{"sunday", "2024-01-21", true},
		// Short day names
		{"mon", "2024-01-15", true},
		{"tue", "2024-01-16", true},
		// German day names
		{"montag", "2024-01-15", true},
		{"dienstag", "2024-01-16", true},
		// ISO date
		{"2024-01-16", "2024-01-16", true},
		// Case insensitive
		{"MONDAY", "2024-01-15", true},
		{"Monday", "2024-01-15", true},
		// Invalid
		{"invalid", "", false},
		{"2024-01-25", "", false}, // outside plan range
	}

	for _, tc := range tests {
		date, ok := plan.ParseDayToDate(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseDayToDate(%q): expected ok=%v, got %v", tc.input, tc.ok, ok)
		}
		if date != tc.expected {
			t.Errorf("ParseDayToDate(%q): expected %q, got %q", tc.input, tc.expected, date)
		}
	}
}

func TestMealPlan_AddEntry(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      7,
		Entries:   []PlanEntry{},
	}

	entry := PlanEntry{
		Date:       "2024-01-16",
		RecipeID:   "recipe1",
		Servings:   4,
		CoversDays: 2,
	}

	plan.AddEntry(entry)

	if len(plan.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(plan.Entries))
	}

	if plan.Entries[0].RecipeID != "recipe1" {
		t.Errorf("expected recipe1, got %s", plan.Entries[0].RecipeID)
	}
}

func TestMealPlan_AddEntry_ReplaceExisting(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-16", RecipeID: "recipe1", Servings: 2, CoversDays: 1},
		},
	}

	entry := PlanEntry{
		Date:       "2024-01-16",
		RecipeID:   "recipe2",
		Servings:   6,
		CoversDays: 3,
	}

	plan.AddEntry(entry)

	if len(plan.Entries) != 1 {
		t.Fatalf("expected 1 entry (replaced), got %d", len(plan.Entries))
	}

	if plan.Entries[0].RecipeID != "recipe2" {
		t.Errorf("expected recipe2 (replaced), got %s", plan.Entries[0].RecipeID)
	}
}

func TestMealPlan_RemoveEntry(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      7,
		Entries: []PlanEntry{
			{Date: "2024-01-16", RecipeID: "recipe1", Servings: 4, CoversDays: 1},
			{Date: "2024-01-17", RecipeID: "recipe2", Servings: 2, CoversDays: 1},
		},
	}

	removed := plan.RemoveEntry("2024-01-16")
	if !removed {
		t.Error("expected entry to be removed")
	}

	if len(plan.Entries) != 1 {
		t.Fatalf("expected 1 entry remaining, got %d", len(plan.Entries))
	}

	if plan.Entries[0].RecipeID != "recipe2" {
		t.Errorf("expected recipe2 to remain, got %s", plan.Entries[0].RecipeID)
	}
}

func TestMealPlan_RemoveEntry_NotFound(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      7,
		Entries:   []PlanEntry{},
	}

	removed := plan.RemoveEntry("2024-01-16")
	if removed {
		t.Error("expected entry not to be found")
	}
}

func TestMealPlan_IsDateInPlan(t *testing.T) {
	plan := &MealPlan{
		StartDate: "2024-01-15",
		Days:      3,
	}

	tests := []struct {
		date     string
		expected bool
	}{
		{"2024-01-15", true},
		{"2024-01-16", true},
		{"2024-01-17", true},
		{"2024-01-14", false},
		{"2024-01-18", false},
	}

	for _, tc := range tests {
		result := plan.IsDateInPlan(tc.date)
		if result != tc.expected {
			t.Errorf("IsDateInPlan(%q): expected %v, got %v", tc.date, tc.expected, result)
		}
	}
}
