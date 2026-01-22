package plan

import (
	"strings"
	"time"
)

// Day name mappings (English and German)
var dayNames = map[string]time.Weekday{
	// English
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sun":       time.Sunday,
	"mon":       time.Monday,
	"tue":       time.Tuesday,
	"wed":       time.Wednesday,
	"thu":       time.Thursday,
	"fri":       time.Friday,
	"sat":       time.Saturday,
	// German
	"sonntag":    time.Sunday,
	"montag":     time.Monday,
	"dienstag":   time.Tuesday,
	"mittwoch":   time.Wednesday,
	"donnerstag": time.Thursday,
	"freitag":    time.Friday,
	"samstag":    time.Saturday,
	"so":         time.Sunday,
	"mo":         time.Monday,
	"di":         time.Tuesday,
	"mi":         time.Wednesday,
	"do":         time.Thursday,
	"fr":         time.Friday,
	"sa":         time.Saturday,
}

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

// ParseDayToDate converts a day name or ISO date to an ISO date within the plan
// Supports: "monday", "mon", "montag", "mo", or "2024-01-15"
func (p *MealPlan) ParseDayToDate(dayInput string) (string, bool) {
	input := strings.ToLower(strings.TrimSpace(dayInput))

	// Check if it's already an ISO date
	if _, err := time.Parse("2006-01-02", input); err == nil {
		// Verify it's in the plan
		for _, date := range p.GetDates() {
			if date == input {
				return input, true
			}
		}
		return "", false
	}

	// Check if it's a day name
	weekday, ok := dayNames[input]
	if !ok {
		return "", false
	}

	// Find the matching date in the plan
	for _, date := range p.GetDates() {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}
		if t.Weekday() == weekday {
			return date, true
		}
	}

	return "", false
}

// AddEntry adds or replaces an entry for a specific date
func (p *MealPlan) AddEntry(entry PlanEntry) {
	// Remove existing entry for this date if any
	p.RemoveEntry(entry.Date)

	// Add the new entry
	p.Entries = append(p.Entries, entry)
}

// RemoveEntry removes the entry for a specific date
func (p *MealPlan) RemoveEntry(date string) bool {
	for i, entry := range p.Entries {
		if entry.Date == date {
			p.Entries = append(p.Entries[:i], p.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// IsDateInPlan checks if a date is within the plan's date range
func (p *MealPlan) IsDateInPlan(date string) bool {
	for _, d := range p.GetDates() {
		if d == date {
			return true
		}
	}
	return false
}

// generateID creates a short unique ID
func generateID() string {
	return time.Now().Format("20060102150405")
}
