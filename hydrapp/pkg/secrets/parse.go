package secrets

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Root struct {
	JavaSecrets JavaSecrets `yaml:"java"`
	PGPSecrets  PGPSecrets  `yaml:"pgp"`
}

type JavaSecrets struct {
	Keystore            []byte `yaml:"keystore"`
	KeystorePassword    string `yaml:"keystorePassword"`
	CertificatePassword string `yaml:"certificatePassword"`
}

type PGPSecrets struct {
	Key         string `yaml:"key"`
	KeyID       string `yaml:"keyID"`
	KeyPassword string `yaml:"keyPassword"`
}

func Parse(r io.Reader) (*Root, error) {
	var root Root
	if err := yaml.NewDecoder(r).Decode(&root); err != nil {
		return nil, err
	}

	return &root, nil
}
