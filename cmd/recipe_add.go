package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var recipeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new recipe",
	Long:  `Interactively add a new recipe to your collection.`,
	RunE:  runRecipeAdd,
}

func init() {
	recipeCmd.AddCommand(recipeAddCmd)
}

func runRecipeAdd(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Add a new recipe")
	fmt.Println("================")
	fmt.Println()

	// Title
	title, err := prompt(reader, "Title")
	if err != nil {
		return err
	}

	// Description
	description, err := prompt(reader, "Description (optional)")
	if err != nil {
		return err
	}

	// Servings
	servingsStr, err := prompt(reader, "Servings")
	if err != nil {
		return err
	}
	servings, _ := strconv.Atoi(servingsStr)
	if servings == 0 {
		servings = 4
	}

	// Prep time
	prepTimeStr, err := prompt(reader, "Prep time (minutes)")
	if err != nil {
		return err
	}
	prepTime, _ := strconv.Atoi(prepTimeStr)

	// Cook time
	cookTimeStr, err := prompt(reader, "Cook time (minutes)")
	if err != nil {
		return err
	}
	cookTime, _ := strconv.Atoi(cookTimeStr)

	// Ingredients
	fmt.Println("\nIngredients (enter empty line when done):")
	var ingredients []recipe.Ingredient
	for {
		line, err := prompt(reader, "  Ingredient (e.g., '2 cups flour')")
		if err != nil {
			return err
		}
		if line == "" {
			break
		}
		ingredients = append(ingredients, parseIngredient(line))
	}

	// Steps
	fmt.Println("\nSteps (enter empty line when done):")
	var steps []string
	stepNum := 1
	for {
		line, err := prompt(reader, fmt.Sprintf("  Step %d", stepNum))
		if err != nil {
			return err
		}
		if line == "" {
			break
		}
		steps = append(steps, line)
		stepNum++
	}

	// Tags
	tagsStr, err := prompt(reader, "Tags (comma-separated, optional)")
	if err != nil {
		return err
	}
	var tags []string
	if tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	// Create recipe
	r := &recipe.Recipe{
		Title:       title,
		Description: description,
		Servings:    servings,
		PrepTime:    prepTime,
		CookTime:    cookTime,
		Ingredients: ingredients,
		Steps:       steps,
		Tags:        tags,
		Source:      recipe.SourceManual,
		Language:    recipe.LangEN,
	}

	// Save
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := recipe.NewService(store)
	if err := svc.Add(r); err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	fmt.Printf("\nRecipe '%s' saved with ID: %s\n", r.Title, r.ID)
	return nil
}

func prompt(reader *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s: ", label)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func parseIngredient(line string) recipe.Ingredient {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return recipe.Ingredient{Name: line}
	}

	// Try to parse quantity
	quantity, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return recipe.Ingredient{Name: line}
	}

	if len(parts) == 1 {
		return recipe.Ingredient{Quantity: quantity}
	}

	// Check if second part is a unit
	units := map[string]bool{
		"cup": true, "cups": true,
		"tbsp": true, "tablespoon": true, "tablespoons": true,
		"tsp": true, "teaspoon": true, "teaspoons": true,
		"oz": true, "ounce": true, "ounces": true,
		"lb": true, "pound": true, "pounds": true,
		"g": true, "gram": true, "grams": true,
		"kg": true, "kilogram": true, "kilograms": true,
		"ml": true, "milliliter": true, "milliliters": true,
		"l": true, "liter": true, "liters": true,
		"piece": true, "pieces": true,
	}

	unit := ""
	nameStart := 1
	if units[strings.ToLower(parts[1])] {
		unit = parts[1]
		nameStart = 2
	}

	name := strings.Join(parts[nameStart:], " ")

	return recipe.Ingredient{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
	}
}
