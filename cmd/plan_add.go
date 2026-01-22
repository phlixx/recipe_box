package cmd

import (
	"errors"
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var (
	planAddServings int
	planAddDays     int
)

var planAddCmd = &cobra.Command{
	Use:  "add <day> <recipe-id>",
	Args: cobra.ExactArgs(2),
	RunE: runPlanAdd,
}

func init() {
	planAddCmd.Short = i18n.T(i18n.MsgCmdPlanAddShort)
	planAddCmd.Long = i18n.T(i18n.MsgCmdPlanAddLong)
	planCmd.AddCommand(planAddCmd)
	planAddCmd.Flags().IntVarP(&planAddServings, "servings", "s", 0, "Number of servings")
	planAddCmd.Flags().IntVarP(&planAddDays, "days", "d", 1, "Number of days this meal covers (for leftovers)")
}

func runPlanAdd(cmd *cobra.Command, args []string) error {
	dayInput := args[0]
	recipeID := args[1]

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	planSvc := plan.NewService(store)
	recipeSvc := recipe.NewService(store)

	// Get the plan to parse day name
	p, err := planSvc.Get()
	if err != nil {
		return fmt.Errorf("failed to get meal plan: %w", err)
	}
	if p == nil {
		ui.ErrorPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlan))
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlanHint))
		return nil
	}

	// Parse day input to date
	date, ok := p.ParseDayToDate(dayInput)
	if !ok {
		return fmt.Errorf(i18n.T(i18n.MsgPlanDateNotInPlan), dayInput)
	}

	// Verify recipe exists
	r, err := recipeSvc.Get(recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}
	if r == nil {
		return fmt.Errorf(i18n.T(i18n.MsgRecipeNotFound), recipeID)
	}

	// Use recipe's default servings if not specified
	servings := planAddServings
	if servings == 0 {
		servings = r.Servings
	}

	// Add entry to plan
	err = planSvc.AddEntry(date, recipeID, servings, planAddDays)
	if err != nil {
		if errors.Is(err, plan.ErrNoPlan) {
			ui.ErrorPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlan))
			return nil
		}
		if errors.Is(err, plan.ErrDateNotInPlan) {
			return fmt.Errorf(i18n.T(i18n.MsgPlanDateNotInPlan), date)
		}
		return fmt.Errorf("failed to add entry: %w", err)
	}

	// Format success message
	dayName := formatDayName(date)
	ui.SuccessPrintf(i18n.T(i18n.MsgPlanEntryAdded)+"\n", r.Title, dayName)
	return nil
}
