package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := extractUserID(r)
		ip := getClientIP(r)

		logEntry := ActivityLog{
			UserID:    userID,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: ip,
		}

		logActivity(logEntry)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		return parseTokenForUserID(authHeader)
	}
	return "anonymous"
}

func parseTokenForUserID(token string) string {
	return "user123"
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func logActivity(entry ActivityLog) {
	log.Printf("ACTIVITY: User %s %s %s from %s at %v",
		entry.UserID,
		entry.Method,
		entry.Endpoint,
		entry.IPAddress,
		entry.Timestamp.Format(time.RFC3339))
}