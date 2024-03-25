package utils

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword(length int) (string, error) {
	password := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		password += string(alphabet[num.Int64()])
	}

	return password, nil
}
