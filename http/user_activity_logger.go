package main

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
	Metadata  string    `json:"metadata,omitempty"`
}

func NewActivity(userID, action, metadata string) *Activity {
	return &Activity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().UTC(),
		Metadata:  metadata,
	}
}

func (a *Activity) ToJSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

func LogActivity(activity *Activity, logFile string) error {
	jsonData, err := activity.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal activity: %w", err)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	jsonData = append(jsonData, '\n')
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

func main() {
	activity := NewActivity("user123", "login", "from_ip:192.168.1.100")
	
	if err := LogActivity(activity, "activity.log"); err != nil {
		log.Fatalf("Failed to log activity: %v", err)
	}
	
	fmt.Println("Activity logged successfully")
}