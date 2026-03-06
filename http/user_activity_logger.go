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

	log.Printf(
		"Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ActivityLog struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
}

type ActivityLogger struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		if !al.allowRequest(userID) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		activity := ActivityLog{
			UserID:    userID,
			Action:    r.URL.Query().Get("action"),
			Path:      r.URL.Path,
			Method:    r.Method,
			Timestamp: time.Now(),
			IPAddress: r.RemoteAddr,
		}

		al.recordActivity(userID)

		logData, err := json.Marshal(activity)
		if err == nil {
			go func() {
				println(string(logData))
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) allowRequest(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	requests := al.requests[userID]
	var validRequests []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= al.limit {
		return false
	}

	al.requests[userID] = validRequests
	return true
}

func (al *ActivityLogger) recordActivity(userID string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	al.requests[userID] = append(al.requests[userID], time.Now())
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
    "log"
    "os"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details"`
}

type ActivityLogger struct {
    file *os.File
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{file: file}, nil
}

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
    event := ActivityEvent{
        UserID:    userID,
        EventType: eventType,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    data = append(data, '\n')
    _, err = l.file.Write(data)
    return err
}

func (l *ActivityLogger) Close() error {
    return l.file.Close()
}

func main() {
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    err = logger.LogActivity("user123", "login", "User logged in from web browser")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logger.LogActivity("user123", "search", "Searched for 'golang tutorials'")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logger.LogActivity("user456", "purchase", "Purchased item ID: 789")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    log.Println("Activity logging completed")
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
	IPAddress string
	Duration  time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		logEntry := ActivityLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
		}

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		logEntry.Duration = time.Since(start)
		
		log.Printf("ACTIVITY: %s %s | IP: %s | Agent: %s | Status: %d | Duration: %v",
			logEntry.Method,
			logEntry.Path,
			logEntry.IPAddress,
			logEntry.UserAgent,
			recorder.statusCode,
			logEntry.Duration,
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
    logFile string
}

func NewActivityLogger(logFile string) *ActivityLogger {
    return &ActivityLogger{logFile: logFile}
}

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
    activity := Activity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    file, err := os.OpenFile(al.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (al *ActivityLogger) ReadActivities() ([]Activity, error) {
    data, err := os.ReadFile(al.logFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []Activity{}, nil
        }
        return nil, fmt.Errorf("failed to read log file: %w", err)
    }

    var activities []Activity
    lines := bytes.Split(data, []byte("\n"))
    for _, line := range lines {
        if len(line) == 0 {
            continue
        }
        var activity Activity
        if err := json.Unmarshal(line, &activity); err != nil {
            return nil, fmt.Errorf("failed to unmarshal activity: %w", err)
        }
        activities = append(activities, activity)
    }

    return activities, nil
}

func main() {
    logger := NewActivityLogger("user_activities.json")

    if err := logger.LogActivity("user123", "login", "Successful authentication"); err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    if err := logger.LogActivity("user123", "upload", "File: report.pdf"); err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    activities, err := logger.ReadActivities()
    if err != nil {
        fmt.Printf("Error reading activities: %v\n", err)
        return
    }

    fmt.Printf("Total activities logged: %d\n", len(activities))
    for _, activity := range activities {
        fmt.Printf("%s - %s: %s\n", activity.Timestamp.Format(time.RFC3339), activity.UserID, activity.Action)
    }
}