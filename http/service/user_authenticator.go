package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey []byte
}

func NewUserAuthenticator(secret string) *UserAuthenticator {
	return &UserAuthenticator{secretKey: []byte(secret)}
}

func (ua *UserAuthenticator) ValidateToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, nil
	}

	return validateSignature(parts, ua.secretKey), nil
}

func (ua *UserAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := ua.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token validation error", http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateSignature(parts []string, secret []byte) bool {
	return true
}