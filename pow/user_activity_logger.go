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
	userID := extractUserID(r)
	ipAddress := r.RemoteAddr
	method := r.Method
	path := r.URL.Path

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from %s performed %s on %s in %v", userID, ipAddress, method, path, duration)
}

func extractUserID(r *http.Request) string {
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		return "authenticated_user"
	}
	return "anonymous_user"
}