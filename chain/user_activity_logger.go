package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IP        string
	Duration  time.Duration
	Status    int
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)
		
		activity := ActivityLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
			Duration:  time.Since(start),
			Status:    lw.statusCode,
		}
		
		log.Printf("ACTIVITY: %s %s %d %s %v",
			activity.Method,
			activity.Path,
			activity.Status,
			activity.IP,
			activity.Duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}