package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "recipe_box",
	Short: "A CLI for managing recipes and shopping lists",
	Long: `Recipe Box is a command-line tool for managing your personal recipe collection.

You can:
  - Add, list, and view recipes
  - Generate shopping lists from selected recipes
  - Organize recipes by tags and categories`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
