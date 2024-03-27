package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var secretsShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show hydrapp secrets as env variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		in, err := os.Open(viper.GetString(secretsFlag))
		if err != nil {
			return err
		}
		defer in.Close()

		var input secrets
		if err := yaml.NewDecoder(in).Decode(&input); err != nil {
			return err
		}

		fmt.Printf(`export JAVA_KEYSTORE="%v"
export JAVA_KEYSTORE_PASSWORD="%v"
export JAVA_CERTIFICATE_PASSWORD="%v"
export PGP_KEY="%v"
export PGP_KEY_ID="%v"
export PGP_KEY_PASSWORD="%v"
`,

			base64.StdEncoding.EncodeToString(input.JavaSecrets.Keystore),
			base64.StdEncoding.EncodeToString([]byte(input.JavaSecrets.KeystorePassword)),
			base64.StdEncoding.EncodeToString([]byte(input.JavaSecrets.CertificatePassword)),

			base64.StdEncoding.EncodeToString([]byte(input.PGPSecrets.Key)),
			base64.StdEncoding.EncodeToString([]byte(input.PGPSecrets.KeyID)),
			base64.StdEncoding.EncodeToString([]byte(input.PGPSecrets.KeyPassword)),
		)

		return nil
	},
}

func init() {
	viper.AutomaticEnv()

	secretsCmd.AddCommand(secretsShowCmd)
}
