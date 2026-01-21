package cmd

import (
	"github.com/spf13/cobra"
)

var shoppingCmd = &cobra.Command{
	Use:   "shopping",
	Short: "Manage shopping lists",
	Long:  `Generate, view, and clear shopping lists from your recipes.`,
}

func init() {
	rootCmd.AddCommand(shoppingCmd)
}
