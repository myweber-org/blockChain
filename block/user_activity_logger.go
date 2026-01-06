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

func NewActivityLogger(window time.Duration, max int) *ActivityLogger {
	return &ActivityLogger{
		userHits:    make(map[string][]time.Time),
		windowSize:  window,
		maxRequests: max,
	}
}

func (al *ActivityLogger) LogActivity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIP := r.RemoteAddr
		currentTime := time.Now()

		al.mu.Lock()
		defer al.mu.Unlock()

		hits := al.userHits[userIP]
		validHits := []time.Time{}
		for _, hit := range hits {
			if currentTime.Sub(hit) <= al.windowSize {
				validHits = append(validHits, hit)
			}
		}

		if len(validHits) >= al.maxRequests {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for IP: %s", userIP)
			return
		}

		validHits = append(validHits, currentTime)
		al.userHits[userIP] = validHits

		log.Printf("Activity from IP %s at %s", userIP, currentTime.Format(time.RFC3339))
		next.ServeHTTP(w, r)
	}
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(al.windowSize * 2)
	go func() {
		for range ticker.C {
			al.mu.Lock()
			cutoff := time.Now().Add(-al.windowSize)
			for ip, hits := range al.userHits {
				validHits := []time.Time{}
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
			al.mu.Unlock()
		}
	}()
}