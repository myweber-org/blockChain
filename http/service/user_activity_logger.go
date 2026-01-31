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
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("Activity: %s %s | User-Agent: %s | IP: %s | Duration: %v",
		r.Method, r.URL.Path, userAgent, ipAddress, duration)
}