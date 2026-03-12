package middleware

import (
    "net/http"
    "strings"
)

type AuthMiddleware struct {
    sessionStore SessionStore
}

func NewAuthMiddleware(store SessionStore) *AuthMiddleware {
    return &AuthMiddleware{sessionStore: store}
}

func (m *AuthMiddleware) ValidateSession(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }

        token := parts[1]
        userID, valid := m.sessionStore.ValidateToken(token)
        if !valid {
            http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-User-ID", userID)
        next.ServeHTTP(w, r)
    })
}

type SessionStore interface {
    ValidateToken(token string) (userID string, valid bool)
}