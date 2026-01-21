package cmd

import (
	"errors"
	"fmt"
	"strings"

	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var recipeViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View a recipe",
	Long:  `Display the full details of a recipe.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRecipeView,
}

func init() {
	recipeCmd.AddCommand(recipeViewCmd)
}

func runRecipeView(cmd *cobra.Command, args []string) error {
	id := args[0]

	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := recipe.NewService(store)
	r, err := svc.Get(id)
	if err != nil {
		if errors.Is(err, recipe.ErrNotFound) {
			return fmt.Errorf("recipe %q not found", id)
		}
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	printRecipe(r)
	return nil
}

func printRecipe(r *recipe.Recipe) {
	fmt.Printf("%s\n", r.Title)
	fmt.Println(strings.Repeat("=", len(r.Title)))
	fmt.Println()

	if r.Description != "" {
		fmt.Printf("%s\n\n", r.Description)
	}

	// Metadata
	fmt.Printf("Servings: %d\n", r.Servings)
	if r.PrepTime > 0 {
		fmt.Printf("Prep time: %d min\n", r.PrepTime)
	}
	if r.CookTime > 0 {
		fmt.Printf("Cook time: %d min\n", r.CookTime)
	}
	if r.PrepTime > 0 || r.CookTime > 0 {
		fmt.Printf("Total time: %d min\n", r.PrepTime+r.CookTime)
	}

	// Tags
	if len(r.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(r.Tags, ", "))
	}

	// Ingredients
	fmt.Println("\nIngredients:")
	for _, ing := range r.Ingredients {
		if ing.Quantity > 0 && ing.Unit != "" {
			fmt.Printf("  - %.4g %s %s\n", ing.Quantity, ing.Unit, ing.Name)
		} else if ing.Quantity > 0 {
			fmt.Printf("  - %.4g %s\n", ing.Quantity, ing.Name)
		} else {
			fmt.Printf("  - %s\n", ing.Name)
		}
	}

	// Steps
	fmt.Println("\nSteps:")
	for i, step := range r.Steps {
		fmt.Printf("  %d. %s\n", i+1, step)
	}

	fmt.Printf("\n[ID: %s]\n", r.ID)
}
