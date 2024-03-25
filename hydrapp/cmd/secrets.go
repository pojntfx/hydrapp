package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"s"},
	Short:   "Manage secrets",
}

func init() {
	viper.AutomaticEnv()

	rootCmd.AddCommand(secretsCmd)
}
