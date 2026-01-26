package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
}

type ActivityLogger struct {
	logFile *os.File
	encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{
		logFile: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (al *ActivityLogger) LogActivity(userID, action, resource string) error {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}
	return al.encoder.Encode(logEntry)
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	activities := []struct {
		userID   string
		action   string
		resource string
	}{
		{"user_001", "LOGIN", "auth_system"},
		{"user_002", "CREATE", "document_123"},
		{"user_001", "UPDATE", "profile_settings"},
		{"user_003", "DELETE", "comment_456"},
	}

	for _, act := range activities {
		if err := logger.LogActivity(act.userID, act.action, act.resource); err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		} else {
			fmt.Printf("Logged: %s performed %s on %s\n", act.userID, act.action, act.resource)
		}
	}
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

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"context"
	"net/http"
	"time"
)

type ActivityLogger struct {
	store    ActivityStore
	limiter  RateLimiter
	timeout  time.Duration
}

type ActivityStore interface {
	Log(ctx context.Context, userID string, action string, metadata map[string]interface{}) error
}

type RateLimiter interface {
	Allow(userID string) bool
}

func NewActivityLogger(store ActivityStore, limiter RateLimiter, timeout time.Duration) *ActivityLogger {
	return &ActivityLogger{
		store:   store,
		limiter: limiter,
		timeout: timeout,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !al.limiter.Allow(userID) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), al.timeout)
		defer cancel()

		action := r.Method + " " + r.URL.Path
		metadata := map[string]interface{}{
			"user_agent": r.UserAgent(),
			"ip_address": r.RemoteAddr,
			"timestamp":  time.Now().UTC(),
		}

		go func() {
			if err := al.store.Log(ctx, userID, action, metadata); err != nil {
				logError(ctx, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return parseToken(auth)
	}
	return ""
}

func parseToken(token string) string {
	return token
}

func logError(ctx context.Context, err error) {
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
	requests, exists := al.rateLimiter[clientIP]

	if !exists {
		al.rateLimiter[clientIP] = []time.Time{now}
		return true
	}

	var validRequests []time.Time
	for _, t := range requests {
		if now.Sub(t) <= al.window {
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
			al.cleanupOldEntries()
		}
	}()
}

func (al *ActivityLogger) cleanupOldEntries() {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	for ip, requests := range al.rateLimiter {
		var validRequests []time.Time
		for _, t := range requests {
			if now.Sub(t) <= al.window {
				validRequests = append(validRequests, t)
			}
		}
		if len(validRequests) == 0 {
			delete(al.rateLimiter, ip)
		} else {
			al.rateLimiter[ip] = validRequests
		}
	}
}