package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey []byte
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{
		secretKey: []byte(secretKey),
	}
}

func (a *Authenticator) ValidateToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, nil
	}

	return validateTokenParts(parts, a.secretKey)
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := a.ValidateToken(tokenString)
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

func validateTokenParts(parts []string, secretKey []byte) (bool, error) {
	// Simplified validation logic for demonstration
	// In production, use proper JWT library like github.com/golang-jwt/jwt
	if len(parts[0]) > 0 && len(parts[1]) > 0 && len(parts[2]) > 0 {
		return true, nil
	}
	return false, nil
}