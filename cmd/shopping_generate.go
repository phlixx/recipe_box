package cmd

import (
	"fmt"

	"recipe_box/internal/recipe"
	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var shoppingGenerateCmd = &cobra.Command{
	Use:   "generate <recipe-id>",
	Short: "Generate shopping list from a recipe",
	Long: `Add ingredients from a recipe to your shopping list.

Examples:
  recipe_box shopping generate abc123
  recipe_box shopping generate abc123 def456`,
	Args: cobra.MinimumNArgs(1),
	RunE: runShoppingGenerate,
}

func init() {
	shoppingCmd.AddCommand(shoppingGenerateCmd)
}

func runShoppingGenerate(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	recipeSvc := recipe.NewService(store)
	shoppingSvc := shopping.NewService(store)

	for _, id := range args {
		r, err := recipeSvc.Get(id)
		if err != nil {
			if err == recipe.ErrNotFound {
				fmt.Printf("Recipe %s not found, skipping\n", id)
				continue
			}
			return fmt.Errorf("failed to get recipe %s: %w", id, err)
		}

		if err := shoppingSvc.GenerateFromRecipe(r); err != nil {
			return fmt.Errorf("failed to add ingredients from %s: %w", id, err)
		}

		fmt.Printf("Added %d ingredients from '%s' to shopping list\n", len(r.Ingredients), r.Title)
	}

	fmt.Println("\nUse 'recipe_box shopping show' to view your shopping list.")
	return nil
}
