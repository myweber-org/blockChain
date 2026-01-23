
package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userHits    map[string][]time.Time
	windowSize  time.Duration
	maxRequests int
}

func NewActivityLogger(window time.Duration, limit int) *ActivityLogger {
	return &ActivityLogger{
		userHits:    make(map[string][]time.Time),
		windowSize:  window,
		maxRequests: limit,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP := r.RemoteAddr
		if !al.allowRequest(userIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		al.logActivity(userIP, r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) allowRequest(ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	hits := al.userHits[ip]

	var validHits []time.Time
	for _, hit := range hits {
		if now.Sub(hit) <= al.windowSize {
			validHits = append(validHits, hit)
		}
	}

	if len(validHits) >= al.maxRequests {
		return false
	}

	validHits = append(validHits, now)
	al.userHits[ip] = validHits
	return true
}

func (al *ActivityLogger) logActivity(ip, path, method string) {
	log.Printf("Activity: %s %s %s", ip, method, path)
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(al.windowSize * 2)
	go func() {
		for range ticker.C {
			al.cleanOldEntries()
		}
	}()
}

func (al *ActivityLogger) cleanOldEntries() {
	al.mu.Lock()
	defer al.mu.Unlock()

	cutoff := time.Now().Add(-al.windowSize)
	for ip, hits := range al.userHits {
		var validHits []time.Time
		for _, hit := range hits {
			if hit.After(cutoff) {
				validHits = append(validHits, hit)
			}
		}
		if len(validHits) == 0 {
			delete(al.userHits, ip)
		} else {
			al.userHits[ip] = validHits
		}
	}
}