package cmd

import (
	"errors"
	"fmt"

	"recipe_box/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify Recipe Box configuration settings.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		value, err := cfg.Get(args[0])
		if err != nil {
			if errors.Is(err, config.ErrKeyNotFound) {
				return fmt.Errorf("key %q not found", args[0])
			}
			return err
		}

		fmt.Println(value)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return err
		}

		if err := cfg.Set(args[0], args[1]); err != nil {
			return err
		}

		fmt.Printf("Set %s = %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
