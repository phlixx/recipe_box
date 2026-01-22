package cmd

import (
	"errors"
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var planRemoveCmd = &cobra.Command{
	Use:  "remove <day>",
	Args: cobra.ExactArgs(1),
	RunE: runPlanRemove,
}

func init() {
	planRemoveCmd.Short = i18n.T(i18n.MsgCmdPlanRemoveShort)
	planRemoveCmd.Long = i18n.T(i18n.MsgCmdPlanRemoveLong)
	planCmd.AddCommand(planRemoveCmd)
}

func runPlanRemove(cmd *cobra.Command, args []string) error {
	dayInput := args[0]

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	planSvc := plan.NewService(store)

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

	// Remove entry from plan
	removed, err := planSvc.RemoveEntry(date)
	if err != nil {
		if errors.Is(err, plan.ErrNoPlan) {
			ui.ErrorPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlan))
			return nil
		}
		if errors.Is(err, plan.ErrDateNotInPlan) {
			return fmt.Errorf(i18n.T(i18n.MsgPlanDateNotInPlan), date)
		}
		return fmt.Errorf("failed to remove entry: %w", err)
	}

	if !removed {
		// No entry existed for this date, but that's not an error
		dayName := formatDayName(date)
		ui.LabelPrintf("%s\n", fmt.Sprintf(i18n.T(i18n.MsgPlanDayFormat), dayName, date)+" "+i18n.T(i18n.MsgPlanEmptyDay))
		return nil
	}

	// Format success message
	dayName := formatDayName(date)
	ui.SuccessPrintf(i18n.T(i18n.MsgPlanEntryRemoved)+"\n", dayName)
	return nil
}
