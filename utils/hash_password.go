package utils

import "golang.org/x/crypto/bcrypt"

// GenerateHash digunakan untuk menghasilkan hash bcrypt untuk password yang diberikan
func GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
