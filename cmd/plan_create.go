package cmd

import (
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var planDays int

var planCreateCmd = &cobra.Command{
	Use:  "create",
	RunE: runPlanCreate,
}

func init() {
	planCreateCmd.Short = i18n.T(i18n.MsgCmdPlanCreateShort)
	planCreateCmd.Long = i18n.T(i18n.MsgCmdPlanCreateLong)
	planCmd.AddCommand(planCreateCmd)
	planCreateCmd.Flags().IntVarP(&planDays, "days", "d", 7, "Number of days for the meal plan")
}

func runPlanCreate(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := plan.NewService(store)
	p, err := svc.Create(planDays)
	if err != nil {
		return fmt.Errorf("failed to create meal plan: %w", err)
	}

	ui.SuccessPrintf(i18n.T(i18n.MsgPlanCreated)+"\n", p.Days, p.StartDate)
	fmt.Println()
	ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanEmptyHint))
	return nil
}
