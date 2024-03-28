package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show hydrapp secrets as env variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		secretsFile, err := os.Open(viper.GetString(secretsFlag))
		if err != nil {
			return err
		}
		defer secretsFile.Close()

		scs, err := secrets.Parse(secretsFile)
		if err != nil {
			return err
		}

		fmt.Printf(`export JAVA_KEYSTORE="%v"
export JAVA_KEYSTORE_PASSWORD="%v"
export JAVA_CERTIFICATE_PASSWORD="%v"
export PGP_KEY="%v"
export PGP_KEY_ID="%v"
export PGP_KEY_PASSWORD="%v"
`,

			base64.StdEncoding.EncodeToString(scs.JavaSecrets.Keystore),
			base64.StdEncoding.EncodeToString([]byte(scs.JavaSecrets.KeystorePassword)),
			base64.StdEncoding.EncodeToString([]byte(scs.JavaSecrets.CertificatePassword)),

			base64.StdEncoding.EncodeToString([]byte(scs.PGPSecrets.Key)),
			base64.StdEncoding.EncodeToString([]byte(scs.PGPSecrets.KeyID)),
			base64.StdEncoding.EncodeToString([]byte(scs.PGPSecrets.KeyPassword)),
		)

		return nil
	},
}

func init() {
	viper.AutomaticEnv()

	secretsCmd.AddCommand(secretsShowCmd)
}
