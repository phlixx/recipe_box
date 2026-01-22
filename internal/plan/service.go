package plan

import (
	"errors"

	"recipe_box/internal/storage"
)

// Service errors
var (
	ErrNoPlan        = errors.New("no meal plan found")
	ErrDateNotInPlan = errors.New("date not in meal plan")
)

const mealPlanFile = "meal_plan.json"

// Storage defines the interface for persistence
type Storage interface {
	Save(filename string, data any) error
	Load(filename string, dest any) error
	Exists(filename string) bool
	Delete(filename string) error
}

// Service handles meal plan operations
type Service struct {
	storage Storage
}

// NewService creates a new plan service
func NewService(s *storage.Storage) *Service {
	return &Service{storage: s}
}

// Create creates a new meal plan (replaces existing)
func (s *Service) Create(days int) (*MealPlan, error) {
	plan := NewMealPlan(days)
	if err := s.storage.Save(mealPlanFile, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

// Get returns the current meal plan
func (s *Service) Get() (*MealPlan, error) {
	var plan MealPlan

	if !s.storage.Exists(mealPlanFile) {
		return nil, nil
	}

	err := s.storage.Load(mealPlanFile, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// Exists checks if a meal plan exists
func (s *Service) Exists() bool {
	return s.storage.Exists(mealPlanFile)
}

// Clear removes the current meal plan
func (s *Service) Clear() error {
	if !s.storage.Exists(mealPlanFile) {
		return nil
	}
	return s.storage.Delete(mealPlanFile)
}

// Save persists the meal plan
func (s *Service) Save(plan *MealPlan) error {
	return s.storage.Save(mealPlanFile, plan)
}

// AddEntry adds a recipe to a specific date in the plan
func (s *Service) AddEntry(date string, recipeID string, servings int, coversDays int) error {
	plan, err := s.Get()
	if err != nil {
		return err
	}
	if plan == nil {
		return ErrNoPlan
	}

	if !plan.IsDateInPlan(date) {
		return ErrDateNotInPlan
	}

	if coversDays < 1 {
		coversDays = 1
	}

	entry := PlanEntry{
		Date:       date,
		RecipeID:   recipeID,
		Servings:   servings,
		CoversDays: coversDays,
	}

	plan.AddEntry(entry)
	return s.Save(plan)
}

// RemoveEntry removes the entry for a specific date
func (s *Service) RemoveEntry(date string) (bool, error) {
	plan, err := s.Get()
	if err != nil {
		return false, err
	}
	if plan == nil {
		return false, ErrNoPlan
	}

	if !plan.IsDateInPlan(date) {
		return false, ErrDateNotInPlan
	}

	removed := plan.RemoveEntry(date)
	if !removed {
		return false, nil
	}

	return true, s.Save(plan)
}
