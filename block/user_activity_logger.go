
package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	activities  map[string][]time.Time
	rateLimit   int
	window      time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		activities: make(map[string][]time.Time),
		rateLimit:  limit,
		window:     window,
	}
}

func (al *ActivityLogger) recordActivity(userID string) bool {
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

func (al *ActivityLogger) cleanupOldActivities() {
	ticker := time.NewTicker(al.window * 2)
	defer ticker.Stop()

	for range ticker.C {
		al.mu.Lock()
		windowStart := time.Now().Add(-al.window)
		for userID, activities := range al.activities {
			var validActivities []time.Time
			for _, t := range activities {
				if t.After(windowStart) {
					validActivities = append(validActivities, t)
				}
			}
			if len(validActivities) == 0 {
				delete(al.activities, userID)
			} else {
				al.activities[userID] = validActivities
			}
		}
		al.mu.Unlock()
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !al.recordActivity(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		log.Printf("Activity recorded for user %s: %s %s", userID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := NewActivityLogger(10, time.Minute)
	go logger.cleanupOldActivities()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Data response"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: logger.Middleware(mux),
	}

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}