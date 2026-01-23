package cmd

import (
	"fmt"
	"strings"
	"time"

	"recipe_box/internal/ai"
	"recipe_box/internal/config"
	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var (
	planGenerateDays       int
	planGenerateCuisine    string
	planGenerateVegetarian bool
	planGenerateQuick      bool
)

var planGenerateCmd = &cobra.Command{
	Use:  "generate",
	RunE: runPlanGenerate,
}

func init() {
	planGenerateCmd.Short = i18n.T(i18n.MsgCmdPlanGenerateShort)
	planGenerateCmd.Long = i18n.T(i18n.MsgCmdPlanGenerateLong)
	planCmd.AddCommand(planGenerateCmd)

	planGenerateCmd.Flags().IntVarP(&planGenerateDays, "days", "d", 7, "Number of days to plan")
	planGenerateCmd.Flags().StringVarP(&planGenerateCuisine, "cuisine", "c", "", "Cuisine preference (e.g., italian, asian)")
	planGenerateCmd.Flags().BoolVarP(&planGenerateVegetarian, "vegetarian", "v", false, "Generate vegetarian meals only")
	planGenerateCmd.Flags().BoolVarP(&planGenerateQuick, "quick", "q", false, "Quick meals only (under 30 min)")
}

func runPlanGenerate(cmd *cobra.Command, args []string) error {
	// Get API key
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	apiKey, err := cfg.Get("api_key")
	if err != nil {
		return fmt.Errorf(i18n.T(i18n.MsgErrNoAPIKey))
	}

	client, err := ai.NewClient(apiKey)
	if err != nil {
		return err
	}

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	planSvc := plan.NewService(store)

	// Generate dates starting from today
	dates := generateDates(planGenerateDays)

	opts := ai.MealPlanOptions{
		Days:       planGenerateDays,
		Cuisine:    planGenerateCuisine,
		Vegetarian: planGenerateVegetarian,
		Quick:      planGenerateQuick,
		Language:   i18n.GetLanguage(),
	}

	fmt.Println()
	ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanGenerating))

	suggestions, err := client.GenerateMealPlan(opts, dates)
	if err != nil {
		return fmt.Errorf("failed to generate meal plan: %w", err)
	}

	// Display suggestions
	fmt.Println()
	displayMealSuggestions(suggestions)

	// Ask for approval
	fmt.Println()
	if !promptYesNo(i18n.T(i18n.MsgPlanApprovePrompt)) {
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanDiscarded))
		resetPlanGenerateFlags()
		return nil
	}

	// Create the meal plan
	p, err := planSvc.Create(planGenerateDays)
	if err != nil {
		return fmt.Errorf("failed to create meal plan: %w", err)
	}

	// Ask if user wants to create full recipes
	fmt.Println()
	createRecipes := promptYesNo(i18n.T(i18n.MsgPlanCreateRecipesPrompt))

	if createRecipes {
		// Create recipes and add to plan
		recipeSvc := recipe.NewService(store)
		for _, suggestion := range suggestions.Suggestions {
			ui.LabelPrintf(i18n.T(i18n.MsgPlanCreatingRecipes)+"\n", suggestion.MealName)

			// Generate full recipe for this meal
			recipeOpts := ai.GenerateOptions{
				Prompt:     suggestion.MealName + ": " + suggestion.Description,
				Cuisine:    planGenerateCuisine,
				Vegetarian: planGenerateVegetarian,
				Quick:      planGenerateQuick,
				Language:   i18n.GetLanguage(),
				Servings:   suggestion.Servings,
			}

			r, err := client.GenerateRecipe(recipeOpts)
			if err != nil {
				ui.ErrorPrintf("Failed to create recipe for %s: %v\n", suggestion.MealName, err)
				continue
			}

			// Save the recipe
			if err := recipeSvc.Add(r); err != nil {
				ui.ErrorPrintf("Failed to save recipe %s: %v\n", r.Title, err)
				continue
			}

			ui.SuccessPrintf(i18n.T(i18n.MsgPlanRecipeCreated)+"\n", r.Title, r.ID)

			// Add to plan
			if err := planSvc.AddEntry(suggestion.Date, r.ID, suggestion.Servings, 1); err != nil {
				ui.ErrorPrintf("Failed to add recipe to plan: %v\n", err)
			}
		}

		fmt.Println()
		ui.SuccessPrintf("%s\n", i18n.T(i18n.MsgPlanAllRecipesCreated))
	} else {
		// Just save the plan structure without recipes
		// The plan is already created, just inform the user
		ui.SuccessPrintf(i18n.T(i18n.MsgPlanCreated)+"\n", p.Days, p.StartDate)
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanEmptyHint))
	}

	resetPlanGenerateFlags()
	return nil
}

func generateDates(days int) []string {
	dates := make([]string, days)
	start := time.Now()
	for i := 0; i < days; i++ {
		dates[i] = start.AddDate(0, 0, i).Format("2006-01-02")
	}
	return dates
}

func displayMealSuggestions(suggestions *ai.MealPlanSuggestions) {
	ui.TitlePrintf(i18n.T(i18n.MsgPlanGeneratedHeader)+"\n", len(suggestions.Suggestions))
	fmt.Println(strings.Repeat("=", 40))

	for _, s := range suggestions.Suggestions {
		dayName := s.Day
		if dayName == "" {
			// Fall back to parsing date if day name not provided
			t, err := time.Parse("2006-01-02", s.Date)
			if err == nil {
				dayName = t.Weekday().String()
			}
		}

		fmt.Println()
		ui.SectionPrintf("%s (%s):\n", dayName, s.Date)
		fmt.Printf("  %s\n", s.MealName)
		if s.Description != "" {
			ui.LabelPrintf("  %s\n", s.Description)
		}

		// Time info
		totalTime := s.PrepTime + s.CookTime
		if totalTime > 0 {
			ui.LabelPrintf("  %d %s | %s\n", totalTime, i18n.T(i18n.MsgMinutes), fmt.Sprintf(i18n.T(i18n.MsgPlanServingsInfo), s.Servings))
		} else if s.Servings > 0 {
			ui.LabelPrintf("  %s\n", fmt.Sprintf(i18n.T(i18n.MsgPlanServingsInfo), s.Servings))
		}
	}
}

func resetPlanGenerateFlags() {
	planGenerateDays = 7
	planGenerateCuisine = ""
	planGenerateVegetarian = false
	planGenerateQuick = false
}
