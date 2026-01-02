package middleware

import (
	"encoding/json"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	Status    int       `json:"status"`
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(lw, r)
		
		logEntry := ActivityLog{
			Timestamp: start,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
			Status:    lw.statusCode,
		}
		
		logData, _ := json.Marshal(logEntry)
		println(string(logData))
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