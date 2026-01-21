package cmd

import (
	"fmt"

	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var recipeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recipes",
	Long:  `Display all recipes in your collection.`,
	RunE:  runRecipeList,
}

func init() {
	recipeCmd.AddCommand(recipeListCmd)
}

func runRecipeList(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := recipe.NewService(store)
	recipes, err := svc.List()
	if err != nil {
		return fmt.Errorf("failed to list recipes: %w", err)
	}

	if len(recipes) == 0 {
		fmt.Println("No recipes found. Use 'recipe_box recipe add' to add one.")
		return nil
	}

	fmt.Printf("Found %d recipe(s):\n\n", len(recipes))
	for _, r := range recipes {
		totalTime := r.PrepTime + r.CookTime
		fmt.Printf("  %s  %s\n", r.ID, r.Title)
		if totalTime > 0 {
			fmt.Printf("         %d servings, %d min\n", r.Servings, totalTime)
		} else {
			fmt.Printf("         %d servings\n", r.Servings)
		}
	}

	return nil
}
