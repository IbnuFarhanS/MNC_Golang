package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

// AuthMiddleware adalah middleware untuk mengautentikasi permintaan
func AuthMiddleware(repo repository.CustomerRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Mengautentikasi permintaan...")

			// Ekstrak token dari header
			tokenString, err := extractToken(r)
			if err != nil {
				log.Println("Gagal mengekstrak token:", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Verifikasi token
			token, err := utils.VerifyToken(tokenString)
			if err != nil {
				log.Println("Gagal memverifikasi token:", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Konversi token.Claims ke jwt.MapClaims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				log.Println("Token tidak valid")
				http.Error(w, "token tidak valid", http.StatusUnauthorized)
				return
			}

			// Dapatkan ID pengguna dari klaim token
			userID, ok := claims["user_id"].(string)
			if !ok {
				log.Println("Gagal mengekstrak ID pengguna dari klaim token")
				http.Error(w, "klaim token tidak valid", http.StatusUnauthorized)
				return
			}

			log.Println("ID pengguna terautentikasi:", userID)

			// Periksa apakah token masuk dalam daftar token yang diblacklist
			isBlacklisted, err := repo.IsTokenBlacklisted(tokenString)
			if err != nil {
				log.Println("Gagal memeriksa apakah token masuk dalam daftar token yang diblacklist:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if isBlacklisted {
				log.Println("Token masuk dalam daftar token yang diblacklist")
				http.Error(w, "token masuk dalam daftar token yang diblacklist", http.StatusUnauthorized)
				return
			}

			// Tambahkan ID pengguna ke konteks permintaan
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			r = r.WithContext(ctx)

			log.Println("Permintaan terautentikasi.")
			// Panggil handler berikutnya
			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("header Authorization tidak ada")
	}

	// Format token: Bearer <token>
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("format token tidak valid")
	}

	// Hapus awalan "Bearer " dari token, jika ada
	token := strings.TrimPrefix(tokenParts[1], "Bearer ")

	return token, nil
}
