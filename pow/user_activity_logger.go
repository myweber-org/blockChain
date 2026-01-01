package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	requests map[string][]time.Time
	window   time.Duration
	maxReqs  int
}

func NewActivityLogger(window time.Duration, maxReqs int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: &RateLimiter{
			requests: make(map[string][]time.Time),
			window:   window,
			maxReqs:  maxReqs,
		},
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path

		if !al.rateLimiter.Allow(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for IP: %s", clientIP)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("Activity: IP=%s, UA=%s, Path=%s, Duration=%v", 
			clientIP, userAgent, path, duration)
	})
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	requests := rl.requests[ip]

	var validReqs []time.Time
	for _, t := range requests {
		if now.Sub(t) <= rl.window {
			validReqs = append(validReqs, t)
		}
	}

	if len(validReqs) >= rl.maxReqs {
		return false
	}

	validReqs = append(validReqs, now)
	rl.requests[ip] = validReqs
	return true
}

func (rl *RateLimiter) Cleanup() {
	ticker := time.NewTicker(rl.window * 2)
	go func() {
		for range ticker.C {
			now := time.Now()
			for ip, requests := range rl.requests {
				var validReqs []time.Time
				for _, t := range requests {
					if now.Sub(t) <= rl.window {
						validReqs = append(validReqs, t)
					}
				}
				if len(validReqs) == 0 {
					delete(rl.requests, ip)
				} else {
					rl.requests[ip] = validReqs
				}
			}
		}
	}()
}