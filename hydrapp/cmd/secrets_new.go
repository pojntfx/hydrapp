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
	JavaSecrets javaSecrets `yaml:"java"`
	PGPSecrets  pgpSecrets  `yaml:"pgp"`
}

type javaSecrets struct {
	Keystore            []byte `yaml:"keystore"`
	KeystorePassword    string `yaml:"keystorePassword"`
	CertificatePassword string `yaml:"certificatePassword"`
}

type pgpSecrets struct {
	Key         string `yaml:"key"`
	KeyID       string `yaml:"keyID"`
	KeyPassword string `yaml:"keyPassword"`
}

const (
	javaKeystorePasswordFlag    = "java-keystore-password"
	javaCertificatePasswordFlag = "java-certificate-password"
	javaCertificateAliasFlag    = "java-certificate-alias"
	javaCertificateCNAMEFlag    = "java-certificate-cname"
	javaCertificateValidityFlag = "java-certificate-validity"
	javaRSABitsFlag             = "java-rsa-bits"

	pgpFullNameFlag = "pgp-full-name"
	pgpEmailFlag    = "pgp-email"
	pgpPasswordFlag = "pgp-password"
)

var secretsNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "Generate new hydrapp secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		keystorePassword := viper.GetString(javaKeystorePasswordFlag)
		if strings.TrimSpace(keystorePassword) == "" {
			v, err := s.GeneratePassword(32)
			if err != nil {
				panic(err)
			}

			keystorePassword = v
		}

		certificatePassword := viper.GetString(javaCertificatePasswordFlag)
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
			viper.GetString(javaCertificateAliasFlag),
			viper.GetString(javaCertificateCNAMEFlag),
			viper.GetDuration(javaCertificateValidityFlag),
			viper.GetUint32(javaRSABitsFlag),
			keystoreBuf,
		); err != nil {
			panic(err)
		}

		pgpPassword := viper.GetString(pgpPasswordFlag)
		if strings.TrimSpace(pgpPassword) == "" {
			v, err := s.GeneratePassword(32)
			if err != nil {
				panic(err)
			}

			pgpPassword = v
		}

		pgpKey, pgpKeyID, err := s.GeneratePGPKey(
			viper.GetString(pgpFullNameFlag),
			viper.GetString(pgpEmailFlag),
			pgpPassword,
		)
		if err != nil {
			panic(err)
		}

		output := &secrets{
			JavaSecrets: javaSecrets{
				Keystore:            keystoreBuf.Bytes(),
				KeystorePassword:    keystorePassword,
				CertificatePassword: certificatePassword,
			},
			PGPSecrets: pgpSecrets{
				Key:         pgpKey,
				KeyID:       pgpKeyID,
				KeyPassword: pgpPassword,
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
	secretsNewCmd.PersistentFlags().String(javaKeystorePasswordFlag, "", "Java/APK keystore password (auto-generated if not specified)")
	secretsNewCmd.PersistentFlags().String(javaCertificatePasswordFlag, "", "Java/APK certificate password (auto-generated if not specified)")
	secretsNewCmd.PersistentFlags().String(javaCertificateAliasFlag, "Anonymous Hydrapp Developer", "Java/APK certificate alias")
	secretsNewCmd.PersistentFlags().String(javaCertificateCNAMEFlag, "Anonymous Hydrapp Developer", "Java/APK certificate CNAME")
	secretsNewCmd.PersistentFlags().Duration(javaCertificateValidityFlag, time.Hour*24*365, "Java/APK certificate validty")
	secretsNewCmd.PersistentFlags().Uint32(javaRSABitsFlag, 2048, "Java/APK RSA bits")

	secretsNewCmd.PersistentFlags().String(pgpFullNameFlag, "Anonymous Hydrapp Developer", "PGP full name")
	secretsNewCmd.PersistentFlags().String(pgpEmailFlag, "test@example.com", "PGP E-Mail")
	secretsNewCmd.PersistentFlags().String(pgpPasswordFlag, "", "PGP password (auto-generated if not specified)")

	viper.AutomaticEnv()

	secretsCmd.AddCommand(secretsNewCmd)
}
