package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type ActivityEntry struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
}

type ActivityLogger struct {
	mu          sync.RWMutex
	activities  []ActivityEntry
	rateLimiter map[string]time.Time
	window      time.Duration
	maxEntries  int
}

func NewActivityLogger(window time.Duration, maxEntries int) *ActivityLogger {
	return &ActivityLogger{
		activities:  make([]ActivityEntry, 0, maxEntries),
		rateLimiter: make(map[string]time.Time),
		window:      window,
		maxEntries:  maxEntries,
	}
}

func (al *ActivityLogger) LogActivity(userID, action, ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	key := userID + ":" + action
	if lastTime, exists := al.rateLimiter[key]; exists {
		if time.Since(lastTime) < al.window {
			return false
		}
	}

	entry := ActivityEntry{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
		IPAddress: ip,
	}

	al.activities = append(al.activities, entry)
	al.rateLimiter[key] = entry.Timestamp

	if len(al.activities) > al.maxEntries {
		al.activities = al.activities[1:]
	}

	return true
}

func (al *ActivityLogger) GetActivities(since time.Time) []ActivityEntry {
	al.mu.RLock()
	defer al.mu.RUnlock()

	var result []ActivityEntry
	for _, entry := range al.activities {
		if entry.Timestamp.After(since) {
			result = append(result, entry)
		}
	}
	return result
}

func ActivityLoggingMiddleware(al *ActivityLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		action := r.Method + " " + r.URL.Path
		ip := r.RemoteAddr

		if al.LogActivity(userID, action, ip) {
			ctx = context.WithValue(ctx, "activity_logged", true)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}