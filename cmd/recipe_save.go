package cmd

import (
	"fmt"
	"os"

	internalRecipe "recipe_box/internal/recipe"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var recipeSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the last generated recipe",
	Long: `Save the most recently generated recipe to your collection.

Run 'recipe generate' first to create a recipe, then use this command to save it.`,
	RunE: runRecipeSave,
}

func init() {
	recipeCmd.AddCommand(recipeSaveCmd)
}

func runRecipeSave(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	// Load last generated recipe
	if !store.Exists(lastGeneratedFile) {
		return fmt.Errorf("no generated recipe to save - run 'recipe generate' first")
	}

	var recipe internalRecipe.Recipe
	if err := store.Load(lastGeneratedFile, &recipe); err != nil {
		return fmt.Errorf("failed to load generated recipe: %w", err)
	}

	// Add to collection
	svc := internalRecipe.NewService(store)
	if err := svc.Add(&recipe); err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	// Remove the last_generated file
	if err := store.Delete(lastGeneratedFile); err != nil {
		// Not critical, just warn
		fmt.Fprintf(os.Stderr, "Warning: could not clean up temporary file: %v\n", err)
	}

	fmt.Printf("Saved recipe '%s' (%s)\n", recipe.Title, recipe.ID)
	return nil
}
