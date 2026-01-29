package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type Activity struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

type ActivityLogger struct {
    logFile *os.File
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{logFile: file}, nil
}

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
    activity := Activity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    data, err := json.Marshal(activity)
    if err != nil {
        return err
    }

    data = append(data, '\n')
    _, err = al.logFile.Write(data)
    return err
}

func (al *ActivityLogger) Close() error {
    return al.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("user_activities.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    err = logger.LogActivity("user123", "login", "Successful authentication")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logger.LogActivity("user123", "file_upload", "uploaded profile.jpg")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    fmt.Println("Activities logged successfully")
}package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	rateLimiter map[string][]time.Time
	window      time.Duration
	maxRequests int
}

func NewActivityLogger(window time.Duration, maxRequests int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string][]time.Time),
		window:      window,
		maxRequests: maxRequests,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path
		method := r.Method

		if !al.allowRequest(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("Activity: %s %s %s %s Duration: %v", clientIP, method, path, userAgent, duration)
	})
}

func (al *ActivityLogger) allowRequest(clientIP string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	if _, exists := al.rateLimiter[clientIP]; !exists {
		al.rateLimiter[clientIP] = []time.Time{}
	}

	requests := al.rateLimiter[clientIP]
	var validRequests []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= al.maxRequests {
		return false
	}

	validRequests = append(validRequests, now)
	al.rateLimiter[clientIP] = validRequests
	return true
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(al.window * 2)
	go func() {
		for range ticker.C {
			al.mu.Lock()
			for ip, requests := range al.rateLimiter {
				var validRequests []time.Time
				windowStart := time.Now().Add(-al.window)
				for _, t := range requests {
					if t.After(windowStart) {
						validRequests = append(validRequests, t)
					}
				}
				if len(validRequests) == 0 {
					delete(al.rateLimiter, ip)
				} else {
					al.rateLimiter[ip] = validRequests
				}
			}
			al.mu.Unlock()
		}
	}()
}package middleware

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
			"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s | UserAgent: %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)
	})
}