package middleware

import (
	"net/http"
	"strings"
)

type User struct {
	ID       int
	Username string
	Role     string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		user, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid authentication token", http.StatusForbidden)
			return
		}

		if !hasRequiredRole(user, r.URL.Path) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func validateToken(token string) (*User, error) {
	// Token validation logic would go here
	// This is a simplified example
	if token == "valid-admin-token" {
		return &User{ID: 1, Username: "admin", Role: "admin"}, nil
	}
	if token == "valid-user-token" {
		return &User{ID: 2, Username: "user", Role: "user"}, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func hasRequiredRole(user *User, path string) bool {
	if strings.HasPrefix(path, "/admin") && user.Role != "admin" {
		return false
	}
	return true
}