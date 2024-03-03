package passwords

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("bycrypt.GenerateFromPassword: %w", err)
	}
	return string(enc), nil
}

func Compare(password string, enc string) bool {
	return bcrypt.CompareHashAndPassword([]byte(enc), []byte(password)) == nil
}
