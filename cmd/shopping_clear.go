package cmd

import (
	"fmt"

	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var shoppingClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the shopping list",
	Long:  `Remove all items from your shopping list.`,
	RunE:  runShoppingClear,
}

func init() {
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

	fmt.Println("Shopping list cleared.")
	return nil
}
