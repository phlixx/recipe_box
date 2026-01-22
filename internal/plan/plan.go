package plan

import "time"

// MealPlan represents a weekly meal plan
type MealPlan struct {
	ID        string      `json:"id"`
	StartDate string      `json:"start_date"` // ISO 8601 (YYYY-MM-DD)
	Days      int         `json:"days"`
	Entries   []PlanEntry `json:"entries"`
}

// PlanEntry represents a single meal assignment
type PlanEntry struct {
	Date       string `json:"date"` // ISO 8601 (YYYY-MM-DD)
	RecipeID   string `json:"recipe_id"`
	Servings   int    `json:"servings"`
	CoversDays int    `json:"covers_days"` // leftover tracking (default 1)
}

// NewMealPlan creates a new meal plan starting from today
func NewMealPlan(days int) *MealPlan {
	if days <= 0 {
		days = 7
	}
	return &MealPlan{
		ID:        generateID(),
		StartDate: time.Now().Format("2006-01-02"),
		Days:      days,
		Entries:   []PlanEntry{},
	}
}

// GetDates returns all dates covered by the plan
func (p *MealPlan) GetDates() []string {
	start, err := time.Parse("2006-01-02", p.StartDate)
	if err != nil {
		return nil
	}

	dates := make([]string, p.Days)
	for i := 0; i < p.Days; i++ {
		dates[i] = start.AddDate(0, 0, i).Format("2006-01-02")
	}
	return dates
}

// GetEntryForDate returns the entry for a specific date (if any)
func (p *MealPlan) GetEntryForDate(date string) *PlanEntry {
	for i := range p.Entries {
		if p.Entries[i].Date == date {
			return &p.Entries[i]
		}
	}
	return nil
}

// IsLeftoverDay checks if a date is covered by leftovers from a previous entry
func (p *MealPlan) IsLeftoverDay(date string) (bool, *PlanEntry) {
	for i := range p.Entries {
		entry := &p.Entries[i]
		if entry.CoversDays <= 1 {
			continue
		}

		entryDate, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			continue
		}

		// Check if target date falls within the leftover range (excluding the main date)
		for d := 1; d < entry.CoversDays; d++ {
			leftoverDate := entryDate.AddDate(0, 0, d)
			if leftoverDate.Format("2006-01-02") == date {
				return true, entry
			}
		}
	}
	return false, nil
}

// generateID creates a short unique ID
func generateID() string {
	return time.Now().Format("20060102150405")
}
