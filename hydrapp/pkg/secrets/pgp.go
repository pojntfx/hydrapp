package secrets

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func GeneratePGPKey(
	fullName,
	email,
	password string,
) (string, string, error) {
	rawPGPKey, err := crypto.GenerateKey(
		fullName,
		email,
		"x25519",
		0,
	)
	if err != nil {
		return "", "", err
	}
	defer rawPGPKey.ClearPrivateParams()

	lockedPGPKey, err := rawPGPKey.Lock([]byte(password))
	if err != nil {
		return "", "", err
	}

	armoredPGPKey, err := lockedPGPKey.Armor()
	if err != nil {
		return "", "", err
	}

	return armoredPGPKey, rawPGPKey.GetHexKeyID(), nil
}
