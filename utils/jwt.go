package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken digunakan untuk menghasilkan token JWT dengan menggunakan user ID dan role yang diberikan
func GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Set expiration time, e.g., 1 Jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Set JWT signing key
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken digunakan untuk memverifikasi keabsahan token JWT yang diberikan.
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	// Verify token validity
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
