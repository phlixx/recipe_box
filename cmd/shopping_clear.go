package cmd

import (
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var shoppingClearCmd = &cobra.Command{
	Use:  "clear",
	RunE: runShoppingClear,
}

func init() {
	shoppingClearCmd.Short = i18n.T(i18n.MsgCmdShoppingClearShort)
	shoppingClearCmd.Long = i18n.T(i18n.MsgCmdShoppingClearLong)
	shoppingCmd.AddCommand(shoppingClearCmd)
}

func runShoppingClear(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := shopping.NewService(store)
	if err := svc.Clear(); err != nil {
		return fmt.Errorf("failed to clear shopping list: %w", err)
	}

	ui.SuccessPrintf("%s\n", i18n.T(i18n.MsgShoppingCleared))
	return nil
}
