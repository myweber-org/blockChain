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
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    logData, err := json.MarshalIndent(logEntry, "", "  ")
    if err != nil {
        log.Printf("Failed to marshal log entry: %v", err)
        return
    }

    fmt.Println(string(logData))

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v", err)
        return
    }
    defer file.Close()

    if _, err := file.Write(append(logData, '\n')); err != nil {
        log.Printf("Failed to write to log file: %v", err)
    }
}

func main() {
    logActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
    logActivity("user456", "LOGOUT", "Session duration: 1h 30m")
    logActivity("user789", "UPDATE_PROFILE", "Changed email address")
}