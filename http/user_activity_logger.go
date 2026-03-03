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
		"Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ActivityLog struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
}

type ActivityLogger struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		if !al.allowRequest(userID) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		activity := ActivityLog{
			UserID:    userID,
			Action:    r.URL.Query().Get("action"),
			Path:      r.URL.Path,
			Method:    r.Method,
			Timestamp: time.Now(),
			IPAddress: r.RemoteAddr,
		}

		al.recordActivity(userID)

		logData, err := json.Marshal(activity)
		if err == nil {
			go func() {
				println(string(logData))
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) allowRequest(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	requests := al.requests[userID]
	var validRequests []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= al.limit {
		return false
	}

	al.requests[userID] = validRequests
	return true
}

func (al *ActivityLogger) recordActivity(userID string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	al.requests[userID] = append(al.requests[userID], time.Now())
}