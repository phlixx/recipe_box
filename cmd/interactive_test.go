package cmd

import (
	"testing"
)

// TestPlanSuggestionsComplete verifies that planSuggestions() contains all plan subcommands.
// This test ensures that when new commands are added, they are also added to the interactive
// mode autocompletion.
func TestPlanSuggestionsComplete(t *testing.T) {
	// Get actual subcommands registered under planCmd
	registeredCmds := make(map[string]bool)
	for _, cmd := range planCmd.Commands() {
		registeredCmds[cmd.Name()] = true
	}

	// Get suggestions from planSuggestions()
	suggestions := planSuggestions()
	suggestedCmds := make(map[string]bool)
	for _, s := range suggestions {
		suggestedCmds[s.Text] = true
	}

	// Check that every registered command has a suggestion
	for cmdName := range registeredCmds {
		if !suggestedCmds[cmdName] {
			t.Errorf("plan subcommand %q is registered but missing from planSuggestions()", cmdName)
		}
	}
}

// TestRecipeSuggestionsComplete verifies that recipeSuggestions() contains all recipe subcommands.
func TestRecipeSuggestionsComplete(t *testing.T) {
	registeredCmds := make(map[string]bool)
	for _, cmd := range recipeCmd.Commands() {
		registeredCmds[cmd.Name()] = true
	}

	suggestions := recipeSuggestions()
	suggestedCmds := make(map[string]bool)
	for _, s := range suggestions {
		suggestedCmds[s.Text] = true
	}

	for cmdName := range registeredCmds {
		if !suggestedCmds[cmdName] {
			t.Errorf("recipe subcommand %q is registered but missing from recipeSuggestions()", cmdName)
		}
	}
}

// TestShoppingSuggestionsComplete verifies that shoppingSuggestions() contains all shopping subcommands.
func TestShoppingSuggestionsComplete(t *testing.T) {
	registeredCmds := make(map[string]bool)
	for _, cmd := range shoppingCmd.Commands() {
		registeredCmds[cmd.Name()] = true
	}

	suggestions := shoppingSuggestions()
	suggestedCmds := make(map[string]bool)
	for _, s := range suggestions {
		suggestedCmds[s.Text] = true
	}

	for cmdName := range registeredCmds {
		if !suggestedCmds[cmdName] {
			t.Errorf("shopping subcommand %q is registered but missing from shoppingSuggestions()", cmdName)
		}
	}
}

// TestConfigSuggestionsComplete verifies that configSuggestions() contains all config subcommands.
func TestConfigSuggestionsComplete(t *testing.T) {
	registeredCmds := make(map[string]bool)
	for _, cmd := range configCmd.Commands() {
		registeredCmds[cmd.Name()] = true
	}

	suggestions := configSuggestions()
	suggestedCmds := make(map[string]bool)
	for _, s := range suggestions {
		suggestedCmds[s.Text] = true
	}

	for cmdName := range registeredCmds {
		if !suggestedCmds[cmdName] {
			t.Errorf("config subcommand %q is registered but missing from configSuggestions()", cmdName)
		}
	}
}
