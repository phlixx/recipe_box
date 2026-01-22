package cmd

import (
	"errors"
	"fmt"

	"recipe_box/internal/config"
	"recipe_box/internal/i18n"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
}

var configGetCmd = &cobra.Command{
	Use:  "get <key>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		value, err := cfg.Get(args[0])
		if err != nil {
			if errors.Is(err, config.ErrKeyNotFound) {
				return fmt.Errorf(i18n.T(i18n.MsgConfigNotFound), args[0])
			}
			return err
		}

		fmt.Println(value)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:  "set <key> <value>",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		if err := cfg.Set(args[0], args[1]); err != nil {
			return err
		}

		ui.SuccessPrintf(i18n.T(i18n.MsgConfigSet)+"\n", args[0], args[1])
		return nil
	},
}

func init() {
	configCmd.Short = i18n.T(i18n.MsgCmdConfigShort)
	configCmd.Long = i18n.T(i18n.MsgCmdConfigLong)
	configGetCmd.Short = i18n.T(i18n.MsgCmdConfigGetShort)
	configSetCmd.Short = i18n.T(i18n.MsgCmdConfigSetShort)
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
