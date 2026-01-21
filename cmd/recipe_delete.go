package cmd

import (
	"errors"
	"fmt"

	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"

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
			return fmt.Errorf("recipe %q not found", id)
		}
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if err := svc.Delete(id); err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	fmt.Printf("Deleted recipe '%s' (%s)\n", r.Title, r.ID)
	return nil
}
