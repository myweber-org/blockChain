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

	log.Printf(
		"Method: %s | Path: %s | Duration: %v | User-Agent: %s",
		r.Method,
		r.URL.Path,
		duration,
		r.UserAgent(),
	)
}
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

var activityChannel = make(chan ActivityLog, 100)

func init() {
	go processActivityLogs()
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := extractUserID(r)
		ip := r.RemoteAddr

		activity := ActivityLog{
			UserID:    userID,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: ip,
		}

		select {
		case activityChannel <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return parseToken(auth)
	}
	return "anonymous"
}

func parseToken(token string) string {
	if len(token) > 10 {
		return token[:8] + "..."
	}
	return token
}

func processActivityLogs() {
	for activity := range activityChannel {
		log.Printf("ACTIVITY: User=%s %s %s from %s at %v",
			activity.UserID,
			activity.Method,
			activity.Endpoint,
			activity.IPAddress,
			activity.Timestamp.Format(time.RFC3339))
	}
}