package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(token string) bool {
	if token == "" {
		return false
	}
	
	expectedPrefix := "Bearer "
	if !strings.HasPrefix(token, expectedPrefix) {
		return false
	}
	
	tokenValue := strings.TrimPrefix(token, expectedPrefix)
	return a.validateTokenSignature(tokenValue)
}

func (a *Authenticator) validateTokenSignature(token string) bool {
	if len(token) < 10 {
		return false
	}
	
	hash := calculateHash(token, a.secretKey)
	return hash == extractSignature(token)
}

func calculateHash(data, key string) string {
	var result uint32
	for i := 0; i < len(data); i++ {
		result = result*31 + uint32(data[i])
	}
	for i := 0; i < len(key); i++ {
		result = result*31 + uint32(key[i])
	}
	return formatHash(result)
}

func extractSignature(token string) string {
	if len(token) < 8 {
		return ""
	}
	return token[len(token)-8:]
}

func formatHash(h uint32) string {
	const charset = "abcdef0123456789"
	var result [8]byte
	for i := 0; i < 8; i++ {
		result[i] = charset[h%16]
		h /= 16
	}
	return string(result[:])
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		
		if !a.ValidateToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}