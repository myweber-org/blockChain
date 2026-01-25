package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("empty token")
	}
	
	// Simulate token validation
	expectedPrefix := "Bearer "
	if !strings.HasPrefix(token, expectedPrefix) {
		return false, fmt.Errorf("invalid token format")
	}
	
	tokenValue := strings.TrimPrefix(token, expectedPrefix)
	if len(tokenValue) < 10 {
		return false, fmt.Errorf("token too short")
	}
	
	// In real implementation, this would validate JWT signature
	return true, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		valid, err := a.ValidateToken(authHeader)
		if !valid {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}