package cmd

import (
	"recipe_box/internal/i18n"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use: "plan",
}

func init() {
	planCmd.Short = i18n.T(i18n.MsgCmdPlanShort)
	planCmd.Long = i18n.T(i18n.MsgCmdPlanLong)
	rootCmd.AddCommand(planCmd)
}
