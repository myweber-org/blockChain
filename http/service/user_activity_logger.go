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
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityType string

const (
    Login    ActivityType = "LOGIN"
    Logout   ActivityType = "LOGOUT"
    Purchase ActivityType = "PURCHASE"
    View     ActivityType = "VIEW"
)

type UserActivity struct {
    UserID    string       `json:"user_id"`
    Action    ActivityType `json:"action"`
    Timestamp time.Time    `json:"timestamp"`
    Details   string       `json:"details,omitempty"`
}

func logActivity(activity UserActivity) error {
    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    activityJSON, err := json.Marshal(activity)
    if err != nil {
        return err
    }

    _, err = file.Write(append(activityJSON, '\n'))
    return err
}

func main() {
    activities := []UserActivity{
        {UserID: "user123", Action: Login, Timestamp: time.Now(), Details: "Successful login"},
        {UserID: "user456", Action: View, Timestamp: time.Now().Add(-5 * time.Minute), Details: "Viewed product catalog"},
        {UserID: "user123", Action: Purchase, Timestamp: time.Now().Add(-2 * time.Minute), Details: "Purchased item ID: 789"},
        {UserID: "user456", Action: Logout, Timestamp: time.Now().Add(-1 * time.Minute)},
    }

    for _, activity := range activities {
        if err := logActivity(activity); err != nil {
            log.Printf("Failed to log activity: %v", err)
        } else {
            fmt.Printf("Logged: %s - %s\n", activity.UserID, activity.Action)
        }
    }
}package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityEvent struct {
	Timestamp time.Time
	UserID    string
	EventType string
	Details   string
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

func (al *ActivityLogger) LogActivity(userID, eventType, details string) {
	event := ActivityEvent{
		Timestamp: time.Now(),
		UserID:    userID,
		EventType: eventType,
		Details:   details,
	}

	logEntry := fmt.Sprintf("%s | User: %s | Event: %s | Details: %s\n",
		event.Timestamp.Format("2006-01-02 15:04:05"),
		event.UserID,
		event.EventType,
		event.Details)

	if _, err := al.logFile.WriteString(logEntry); err != nil {
		log.Printf("Failed to write activity log: %v", err)
	}
}

func (al *ActivityLogger) Close() {
	if al.logFile != nil {
		al.logFile.Close()
	}
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.LogActivity("user123", "LOGIN", "Successful authentication")
	logger.LogActivity("user123", "VIEW_PAGE", "Accessed dashboard")
	logger.LogActivity("user456", "UPDATE_PROFILE", "Changed email address")
	logger.LogActivity("user123", "LOGOUT", "Session terminated")

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

	log.Printf(
		"Method: %s | Path: %s | Duration: %v | Timestamp: %s",
		r.Method,
		r.URL.Path,
		duration,
		time.Now().Format(time.RFC3339),
	)
}
package middleware

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
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s | UserAgent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
		r.UserAgent(),
	)
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

type Activity struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  string    `json:"metadata,omitempty"`
}

func logActivity(userID, action, metadata string) {
	activity := Activity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().UTC(),
		Metadata:  metadata,
	}

	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(activity); err != nil {
		log.Printf("Failed to encode activity: %v", err)
		return
	}

	fmt.Printf("Logged: %s performed %s at %s\n", userID, action, activity.Timestamp.Format(time.RFC3339))
}

func main() {
	logActivity("user123", "login", "browser:chrome")
	logActivity("user456", "upload", "file_size:2.5MB")
	logActivity("user123", "logout", "session_duration:15m")
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
		"Method: %s | Path: %s | Status: %d | Duration: %v | UserAgent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.UserAgent(),
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
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(writer, r)
	
	duration := time.Since(start)
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		writer.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
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

func logActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(logEntry); err != nil {
        return fmt.Errorf("failed to encode log entry: %w", err)
    }

    return nil
}

func main() {
    if err := logActivity("user123", "login", "User logged in from IP 192.168.1.100"); err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    if err := logActivity("user123", "file_upload", "Uploaded document.pdf"); err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    fmt.Println("Activity logging completed")
}