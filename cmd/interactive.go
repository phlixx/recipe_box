package cmd

import (
	"fmt"
	"os"
	"strings"

	"recipe_box/internal/i18n"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/c-bata/go-prompt"
)

// RunInteractive starts the interactive REPL mode
func RunInteractive() {
	// Set interactive mode flag for commands to use
	IsInteractive = true

	ui.TitlePrintf("Recipe Box")
	fmt.Printf(" - %s\n", i18n.T(i18n.MsgInteractiveMode))
	ui.LabelPrintf("%s\n\n", i18n.T(i18n.MsgInteractiveHint))

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
		fmt.Println(i18n.T(i18n.MsgGoodbye))
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
		case "plan":
			return prompt.FilterHasPrefix(planSuggestions(), words[1], true)
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
		{Text: "recipe", Description: i18n.T(i18n.MsgCmdRecipeShort)},
		{Text: "shopping", Description: i18n.T(i18n.MsgCmdShoppingShort)},
		{Text: "plan", Description: i18n.T(i18n.MsgCmdPlanShort)},
		{Text: "config", Description: i18n.T(i18n.MsgCmdConfigShort)},
		{Text: "help", Description: i18n.T(i18n.MsgCmdHelpShort)},
		{Text: "exit", Description: i18n.T(i18n.MsgCmdExitShort)},
	}
}

func recipeSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "list", Description: i18n.T(i18n.MsgCmdRecipeListShort)},
		{Text: "view", Description: i18n.T(i18n.MsgCmdRecipeViewShort)},
		{Text: "add", Description: i18n.T(i18n.MsgCmdRecipeAddShort)},
		{Text: "generate", Description: i18n.T(i18n.MsgCmdRecipeGenerateShort)},
		{Text: "save", Description: i18n.T(i18n.MsgCmdRecipeSaveShort)},
		{Text: "delete", Description: i18n.T(i18n.MsgCmdRecipeDeleteShort)},
	}
}

func shoppingSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "generate", Description: i18n.T(i18n.MsgCmdShoppingGenerateShort)},
		{Text: "show", Description: i18n.T(i18n.MsgCmdShoppingShowShort)},
		{Text: "clear", Description: i18n.T(i18n.MsgCmdShoppingClearShort)},
	}
}

func planSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "create", Description: i18n.T(i18n.MsgCmdPlanCreateShort)},
		{Text: "show", Description: i18n.T(i18n.MsgCmdPlanShowShort)},
		{Text: "clear", Description: i18n.T(i18n.MsgCmdPlanClearShort)},
	}
}

func configSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "get", Description: i18n.T(i18n.MsgCmdConfigGetShort)},
		{Text: "set", Description: i18n.T(i18n.MsgCmdConfigSetShort)},
	}
}

func configKeySuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "api_key", Description: i18n.T(i18n.MsgConfigKeyAPIKey)},
		{Text: "language", Description: i18n.T(i18n.MsgConfigKeyLanguage)},
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
	ui.TitlePrintf("%s\n\n", i18n.T(i18n.MsgAvailableCmd))

	ui.SectionPrintf("recipe\n")
	fmt.Printf("  list              %s\n", i18n.T(i18n.MsgCmdRecipeListShort))
	fmt.Printf("  view <id>         %s\n", i18n.T(i18n.MsgCmdRecipeViewShort))
	fmt.Printf("  add               %s\n", i18n.T(i18n.MsgCmdRecipeAddShort))
	fmt.Printf("  generate          %s\n", i18n.T(i18n.MsgCmdRecipeGenerateShort))
	fmt.Printf("  save              %s\n", i18n.T(i18n.MsgCmdRecipeSaveShort))
	fmt.Printf("  delete <id>       %s\n", i18n.T(i18n.MsgCmdRecipeDeleteShort))

	fmt.Println()
	ui.SectionPrintf("shopping\n")
	fmt.Printf("  generate <id>     %s\n", i18n.T(i18n.MsgCmdShoppingGenerateShort))
	fmt.Printf("  show              %s\n", i18n.T(i18n.MsgCmdShoppingShowShort))
	fmt.Printf("  clear             %s\n", i18n.T(i18n.MsgCmdShoppingClearShort))

	fmt.Println()
	ui.SectionPrintf("plan\n")
	fmt.Printf("  create            %s\n", i18n.T(i18n.MsgCmdPlanCreateShort))
	fmt.Printf("  show              %s\n", i18n.T(i18n.MsgCmdPlanShowShort))
	fmt.Printf("  clear             %s\n", i18n.T(i18n.MsgCmdPlanClearShort))

	fmt.Println()
	ui.SectionPrintf("config\n")
	fmt.Printf("  get <key>         %s\n", i18n.T(i18n.MsgCmdConfigGetShort))
	fmt.Printf("  set <key> <value> %s\n", i18n.T(i18n.MsgCmdConfigSetShort))

	fmt.Println()
	ui.SectionPrintf("%s\n", i18n.T(i18n.MsgOther))
	fmt.Printf("  help              %s\n", i18n.T(i18n.MsgCmdHelpShort))
	fmt.Printf("  exit              %s\n", i18n.T(i18n.MsgCmdExitShort))
	fmt.Println()
}
