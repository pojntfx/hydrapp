package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	secretsFlag = "secrets"
)

var rootCmd = &cobra.Command{
	Use:   "hydrapp",
	Short: "Build apps that run everywhere with Go",
	Long: `Build apps that run everywhere with Go and a browser engine of your choice (Chrome, Firefox, Epiphany or Android WebView).
Find more information at:
https://github.com/pojntfx/hydrapp`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("hydrapp")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		return nil
	},
}

func Execute() error {
	dataHomeDir := os.Getenv("XDG_DATA_HOME")
	if strings.TrimSpace(dataHomeDir) == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		dataHomeDir = filepath.Join(userHomeDir, ".local", "share")
	}

	rootCmd.PersistentFlags().String(secretsFlag, filepath.Join(dataHomeDir, "hydrapp", "secrets.yaml"), "Secrets file to use")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		return err
	}

	viper.AutomaticEnv()

	return rootCmd.Execute()
}
