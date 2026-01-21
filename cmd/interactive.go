package cmd

import (
	"fmt"
	"os"
	"strings"

	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/c-bata/go-prompt"
)

// RunInteractive starts the interactive REPL mode
func RunInteractive() {
	ui.TitlePrintf("Recipe Box")
	fmt.Println(" - Interactive Mode")
	ui.LabelPrintf("Type 'help' for commands, 'exit' to quit\n\n")

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("\U0001F373 > "), // frying pan emoji
		prompt.OptionTitle("Recipe Box"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
	p.Run()
}

// executor processes the user input
func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Handle exit command
	if input == "exit" || input == "quit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	// Handle help command
	if input == "help" {
		printHelp()
		return
	}

	// Parse and execute as Cobra command
	args := strings.Fields(input)
	rootCmd.SetArgs(args)

	// Capture and display any errors
	if err := rootCmd.Execute(); err != nil {
		ui.ErrorPrintf("Error: %s\n", err)
	}

	// Reset args for next command
	rootCmd.SetArgs(nil)
	fmt.Println()
}

// completer provides tab completion suggestions
func completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	words := strings.Fields(text)

	// If cursor is right after a space, add an empty word to indicate we're completing next arg
	if len(text) > 0 && text[len(text)-1] == ' ' {
		words = append(words, "")
	}

	if len(words) == 0 {
		return topLevelSuggestions()
	}

	// First word: suggest top-level commands
	if len(words) == 1 {
		return prompt.FilterHasPrefix(topLevelSuggestions(), words[0], true)
	}

	// Second word: suggest subcommands based on first word
	if len(words) == 2 {
		switch words[0] {
		case "recipe":
			return prompt.FilterHasPrefix(recipeSuggestions(), words[1], true)
		case "shopping":
			return prompt.FilterHasPrefix(shoppingSuggestions(), words[1], true)
		case "config":
			return prompt.FilterHasPrefix(configSuggestions(), words[1], true)
		}
	}

	// Third word: suggest recipe IDs for commands that need them
	if len(words) == 3 {
		switch words[0] {
		case "recipe":
			if words[1] == "view" || words[1] == "delete" {
				return prompt.FilterHasPrefix(recipeIDSuggestions(), words[2], true)
			}
		case "shopping":
			if words[1] == "generate" {
				return prompt.FilterHasPrefix(recipeIDSuggestions(), words[2], true)
			}
		case "config":
			if words[1] == "get" || words[1] == "set" {
				return prompt.FilterHasPrefix(configKeySuggestions(), words[2], true)
			}
		}
	}

	return nil
}

func topLevelSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "recipe", Description: "Manage recipes"},
		{Text: "shopping", Description: "Manage shopping lists"},
		{Text: "config", Description: "Manage configuration"},
		{Text: "help", Description: "Show available commands"},
		{Text: "exit", Description: "Exit interactive mode"},
	}
}

func recipeSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "list", Description: "List all recipes"},
		{Text: "view", Description: "View a recipe"},
		{Text: "add", Description: "Add a new recipe manually"},
		{Text: "generate", Description: "Generate a recipe with AI"},
		{Text: "save", Description: "Save the last generated recipe"},
		{Text: "delete", Description: "Delete a recipe"},
	}
}

func shoppingSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "generate", Description: "Generate shopping list from recipe"},
		{Text: "show", Description: "Show current shopping list"},
		{Text: "clear", Description: "Clear the shopping list"},
	}
}

func configSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "get", Description: "Get a configuration value"},
		{Text: "set", Description: "Set a configuration value"},
	}
}

func configKeySuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "api_key", Description: "Claude API key"},
		{Text: "language", Description: "Language (en/de)"},
	}
}

// recipeIDSuggestions returns suggestions for recipe IDs
func recipeIDSuggestions() []prompt.Suggest {
	store, err := storage.New()
	if err != nil {
		return nil
	}

	svc := recipe.NewService(store)
	recipes, err := svc.List()
	if err != nil {
		return nil
	}

	suggestions := make([]prompt.Suggest, 0, len(recipes))
	for _, r := range recipes {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        r.ID,
			Description: r.Title,
		})
	}

	return suggestions
}

func printHelp() {
	ui.TitlePrintf("Available Commands:\n\n")

	ui.SectionPrintf("recipe\n")
	fmt.Println("  list              List all recipes")
	fmt.Println("  view <id>         View a recipe")
	fmt.Println("  add               Add a new recipe manually")
	fmt.Println("  generate          Generate a recipe with AI")
	fmt.Println("  save              Save the last generated recipe")
	fmt.Println("  delete <id>       Delete a recipe")

	fmt.Println()
	ui.SectionPrintf("shopping\n")
	fmt.Println("  generate <id>     Generate shopping list from recipe")
	fmt.Println("  show              Show current shopping list")
	fmt.Println("  clear             Clear the shopping list")

	fmt.Println()
	ui.SectionPrintf("config\n")
	fmt.Println("  get <key>         Get a configuration value")
	fmt.Println("  set <key> <value> Set a configuration value")

	fmt.Println()
	ui.SectionPrintf("Other\n")
	fmt.Println("  help              Show this help message")
	fmt.Println("  exit              Exit interactive mode")
	fmt.Println()
}
