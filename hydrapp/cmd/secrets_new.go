package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	s "github.com/pojntfx/hydrapp/hydrapp/pkg/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type secrets struct {
	APKSecrets apkSecrets `yaml:"apk"`
}

type apkSecrets struct {
	Keystore            []byte `yaml:"keystore"`
	KeystorePassword    string `yaml:"keystorePassword"`
	CertificatePassword string `yaml:"certificatePassword"`
}

const (
	apkKeystorePasswordFlag    = "apk-keystore-password"
	apkCertificatePasswordFlag = "apk-certificate-password"
	apkCertificateAliasFlag    = "apk-certificate-alias"
	apkCertificateCNAMEFlag    = "apk-certificate-cname"
	apkCertificateValidityFlag = "apk-certificate-validity"
	apkRSABitsFlag             = "apk-rsa-bits"
)

var secretsNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "Generate new hydrapp secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		keystorePassword := viper.GetString(apkKeystorePasswordFlag)
		if strings.TrimSpace(keystorePassword) == "" {
			v, err := s.GeneratePassword(32)
			if err != nil {
				panic(err)
			}

			keystorePassword = v
		}

		certificatePassword := viper.GetString(apkCertificatePasswordFlag)
		if strings.TrimSpace(certificatePassword) == "" {
			v, err := s.GeneratePassword(32)
			if err != nil {
				panic(err)
			}

			certificatePassword = v
		}

		keystoreBuf := &bytes.Buffer{}
		if err := s.GenerateKeystore(
			keystorePassword,
			certificatePassword,
			viper.GetString(apkCertificateAliasFlag),
			viper.GetString(apkCertificateCNAMEFlag),
			viper.GetDuration(apkCertificateValidityFlag),
			viper.GetUint32(apkRSABitsFlag),
			keystoreBuf,
		); err != nil {
			panic(err)
		}

		output := &secrets{
			APKSecrets: apkSecrets{
				Keystore:            keystoreBuf.Bytes(),
				KeystorePassword:    keystorePassword,
				CertificatePassword: certificatePassword,
			},
		}

		if err := os.MkdirAll(filepath.Dir(viper.GetString(secretsFlag)), os.ModePerm); err != nil {
			return err
		}

		out, err := os.OpenFile(viper.GetString(secretsFlag), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer out.Close()

		return yaml.NewEncoder(out).Encode(output)
	},
}

func init() {
	secretsNewCmd.PersistentFlags().String(apkKeystorePasswordFlag, "", "APK/Java keystore password (auto-generated if not specified)")
	secretsNewCmd.PersistentFlags().String(apkCertificatePasswordFlag, "", "APK/Java certificate password (auto-generated if not specified)")
	secretsNewCmd.PersistentFlags().String(apkCertificateAliasFlag, "Anonymous Hydrapp Developer", "APK/Java certificate alias")
	secretsNewCmd.PersistentFlags().String(apkCertificateCNAMEFlag, "Anonymous Hydrapp Developer", "APK/Java certificate CNAME")
	secretsNewCmd.PersistentFlags().Duration(apkCertificateValidityFlag, time.Hour*24*365, "APK/Java certificate validty")
	secretsNewCmd.PersistentFlags().Uint32(apkRSABitsFlag, 2048, "APK/Java RSA bits")

	viper.AutomaticEnv()

	secretsCmd.AddCommand(secretsNewCmd)
}
