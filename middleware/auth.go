package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/IbnuFarhanS/Golang_MNC/utils"
	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

// AuthMiddleware is a middleware to authenticate the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Authenticating request...")

		// Extract token from header
		tokenString, err := extractToken(r)
		if err != nil {
			log.Println("Failed to extract token:", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify token
		token, err := utils.VerifyToken(tokenString)
		if err != nil {
			log.Println("Failed to verify token:", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Convert token.Claims to jwt.MapClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("Invalid token")
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Get user ID from token claim
		userID, ok := claims["user_id"].(string)
		if !ok {
			log.Println("Failed to extract user ID from token claims")
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		log.Println("Authenticated user ID:", userID)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		log.Println("Request authenticated.")
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	// Token format: Bearer <token>
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	return tokenParts[1], nil
}
