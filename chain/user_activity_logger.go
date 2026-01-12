package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

func logActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(logEntry)
}

func main() {
    err := logActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
    if err != nil {
        fmt.Printf("Failed to log activity: %v\n", err)
        return
    }
    fmt.Println("Activity logged successfully")
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IP        string
	Duration  time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		activity := ActivityLog{
			Timestamp: start,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
			Duration:  duration,
		}
		
		logActivity(activity, recorder.statusCode)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func logActivity(activity ActivityLog, status int) {
	log.Printf(
		"ACTIVITY | %s | %s %s | Status: %d | Duration: %v | IP: %s | Agent: %s",
		activity.Timestamp.Format("2006-01-02 15:04:05"),
		activity.Method,
		activity.Path,
		status,
		activity.Duration,
		activity.IP,
		activity.UserAgent,
	)
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

	log.Printf(
		"Method: %s | Path: %s | Duration: %v | Timestamp: %s",
		r.Method,
		r.URL.Path,
		duration,
		time.Now().Format(time.RFC3339),
	)
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
}