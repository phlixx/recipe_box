package cmd

import (
	"fmt"

	"recipe_box/internal/i18n"
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
		fmt.Println(i18n.T(i18n.MsgRecipeListEmpty))
		return nil
	}

	fmt.Printf(i18n.T(i18n.MsgRecipeListFound)+"\n\n", len(recipes))
	for _, r := range recipes {
		totalTime := r.PrepTime + r.CookTime
		fmt.Printf("  %s  %s\n", r.ID, r.Title)
		if totalTime > 0 {
			fmt.Printf("         %d %s, %d %s\n", r.Servings, i18n.T(i18n.MsgServings), totalTime, i18n.T(i18n.MsgMinutes))
		} else {
			fmt.Printf("         %d %s\n", r.Servings, i18n.T(i18n.MsgServings))
		}
	}

	return nil
}
