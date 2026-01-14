package main

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
	if err := logActivity("user123", "login", "successful authentication"); err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	if err := logActivity("user456", "file_upload", "uploaded profile.jpg"); err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	fmt.Println("Activity logging completed")
}package main

import (
    "encoding/json"
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
    Metadata  interface{}  `json:"metadata,omitempty"`
}

type ActivityLogger struct {
    logFile *os.File
    encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{
        logFile: file,
        encoder: json.NewEncoder(file),
    }, nil
}

func (l *ActivityLogger) LogActivity(userID string, action ActivityType, metadata interface{}) error {
    activity := UserActivity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Metadata:  metadata,
    }
    return l.encoder.Encode(activity)
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("user_activities.jsonl")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    err = logger.LogActivity("user_123", Login, map[string]string{"ip": "192.168.1.1"})
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logger.LogActivity("user_456", View, map[string]string{"page": "/products/42"})
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logger.LogActivity("user_123", Purchase, map[string]interface{}{
        "item_id": "prod_789",
        "amount":  49.99,
    })
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }
}