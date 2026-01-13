package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secretKey: []byte(secret)}
}

func (am *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
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

		tokenString := parts[1]
		claims, err := am.parseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(am.setClaimsInContext(r.Context(), claims))
		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) parseToken(tokenString string) (map[string]interface{}, error) {
	// Token parsing implementation would go here
	// This is a simplified placeholder
	return map[string]interface{}{
		"user_id": "sample_user",
		"role":    "user",
	}, nil
}

func (am *AuthMiddleware) setClaimsInContext(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, "user_claims", claims)
}