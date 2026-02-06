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
            "%s %s %d %s %s",
            r.Method,
            r.URL.Path,
            rw.statusCode,
            duration.String(),
            r.RemoteAddr,
        )
    })
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter map[string]time.Time
	window      time.Duration
}

func NewActivityLogger(window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string]time.Time),
		window:      window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		now := time.Now()

		if lastSeen, exists := al.rateLimiter[clientIP]; exists {
			if now.Sub(lastSeen) < al.window {
				log.Printf("Rate limit exceeded for IP: %s", clientIP)
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}

		al.rateLimiter[clientIP] = now
		log.Printf("Activity from %s: %s %s", clientIP, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) CleanupOldEntries() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for ip, lastSeen := range al.rateLimiter {
			if now.Sub(lastSeen) > 24*time.Hour {
				delete(al.rateLimiter, ip)
			}
		}
	}
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type UserActivity struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) error {
    activity := UserActivity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    file, err := os.OpenFile("user_activities.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(activity); err != nil {
        return fmt.Errorf("failed to encode activity: %w", err)
    }

    return nil
}

func main() {
    activities := []struct {
        userID, action, details string
    }{
        {"user_123", "login", "Successful authentication"},
        {"user_456", "purchase", "Order ID: ORD-78910"},
        {"user_123", "logout", "Session ended"},
    }

    for _, act := range activities {
        if err := logActivity(act.userID, act.action, act.details); err != nil {
            log.Printf("Failed to log activity: %v", err)
        } else {
            fmt.Printf("Logged %s for user %s\n", act.action, act.userID)
        }
    }
}