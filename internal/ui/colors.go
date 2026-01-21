package ui

import (
	"os"

	"github.com/fatih/color"
)

// Color definitions based on ADR-008
var (
	// Title/Header styling - bold cyan
	Title = color.New(color.FgCyan, color.Bold)

	// Section headers - cyan (no bold)
	Section = color.New(color.FgCyan)

	// Category headers - yellow for visual grouping
	Category = color.New(color.FgYellow)

	// Success messages - green
	Success = color.New(color.FgGreen)

	// Error messages - red
	Error = color.New(color.FgRed)

	// Labels/prompts - dim/gray
	Label = color.New(color.Faint)

	// Highlighted values - bold white
	Highlight = color.New(color.Bold)

	// Muted text - dim
	Muted = color.New(color.Faint)

	// ID display - blue
	ID = color.New(color.FgBlue)
)

func init() {
	// Respect NO_COLOR environment variable (https://no-color.org/)
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		color.NoColor = true
	}
}

// Helper functions for formatted printing

// TitlePrintf prints a title/header
func TitlePrintf(format string, a ...any) {
	Title.Printf(format, a...)
}

// SectionPrintf prints a section header
func SectionPrintf(format string, a ...any) {
	Section.Printf(format, a...)
}

// CategoryPrintf prints a category header
func CategoryPrintf(format string, a ...any) {
	Category.Printf(format, a...)
}

// SuccessPrintf prints a success message
func SuccessPrintf(format string, a ...any) {
	Success.Printf(format, a...)
}

// ErrorPrintf prints an error message
func ErrorPrintf(format string, a ...any) {
	Error.Printf(format, a...)
}

// LabelPrintf prints a label/prompt
func LabelPrintf(format string, a ...any) {
	Label.Printf(format, a...)
}

// IDPrintf prints an ID
func IDPrintf(format string, a ...any) {
	ID.Printf(format, a...)
}

// Sprint functions for inline use

// TitleSprint returns a title-styled string
func TitleSprint(a ...any) string {
	return Title.Sprint(a...)
}

// SectionSprint returns a section-styled string
func SectionSprint(a ...any) string {
	return Section.Sprint(a...)
}

// CategorySprint returns a category-styled string
func CategorySprint(a ...any) string {
	return Category.Sprint(a...)
}

// SuccessSprint returns a success-styled string
func SuccessSprint(a ...any) string {
	return Success.Sprint(a...)
}

// ErrorSprint returns an error-styled string
func ErrorSprint(a ...any) string {
	return Error.Sprint(a...)
}

// IDSprint returns an ID-styled string
func IDSprint(a ...any) string {
	return ID.Sprint(a...)
}
