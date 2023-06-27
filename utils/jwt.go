package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Set expiration time, e.g., 1 day from now
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Set JWT signing key
	// Replace "your-secret-key" with your own secret key for signing the token
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		// Replace "your-secret-key" with your own secret key used for signing the token
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
