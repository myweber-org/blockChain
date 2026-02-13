package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s %d %v",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
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
	startTime := time.Now()
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(startTime)
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package middleware

import (
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	activities  map[string][]time.Time
	rateLimit   int
	window      time.Duration
	nextHandler http.Handler
}

func NewActivityLogger(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &ActivityLogger{
			activities:  make(map[string][]time.Time),
			rateLimit:   limit,
			window:      window,
			nextHandler: next,
		}
	}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIP := r.RemoteAddr
	now := time.Now()

	al.mu.Lock()
	windowStart := now.Add(-al.window)
	userActivities := al.activities[userIP]

	var validActivities []time.Time
	for _, activity := range userActivities {
		if activity.After(windowStart) {
			validActivities = append(validActivities, activity)
		}
	}

	if len(validActivities) >= al.rateLimit {
		al.mu.Unlock()
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	validActivities = append(validActivities, now)
	al.activities[userIP] = validActivities
	al.mu.Unlock()

	al.nextHandler.ServeHTTP(w, r)
}