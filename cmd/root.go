package cmd

import (
	"fmt"
	"os"

	"recipe_box/internal/i18n"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "recipe_box",
}

func init() {
	rootCmd.Short = i18n.T(i18n.MsgCmdRootShort)
	rootCmd.Long = i18n.T(i18n.MsgCmdRootLong)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
