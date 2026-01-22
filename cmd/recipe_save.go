package cmd

import (
	"fmt"
	"os"

	"recipe_box/internal/i18n"
	internalRecipe "recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var recipeSaveCmd = &cobra.Command{
	Use:  "save",
	RunE: runRecipeSave,
}

func init() {
	recipeSaveCmd.Short = i18n.T(i18n.MsgCmdRecipeSaveShort)
	recipeSaveCmd.Long = i18n.T(i18n.MsgCmdRecipeSaveLong)
	recipeCmd.AddCommand(recipeSaveCmd)
}

func runRecipeSave(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	// Load last generated recipe
	if !store.Exists(lastGeneratedFile) {
		return fmt.Errorf(i18n.T(i18n.MsgRecipeNoGenerated))
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

	ui.SuccessPrintf(i18n.T(i18n.MsgRecipeSaved)+"\n", recipe.Title, recipe.ID)
	return nil
}
