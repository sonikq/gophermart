package generator

import (
	"crypto/rand"
	"math/big"
)

const (
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	specialChars = "!@#$%^&*"
)

func GeneratePassword(length int) (string, error) {
	charset := lowerLetters + upperLetters + digits + specialChars

	var password []byte
	for i := 0; i < length; i++ {
		char, err := secureRandChar(charset)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	return string(password), nil
}

func secureRandChar(charset string) (byte, error) {
	maxBigInt := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, maxBigInt)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}
