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
		startTime := time.Now()
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(startTime)
		
		al.Logger.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | UserAgent: %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.UserAgent(),
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
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    err = logger.LogActivity("user123", "login", "Successful authentication")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }

    err = logger.LogActivity("user123", "file_upload", "uploaded profile.jpg")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }

    fmt.Println("Activity logging completed")
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityLog struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(logEntry); err != nil {
        return err
    }

    return nil
}

func main() {
    err := logActivity("user123", "login", "Successful authentication")
    if err != nil {
        log.Fatal("Failed to log activity:", err)
    }
    fmt.Println("Activity logged successfully")
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

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

type ActivityLogger struct {
    logFile string
}

func NewActivityLogger(logFile string) *ActivityLogger {
    return &ActivityLogger{logFile: logFile}
}

func (l *ActivityLogger) LogActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (l *ActivityLogger) ReadLogs() ([]ActivityLog, error) {
    data, err := os.ReadFile(l.logFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []ActivityLog{}, nil
        }
        return nil, fmt.Errorf("failed to read log file: %w", err)
    }

    var logs []ActivityLog
    lines := bytes.Split(data, []byte("\n"))
    for _, line := range lines {
        if len(line) == 0 {
            continue
        }
        var logEntry ActivityLog
        if err := json.Unmarshal(line, &logEntry); err != nil {
            return nil, fmt.Errorf("failed to unmarshal log entry: %w", err)
        }
        logs = append(logs, logEntry)
    }

    return logs, nil
}

func main() {
    logger := NewActivityLogger("user_activity.log")

    err := logger.LogActivity("user123", "login", "User logged in from IP 192.168.1.100")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
        return
    }

    err = logger.LogActivity("user123", "view_page", "Viewed dashboard page")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
        return
    }

    logs, err := logger.ReadLogs()
    if err != nil {
        fmt.Printf("Error reading logs: %v\n", err)
        return
    }

    fmt.Printf("Total logged activities: %d\n", len(logs))
    for _, log := range logs {
        fmt.Printf("%s - %s: %s\n", log.Timestamp.Format(time.RFC3339), log.Action, log.Details)
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
}