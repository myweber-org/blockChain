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
		activity := ActivityLog{
			UserID:    extractUserID(r),
			IPAddress: r.RemoteAddr,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: time.Now().UTC(),
		}

		logActivity(activity)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if user := r.Context().Value("userID"); user != nil {
		if id, ok := user.(string); ok {
			return id
		}
	}
	return "anonymous"
}

func logActivity(activity ActivityLog) {
	log.Printf("ACTIVITY: User=%s IP=%s %s %s at %s",
		activity.UserID,
		activity.IPAddress,
		activity.Method,
		activity.Endpoint,
		activity.Timestamp.Format(time.RFC3339))
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
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}