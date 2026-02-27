package main

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
    if err := encoder.Encode(logEntry); err != nil {
        return err
    }

    fmt.Printf("Logged: %s performed %s at %s\n", userID, action, logEntry.Timestamp.Format(time.RFC3339))
    return nil
}

func main() {
    if err := logActivity("user123", "login", "Successful authentication"); err != nil {
        log.Fatal(err)
    }

    time.Sleep(1 * time.Second)

    if err := logActivity("user123", "file_upload", "Uploaded profile.jpg"); err != nil {
        log.Fatal(err)
    }
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
	Details   string    `json:"details"`
}

func logActivity(userID, action, details string) {
	activity := ActivityLog{
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	logEntry, err := json.Marshal(activity)
	if err != nil {
		log.Printf("Failed to marshal activity: %v", err)
		return
	}

	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(append(logEntry, '\n')); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
}

func main() {
	logActivity("user123", "login", "User logged in from IP 192.168.1.100")
	logActivity("user456", "purchase", "Purchased item SKU: PROD-789")
	logActivity("user123", "logout", "Session ended after 15 minutes")

	fmt.Println("Activity logging completed. Check activity.log file.")
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

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
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

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
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
        {"user_123", "logout", "Session duration: 45m"},
    }

    for _, a := range activities {
        if err := logActivity(a.userID, a.action, a.details); err != nil {
            log.Printf("Failed to log activity: %v", err)
        } else {
            fmt.Printf("Logged: %s performed %s\n", a.userID, a.action)
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
			duration,
			r.RemoteAddr,
		)
	})
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
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}