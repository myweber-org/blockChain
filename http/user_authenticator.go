package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: []byte(secret),
	}
}

func (am *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := validateJWTToken(tokenString, am.secretKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func validateJWTToken(tokenString string, secretKey []byte) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.Exp < time.Now().Unix() {
			return nil, fmt.Errorf("Token expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token claims")
}package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("your-secret-key-change-in-production")

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "myapp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}