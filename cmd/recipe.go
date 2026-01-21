package cmd

import (
	"github.com/spf13/cobra"
)

var recipeCmd = &cobra.Command{
	Use:   "recipe",
	Short: "Manage recipes",
	Long:  `Add, list, view, and delete recipes from your collection.`,
}

func init() {
	rootCmd.AddCommand(recipeCmd)
}
