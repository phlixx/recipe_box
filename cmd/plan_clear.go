package cmd

import (
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var planClearCmd = &cobra.Command{
	Use:  "clear",
	RunE: runPlanClear,
}

func init() {
	planClearCmd.Short = i18n.T(i18n.MsgCmdPlanClearShort)
	planClearCmd.Long = i18n.T(i18n.MsgCmdPlanClearLong)
	planCmd.AddCommand(planClearCmd)
}

func runPlanClear(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := plan.NewService(store)
	if err := svc.Clear(); err != nil {
		return fmt.Errorf("failed to clear meal plan: %w", err)
	}

	ui.SuccessPrintf("%s\n", i18n.T(i18n.MsgPlanCleared))
	return nil
}
