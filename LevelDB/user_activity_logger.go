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
	Resource  string    `json:"resource"`
}

func NewActivityLog(userID, action, resource string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.Marshal(al)
}

func LogActivity(userID, action, resource string) error {
	logEntry := NewActivityLog(userID, action, resource)
	
	jsonData, err := logEntry.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}
	
	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()
	
	if _, err := file.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}
	
	log.Printf("Logged activity: %s performed %s on %s", userID, action, resource)
	return nil
}

func main() {
	if err := LogActivity("user123", "CREATE", "document.pdf"); err != nil {
		log.Fatal(err)
	}
	
	if err := LogActivity("user456", "DELETE", "image.jpg"); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Activity logging completed")
}package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityLog struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

type ActivityLogger struct {
    logs []ActivityLog
}

func NewActivityLogger() *ActivityLogger {
    return &ActivityLogger{
        logs: make([]ActivityLog, 0),
    }
}

func (al *ActivityLogger) LogActivity(userID, action, details string) {
    log := ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }
    al.logs = append(al.logs, log)
}

func (al *ActivityLogger) GetLogs() []ActivityLog {
    return al.logs
}

func (al *ActivityLogger) SaveToFile(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(al.logs)
}

func main() {
    logger := NewActivityLogger()
    
    logger.LogActivity("user123", "LOGIN", "User logged in from web browser")
    logger.LogActivity("user123", "VIEW_PROFILE", "Viewed profile page")
    logger.LogActivity("user456", "REGISTER", "New user registration")
    logger.LogActivity("user123", "LOGOUT", "User logged out")
    
    logs := logger.GetLogs()
    for _, log := range logs {
        fmt.Printf("[%s] %s: %s - %s\n", 
            log.Timestamp.Format("2006-01-02 15:04:05"),
            log.UserID,
            log.Action,
            log.Details)
    }
    
    err := logger.SaveToFile("activity_logs.json")
    if err != nil {
        fmt.Printf("Error saving logs: %v\n", err)
    } else {
        fmt.Println("Logs saved to activity_logs.json")
    }
}