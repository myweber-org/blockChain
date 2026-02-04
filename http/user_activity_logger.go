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
        return fmt.Errorf("failed to write log entry: %w", err)
    }

    return nil
}

func (l *ActivityLogger) ReadRecentActivities(limit int) ([]ActivityLog, error) {
    file, err := os.Open(l.logFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []ActivityLog{}, nil
        }
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    var logs []ActivityLog
    decoder := json.NewDecoder(file)
    for decoder.More() {
        var log ActivityLog
        if err := decoder.Decode(&log); err != nil {
            return logs, fmt.Errorf("failed to decode log entry: %w", err)
        }
        logs = append(logs, log)
    }

    if len(logs) > limit {
        logs = logs[len(logs)-limit:]
    }

    return logs, nil
}

func main() {
    logger := NewActivityLogger("user_activity.log")

    if err := logger.LogActivity("user123", "login", "User logged in from web browser"); err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    if err := logger.LogActivity("user123", "upload", "Uploaded file: report.pdf"); err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    activities, err := logger.ReadRecentActivities(5)
    if err != nil {
        fmt.Printf("Error reading activities: %v\n", err)
    }

    fmt.Printf("Recent activities (%d):\n", len(activities))
    for _, activity := range activities {
        fmt.Printf("[%s] %s: %s - %s\n",
            activity.Timestamp.Format("2006-01-02 15:04:05"),
            activity.UserID,
            activity.Action,
            activity.Details)
    }
}