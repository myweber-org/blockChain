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

	log.Printf("Activity: %s %s from %s took %v",
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

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.MarshalIndent(al, "", "  ")
}

func LogActivity(userID, action, details string) error {
	logEntry := NewActivityLog(userID, action, details)
	jsonData, err := logEntry.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	if _, err := logFile.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	log.Printf("Logged activity: %s by user %s", action, userID)
	return nil
}

func main() {
	if err := LogActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100"); err != nil {
		log.Fatal(err)
	}

	if err := LogActivity("user123", "VIEW_PAGE", "Accessed dashboard page"); err != nil {
		log.Fatal(err)
	}

	if err := LogActivity("user456", "UPDATE_PROFILE", "Changed email address"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Activity logging completed successfully")
}