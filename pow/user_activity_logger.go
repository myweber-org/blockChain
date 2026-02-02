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
package middleware

import (
	"context"
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
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = r.RemoteAddr
	}

	if !al.checkRateLimit(userID) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	al.logActivity(userID, r.Method, r.URL.Path)

	ctx := context.WithValue(r.Context(), "userID", userID)
	al.nextHandler.ServeHTTP(w, r.WithContext(ctx))
}

func (al *ActivityLogger) checkRateLimit(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	activities := al.activities[userID]
	var validActivities []time.Time
	for _, t := range activities {
		if t.After(windowStart) {
			validActivities = append(validActivities, t)
		}
	}

	if len(validActivities) >= al.rateLimit {
		return false
	}

	validActivities = append(validActivities, now)
	al.activities[userID] = validActivities
	return true
}

func (al *ActivityLogger) logActivity(userID, method, path string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	entry := struct {
		UserID   string    `json:"user_id"`
		Method   string    `json:"method"`
		Path     string    `json:"path"`
		Time     time.Time `json:"timestamp"`
	}{
		UserID:   userID,
		Method:   method,
		Path:     path,
		Time:     time.Now(),
	}

	_ = entry
}