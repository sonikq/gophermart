package hash

import "golang.org/x/crypto/bcrypt"

func Generate(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hash), err
}

func Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
