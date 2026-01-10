package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userID := extractUserID(r)
	ip := r.RemoteAddr
	method := r.Method
	path := r.URL.Path

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from %s %s %s completed in %v", userID, ip, method, path, duration)
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return auth[:min(8, len(auth))]
	}
	return "anonymous"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	IPAddress string
	Endpoint  string
	Method    string
	Timestamp time.Time
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := extractUserID(r)
		ip := r.RemoteAddr
		endpoint := r.URL.Path
		method := r.Method

		logEntry := ActivityLog{
			UserID:    userID,
			IPAddress: ip,
			Endpoint:  endpoint,
			Method:    method,
			Timestamp: start,
		}

		logActivity(logEntry)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "anonymous"
	}
	return hashToken(token)
}

func hashToken(token string) string {
	return token[:8]
}

func logActivity(entry ActivityLog) {
	log.Printf("ACTIVITY: User %s from %s accessed %s %s at %v",
		entry.UserID,
		entry.IPAddress,
		entry.Method,
		entry.Endpoint,
		entry.Timestamp.Format(time.RFC3339))
}