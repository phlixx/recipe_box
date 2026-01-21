package cmd

import (
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/recipe"
	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

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
