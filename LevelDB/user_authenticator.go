package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{secretKey: secretKey}
}

func (am *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		if !am.isValidToken(tokenString) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) isValidToken(token string) bool {
	// Simplified token validation logic
	// In production, use proper JWT library like github.com/golang-jwt/jwt
	return token != "" && len(token) > 10 && strings.HasPrefix(token, "valid_")
}

func (am *AuthMiddleware) GenerateToken(userID string) (string, error) {
	// Simplified token generation
	// In production, use proper JWT signing
	if userID == "" {
		return "", fmt.Errorf("userID cannot be empty")
	}
	return "valid_" + userID + "_token", nil
}