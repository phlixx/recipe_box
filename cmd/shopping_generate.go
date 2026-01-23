package cmd

import (
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/recipe"
	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var (
	shoppingServings int
	shoppingFromPlan bool
	shoppingPlanDays int
)

var shoppingGenerateCmd = &cobra.Command{
	Use:  "generate [recipe-id...]",
	Args: cobra.ArbitraryArgs,
	RunE: runShoppingGenerate,
}

func init() {
	shoppingGenerateCmd.Short = i18n.T(i18n.MsgCmdShoppingGenerateShort)
	shoppingGenerateCmd.Long = i18n.T(i18n.MsgCmdShoppingGenerateLong)
	shoppingCmd.AddCommand(shoppingGenerateCmd)
	shoppingGenerateCmd.Flags().IntVarP(&shoppingServings, "servings", "s", 0, "Scale recipe to this number of servings")
	shoppingGenerateCmd.Flags().BoolVarP(&shoppingFromPlan, "plan", "p", false, "Generate from meal plan")
	shoppingGenerateCmd.Flags().IntVarP(&shoppingPlanDays, "days", "d", 0, "Limit to first N days of meal plan (requires --plan)")
}

func runShoppingGenerate(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	recipeSvc := recipe.NewService(store)
	shoppingSvc := shopping.NewService(store)

	// Handle --plan flag
	if shoppingFromPlan {
		return runShoppingGenerateFromPlan(store, recipeSvc, shoppingSvc)
	}

	// Require at least one recipe ID when not using --plan
	if len(args) == 0 {
		return fmt.Errorf(i18n.T(i18n.MsgShoppingGenerateNoArgs))
	}

	addedCount := 0
	for _, id := range args {
		r, err := recipeSvc.Get(id)
		if err != nil {
			if err == recipe.ErrNotFound {
				ui.ErrorPrintf(i18n.T(i18n.MsgRecipeNotFound)+"\n", id)
				continue
			}
			return fmt.Errorf("failed to get recipe %s: %w", id, err)
		}

		// Scale recipe if --servings flag is provided
		if shoppingServings > 0 && shoppingServings != r.Servings {
			r = r.Scale(shoppingServings)
		}

		if err := shoppingSvc.GenerateFromRecipe(r); err != nil {
			return fmt.Errorf("failed to add ingredients from %s: %w", id, err)
		}

		ui.SuccessPrintf(i18n.T(i18n.MsgShoppingAdded)+"\n", len(r.Ingredients), r.Title)
		addedCount++
	}

	// Only show hint if items were actually added
	if addedCount > 0 {
		fmt.Println()
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgShoppingShowHint))
	}
	return nil
}

func runShoppingGenerateFromPlan(store *storage.Storage, recipeSvc *recipe.Service, shoppingSvc *shopping.Service) error {
	planSvc := plan.NewService(store)

	mealPlan, err := planSvc.Get()
	if err != nil {
		return fmt.Errorf("failed to get meal plan: %w", err)
	}
	if mealPlan == nil {
		ui.ErrorPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlan))
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlanHint))
		return nil
	}

	if len(mealPlan.Entries) == 0 {
		ui.ErrorPrintf("%s\n", i18n.T(i18n.MsgPlanEmpty))
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanEmptyHint))
		return nil
	}

	// Get dates to process
	dates := mealPlan.GetDates()
	if shoppingPlanDays > 0 && shoppingPlanDays < len(dates) {
		dates = dates[:shoppingPlanDays]
	}

	// Track which recipes we've already processed (to avoid duplicates from CoversDays)
	processedRecipes := make(map[string]bool)
	addedCount := 0

	for _, date := range dates {
		entry := mealPlan.GetEntryForDate(date)
		if entry == nil {
			continue // Empty day
		}

		// Skip if we've already processed this recipe
		if processedRecipes[entry.RecipeID] {
			continue
		}
		processedRecipes[entry.RecipeID] = true

		r, err := recipeSvc.Get(entry.RecipeID)
		if err != nil {
			if err == recipe.ErrNotFound {
				ui.ErrorPrintf(i18n.T(i18n.MsgRecipeNotFound)+"\n", entry.RecipeID)
				continue
			}
			return fmt.Errorf("failed to get recipe %s: %w", entry.RecipeID, err)
		}

		// Scale recipe based on plan entry servings
		if entry.Servings > 0 && entry.Servings != r.Servings {
			r = r.Scale(entry.Servings)
		}

		if err := shoppingSvc.GenerateFromRecipe(r); err != nil {
			return fmt.Errorf("failed to add ingredients from %s: %w", entry.RecipeID, err)
		}

		ui.SuccessPrintf(i18n.T(i18n.MsgShoppingAdded)+"\n", len(r.Ingredients), r.Title)
		addedCount++
	}

	if addedCount == 0 {
		ui.ErrorPrintf("%s\n", i18n.T(i18n.MsgShoppingPlanNoRecipes))
		return nil
	}

	fmt.Println()
	ui.LabelPrintf("%s\n", i18n.T(i18n.MsgShoppingShowHint))
	return nil
}
