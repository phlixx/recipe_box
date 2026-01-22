package cmd

import (
	"errors"
	"fmt"
	"strings"

	"recipe_box/internal/i18n"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var viewServings int

var recipeViewCmd = &cobra.Command{
	Use:  "view <id>",
	Args: cobra.ExactArgs(1),
	RunE: runRecipeView,
}

func init() {
	recipeViewCmd.Short = i18n.T(i18n.MsgCmdRecipeViewShort)
	recipeViewCmd.Long = i18n.T(i18n.MsgCmdRecipeViewLong)
	recipeCmd.AddCommand(recipeViewCmd)
	recipeViewCmd.Flags().IntVarP(&viewServings, "servings", "s", 0, "Scale recipe to this number of servings")
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
			return fmt.Errorf(i18n.T(i18n.MsgRecipeNotFound), id)
		}
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	// Scale if --servings flag is provided
	originalServings := 0
	if viewServings > 0 && viewServings != r.Servings {
		originalServings = r.Servings
		r = r.Scale(viewServings)
	}

	printRecipeWithScale(r, originalServings)
	return nil
}

func printRecipe(r *recipe.Recipe) {
	printRecipeWithScale(r, 0)
}

func printRecipeWithScale(r *recipe.Recipe, originalServings int) {
	// Display ASCII art if available
	if r.AsciiArt != "" {
		fmt.Println()
		fmt.Println(r.AsciiArt)
		fmt.Println()
	}

	ui.TitlePrintf("%s\n", r.Title)
	fmt.Println(strings.Repeat("=", len(r.Title)))
	fmt.Println()

	if r.Description != "" {
		fmt.Printf("%s\n\n", r.Description)
	}

	// Metadata
	ui.LabelPrintf("%s: ", i18n.T(i18n.MsgLabelServings))
	if originalServings > 0 {
		fmt.Printf("%d %s\n", r.Servings, fmt.Sprintf(i18n.T(i18n.MsgScaledFrom), originalServings))
	} else {
		fmt.Printf("%d\n", r.Servings)
	}
	if r.PrepTime > 0 {
		ui.LabelPrintf("%s: ", i18n.T(i18n.MsgLabelPrepTime))
		fmt.Printf("%d %s\n", r.PrepTime, i18n.T(i18n.MsgMinutes))
	}
	if r.CookTime > 0 {
		ui.LabelPrintf("%s: ", i18n.T(i18n.MsgLabelCookTime))
		fmt.Printf("%d %s\n", r.CookTime, i18n.T(i18n.MsgMinutes))
	}
	if r.PrepTime > 0 || r.CookTime > 0 {
		ui.LabelPrintf("%s: ", i18n.T(i18n.MsgLabelTotalTime))
		fmt.Printf("%d %s\n", r.PrepTime+r.CookTime, i18n.T(i18n.MsgMinutes))
	}

	// Tags
	if len(r.Tags) > 0 {
		ui.LabelPrintf("%s: ", i18n.T(i18n.MsgLabelTags))
		fmt.Printf("%s\n", strings.Join(r.Tags, ", "))
	}

	// Ingredients
	fmt.Println()
	ui.SectionPrintf("%s:\n", i18n.T(i18n.MsgLabelIngredients))
	for _, ing := range r.Ingredients {
		if ing.Quantity > 0 && ing.Unit != "" {
			fmt.Printf("  - %.4g %s %s\n", ing.Quantity, i18n.Unit(ing.Unit), ing.Name)
		} else if ing.Quantity > 0 {
			fmt.Printf("  - %.4g %s\n", ing.Quantity, ing.Name)
		} else {
			fmt.Printf("  - %s\n", ing.Name)
		}
	}

	// Steps
	fmt.Println()
	ui.SectionPrintf("%s:\n", i18n.T(i18n.MsgLabelSteps))
	for i, step := range r.Steps {
		fmt.Printf("  %d. %s\n", i+1, step)
	}

	fmt.Printf("\n[ID: %s]\n", ui.IDSprint(r.ID))
}
