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

func NewActivityLog(userID, action, details string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.MarshalIndent(al, "", "  ")
}

func LogActivity(logFile *os.File, userID, action, details string) error {
	activity := NewActivityLog(userID, action, details)
	jsonData, err := activity.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal activity: %w", err)
	}

	if _, err := logFile.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	log.Printf("Logged activity: %s by user %s", action, userID)
	return nil
}

func main() {
	logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	activities := []struct {
		userID string
		action string
		details string
	}{
		{"user123", "LOGIN", "Successful authentication"},
		{"user123", "VIEW_PROFILE", "Accessed profile page"},
		{"user456", "UPDATE_SETTINGS", "Changed notification preferences"},
		{"user123", "LOGOUT", "Session terminated"},
	}

	for _, act := range activities {
		if err := LogActivity(logFile, act.userID, act.action, act.details); err != nil {
			log.Printf("Error logging activity: %v", err)
		}
	}

	fmt.Println("Activity logging completed. Check activity.log for details.")
}