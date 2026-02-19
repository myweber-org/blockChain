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

func (al *ActivityLog) ToJSON() (string, error) {
    data, err := json.MarshalIndent(al, "", "  ")
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func LogActivity(logger *log.Logger, userID, action, details string) {
    activity := NewActivityLog(userID, action, details)
    jsonStr, err := activity.ToJSON()
    if err != nil {
        logger.Printf("Failed to marshal activity log: %v", err)
        return
    }
    logger.Println(jsonStr)
}

func main() {
    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    logger := log.New(file, "", 0)

    LogActivity(logger, "user123", "LOGIN", "User logged in from IP 192.168.1.100")
    LogActivity(logger, "user456", "UPDATE_PROFILE", "Changed email address")
    LogActivity(logger, "user789", "LOGOUT", "Session terminated after 2 hours")

    fmt.Println("Activity logging completed. Check activity.log file.")
}