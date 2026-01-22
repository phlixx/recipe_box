package cmd

import (
	"recipe_box/internal/i18n"

	"github.com/spf13/cobra"
)

var shoppingCmd = &cobra.Command{
	Use: "shopping",
}

func init() {
	shoppingCmd.Short = i18n.T(i18n.MsgCmdShoppingShort)
	shoppingCmd.Long = i18n.T(i18n.MsgCmdShoppingLong)
	rootCmd.AddCommand(shoppingCmd)
}
