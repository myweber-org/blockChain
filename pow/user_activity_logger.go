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
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
}

func NewActivityLog(sessionID, userID, action, details string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		SessionID: sessionID,
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.MarshalIndent(al, "", "  ")
}

func LogActivity(logger *log.Logger, sessionID, userID, action, details string) {
	activity := NewActivityLog(sessionID, userID, action, details)
	jsonData, err := activity.ToJSON()
	if err != nil {
		logger.Printf("Failed to marshal activity log: %v", err)
		return
	}
	logger.Println(string(jsonData))
}

func main() {
	logger := log.New(os.Stdout, "ACTIVITY: ", log.LstdFlags)
	
	LogActivity(logger, "sess_abc123", "user_789", "LOGIN", "User logged in from IP 192.168.1.100")
	LogActivity(logger, "sess_abc123", "user_789", "VIEW_PAGE", "Accessed dashboard page")
	LogActivity(logger, "sess_abc123", "user_789", "LOGOUT", "User logged out after 15 minutes")
}