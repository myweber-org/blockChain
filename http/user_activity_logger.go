
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	IPAddress string
	Method    string
	Path      string
	Timestamp time.Time
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)

		activity := ActivityLog{
			UserID:    extractUserID(r),
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
		}

		logActivity(activity, lw.statusCode, time.Since(start))
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

func logActivity(activity ActivityLog, status int, duration time.Duration) {
	log.Printf("ACTIVITY: user=%s ip=%s method=%s path=%s status=%d duration=%v",
		activity.UserID,
		activity.IPAddress,
		activity.Method,
		activity.Path,
		status,
		duration,
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}