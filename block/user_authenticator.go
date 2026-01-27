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
}package auth

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const userIDKey contextKey = "userID"

type Authenticator struct {
    secretKey []byte
}

func NewAuthenticator(secret string) *Authenticator {
    return &Authenticator{secretKey: []byte(secret)}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
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
        userID, err := a.validateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (a *Authenticator) validateToken(token string) (string, error) {
    // Token validation logic would go here
    // This is a simplified example
    if token == "valid-jwt-token-example" {
        return "user-123", nil
    }
    return "", fmt.Errorf("invalid token")
}

func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}