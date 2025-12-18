package middleware

import (
	"net/http"
	"strings"
)

type User struct {
	ID    string
	Email string
	Role  string
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		user, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) (*User, error) {
	// This is a placeholder for actual token validation logic
	// In production, this would verify JWT signature and claims
	if token == "valid_token_example" {
		return &User{
			ID:    "user123",
			Email: "user@example.com",
			Role:  "admin",
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}