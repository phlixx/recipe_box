package cmd

import (
	"fmt"
	"strings"

	"recipe_box/internal/ai"
	"recipe_box/internal/config"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

const lastGeneratedFile = "last_generated.json"

var (
	generatePrompt      string
	generateIngredients string
	generateCuisine     string
	generateVegetarian  bool
	generateQuick       bool
)

var recipeGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a recipe using AI",
	Long: `Generate a recipe using Ollama (local AI).

Requires Ollama to be running: ollama serve
Default model: llama3.2 (change with: recipe_box config set model <name>)

Examples:
  recipe_box recipe generate --prompt "quick pasta dish"
  recipe_box recipe generate --ingredients "chicken, lemon, garlic"
  recipe_box recipe generate --cuisine italian --vegetarian
  recipe_box recipe generate --quick --cuisine mexican`,
	RunE: runRecipeGenerate,
}

func init() {
	recipeCmd.AddCommand(recipeGenerateCmd)

	recipeGenerateCmd.Flags().StringVarP(&generatePrompt, "prompt", "p", "", "Recipe description or request")
	recipeGenerateCmd.Flags().StringVarP(&generateIngredients, "ingredients", "i", "", "Comma-separated list of ingredients to use")
	recipeGenerateCmd.Flags().StringVarP(&generateCuisine, "cuisine", "c", "", "Cuisine type (e.g., italian, mexican, asian)")
	recipeGenerateCmd.Flags().BoolVarP(&generateVegetarian, "vegetarian", "v", false, "Generate a vegetarian recipe")
	recipeGenerateCmd.Flags().BoolVarP(&generateQuick, "quick", "q", false, "Quick recipe (under 30 min total)")
}

func runRecipeGenerate(cmd *cobra.Command, args []string) error {
	// Get optional model from config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	model, _ := cfg.Get("model") // Empty string uses default

	// Create AI client
	client := ai.NewClient(model)

	// Parse ingredients
	var ingredients []string
	if generateIngredients != "" {
		for _, ing := range strings.Split(generateIngredients, ",") {
			ing = strings.TrimSpace(ing)
			if ing != "" {
				ingredients = append(ingredients, ing)
			}
		}
	}

	// Build options
	opts := ai.GenerateOptions{
		Prompt:      generatePrompt,
		Ingredients: ingredients,
		Cuisine:     generateCuisine,
		Vegetarian:  generateVegetarian,
		Quick:       generateQuick,
	}

	fmt.Println("Generating recipe...")

	recipe, err := client.GenerateRecipe(opts)
	if err != nil {
		return fmt.Errorf("failed to generate recipe: %w", err)
	}

	// Save as last generated for "recipe save" command
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	if err := store.Save(lastGeneratedFile, recipe); err != nil {
		return fmt.Errorf("failed to save generated recipe: %w", err)
	}

	fmt.Println()
	printRecipe(recipe)
	fmt.Println("\nUse 'recipe_box recipe save' to add this recipe to your collection.")

	return nil
}
