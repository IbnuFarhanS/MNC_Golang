package utils

import "golang.org/x/crypto/bcrypt"

// GenerateHash generates the bcrypt hash for the given password
func GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
