package cmd

import (
	"recipe_box/internal/i18n"

	"github.com/spf13/cobra"
)

var recipeCmd = &cobra.Command{
	Use: "recipe",
}

func init() {
	recipeCmd.Short = i18n.T(i18n.MsgCmdRecipeShort)
	recipeCmd.Long = i18n.T(i18n.MsgCmdRecipeLong)
	rootCmd.AddCommand(recipeCmd)
}
