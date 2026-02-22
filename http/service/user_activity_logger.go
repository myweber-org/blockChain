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
}