package cmd

import (
	"errors"
	"fmt"

	"recipe_box/internal/i18n"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var recipeDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a recipe",
	Long:  `Remove a recipe from your collection.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRecipeDelete,
}

func init() {
	recipeCmd.AddCommand(recipeDeleteCmd)
}

func runRecipeDelete(cmd *cobra.Command, args []string) error {
	id := args[0]

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := recipe.NewService(store)

	// Get recipe first to show title in confirmation
	r, err := svc.Get(id)
	if err != nil {
		if errors.Is(err, recipe.ErrNotFound) {
			return fmt.Errorf(i18n.T(i18n.MsgRecipeNotFound), id)
		}
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if err := svc.Delete(id); err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	ui.SuccessPrintf(i18n.T(i18n.MsgRecipeDeleted)+"\n", r.Title)
	return nil
}
