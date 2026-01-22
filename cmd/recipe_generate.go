package cmd

import (
	"fmt"
	"os"
	"strings"

	"recipe_box/internal/ai"
	"recipe_box/internal/config"
	"recipe_box/internal/i18n"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

const lastGeneratedFile = "last_generated.json"

// IsInteractive indicates whether we're running in interactive REPL mode
var IsInteractive bool

var (
	generatePrompt      string
	generateIngredients string
	generateCuisine     string
	generateVegetarian  bool
	generateQuick       bool
	generateServings    int
)

var recipeGenerateCmd = &cobra.Command{
	Use:  "generate",
	RunE: runRecipeGenerate,
}

func init() {
	recipeGenerateCmd.Short = i18n.T(i18n.MsgCmdRecipeGenerateShort)
	recipeGenerateCmd.Long = i18n.T(i18n.MsgCmdRecipeGenerateLong)
	recipeCmd.AddCommand(recipeGenerateCmd)

	recipeGenerateCmd.Flags().StringVarP(&generatePrompt, "prompt", "p", "", "Recipe description or request")
	recipeGenerateCmd.Flags().StringVarP(&generateIngredients, "ingredients", "i", "", "Comma-separated list of ingredients to use")
	recipeGenerateCmd.Flags().StringVarP(&generateCuisine, "cuisine", "c", "", "Cuisine type (e.g., italian, mexican, asian)")
	recipeGenerateCmd.Flags().BoolVarP(&generateVegetarian, "vegetarian", "v", false, "Generate a vegetarian recipe")
	recipeGenerateCmd.Flags().BoolVarP(&generateQuick, "quick", "q", false, "Quick recipe (under 30 min total)")
	recipeGenerateCmd.Flags().IntVarP(&generateServings, "servings", "s", 0, "Number of servings (default: 4)")
}

func runRecipeGenerate(cmd *cobra.Command, args []string) error {
	// Get API key from config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	apiKey, err := cfg.Get("api_key")
	if err != nil {
		return fmt.Errorf(i18n.T(i18n.MsgErrNoAPIKey))
	}

	// Create AI client
	client, err := ai.NewClient(apiKey)
	if err != nil {
		return err
	}

	// Check if any flags were provided
	flagsProvided := cmd.Flags().Changed("prompt") ||
		cmd.Flags().Changed("ingredients") ||
		cmd.Flags().Changed("cuisine") ||
		cmd.Flags().Changed("vegetarian") ||
		cmd.Flags().Changed("quick") ||
		cmd.Flags().Changed("servings")

	// If in interactive mode and no flags provided, prompt for options
	if IsInteractive && !flagsProvided {
		if err := promptForGenerateOptions(); err != nil {
			return err
		}
	}

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
		Language:    i18n.GetLanguage(),
		Servings:    generateServings,
	}

	fmt.Println()
	ui.LabelPrintf("%s\n", i18n.T(i18n.MsgRecipeGenerating))

	generatedRecipe, err := client.GenerateRecipe(opts)
	if err != nil {
		return fmt.Errorf("failed to generate recipe: %w", err)
	}

	// Save as last generated for "recipe save" command
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	if err := store.Save(lastGeneratedFile, generatedRecipe); err != nil {
		return fmt.Errorf("failed to save generated recipe: %w", err)
	}

	fmt.Println()
	printRecipe(generatedRecipe)

	// If in interactive mode, prompt to save
	if IsInteractive {
		fmt.Println()
		if promptYesNo(i18n.T(i18n.MsgGenerateSavePrompt)) {
			if err := saveGeneratedRecipe(store, generatedRecipe); err != nil {
				return err
			}
		}
	} else {
		fmt.Println()
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgRecipeSaveHint))
	}

	// Reset flags for next interactive use
	resetGenerateFlags()

	return nil
}

func promptForGenerateOptions() error {
	header := i18n.T(i18n.MsgGenerateHeader)
	ui.TitlePrintf("%s\n", header)
	fmt.Println(strings.Repeat("=", len(header)))
	fmt.Println()

	// Use go-prompt's Input for each question (works with REPL)
	emptyCompleter := func(d prompt.Document) []prompt.Suggest { return nil }

	// Recipe prompt (main request)
	generatePrompt = prompt.Input(i18n.T(i18n.MsgGeneratePrompt)+": ", emptyCompleter)

	// Ingredients
	generateIngredients = prompt.Input(i18n.T(i18n.MsgGenerateIngredients)+": ", emptyCompleter)

	// Cuisine
	generateCuisine = prompt.Input(i18n.T(i18n.MsgGenerateCuisine)+": ", emptyCompleter)

	// Servings
	servingsStr := prompt.Input(i18n.T(i18n.MsgGenerateServings)+": ", emptyCompleter)
	if servingsStr != "" {
		if n, err := parseServings(servingsStr); err == nil {
			generateServings = n
		}
	}

	// Vegetarian
	vegStr := prompt.Input(i18n.T(i18n.MsgGenerateVegetarian)+" ", emptyCompleter)
	generateVegetarian = isYes(vegStr)

	// Quick
	quickStr := prompt.Input(i18n.T(i18n.MsgGenerateQuick)+" ", emptyCompleter)
	generateQuick = isYes(quickStr)

	return nil
}

func parseServings(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("invalid servings")
	}
	return n, nil
}

func promptYesNo(question string) bool {
	emptyCompleter := func(d prompt.Document) []prompt.Suggest { return nil }
	answer := prompt.Input(question+" ", emptyCompleter)
	return isYes(answer)
}

func isYes(s string) bool {
	s = strings.ToLower(s)
	yes := strings.ToLower(i18n.T(i18n.MsgYes))
	return s == yes || s == "y" || s == "yes" || s == "j" || s == "ja"
}

func saveGeneratedRecipe(store *storage.Storage, r *recipe.Recipe) error {
	svc := recipe.NewService(store)
	if err := svc.Add(r); err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	// Remove the last_generated file
	if err := store.Delete(lastGeneratedFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not clean up temporary file: %v\n", err)
	}

	ui.SuccessPrintf(i18n.T(i18n.MsgRecipeSaved)+"\n", r.Title, r.ID)
	return nil
}

func resetGenerateFlags() {
	generatePrompt = ""
	generateIngredients = ""
	generateCuisine = ""
	generateVegetarian = false
	generateQuick = false
	generateServings = 0
}
