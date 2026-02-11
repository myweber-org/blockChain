package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
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
    _, err = l.logFile.Write(data)
    return err
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
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

    fmt.Println("Activity logging completed")
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
		"Activity: %s %s from %s completed in %v",
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
	
	al.handler.ServeHTTP(w, r)
	
	duration := time.Since(start)
	
	log.Printf("[%s] %s %s - %v - %s",
		time.Now().Format("2006-01-02 15:04:05"),
		r.Method,
		r.URL.Path,
		duration,
		r.RemoteAddr,
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
    return &ActivityLogger{
        logFile: logFile,
    }
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
        var logEntry ActivityLog
        if err := decoder.Decode(&logEntry); err != nil {
            return logs, fmt.Errorf("failed to decode log entry: %w", err)
        }
        logs = append(logs, logEntry)
    }

    if len(logs) > limit {
        logs = logs[len(logs)-limit:]
    }

    return logs, nil
}

func main() {
    logger := NewActivityLogger("user_activity.log")

    // Example usage
    err := logger.LogActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    err = logger.LogActivity("user123", "VIEW_PAGE", "Accessed dashboard page")
    if err != nil {
        fmt.Printf("Error logging activity: %v\n", err)
    }

    recentLogs, err := logger.ReadRecentActivities(5)
    if err != nil {
        fmt.Printf("Error reading logs: %v\n", err)
    }

    fmt.Printf("Recent activities (%d):\n", len(recentLogs))
    for _, log := range recentLogs {
        fmt.Printf("[%s] %s - %s: %s\n", 
            log.Timestamp.Format(time.RFC3339), 
            log.UserID, 
            log.Action, 
            log.Details)
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

func NewActivityLog(userID, action, details string) *ActivityLog {
    return &ActivityLog{
        Timestamp: time.Now().UTC(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }
}

func (al *ActivityLog) SaveToFile(filename string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(al)
}

func main() {
    logger := NewActivityLog("user123", "login", "User logged in from IP 192.168.1.100")
    
    err := logger.SaveToFile("activity_log.json")
    if err != nil {
        log.Fatalf("Failed to save activity log: %v", err)
    }
    
    fmt.Printf("Activity logged successfully at %s\n", logger.Timestamp.Format(time.RFC3339))
}