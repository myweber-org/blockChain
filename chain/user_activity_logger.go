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

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
	return al.encoder.Encode(logEntry)
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activity.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	activities := []struct {
		userID  string
		action  string
		details string
	}{
		{"user_001", "LOGIN", "Successful authentication"},
		{"user_002", "VIEW_PROFILE", "Accessed user profile page"},
		{"user_001", "UPDATE_SETTINGS", "Changed notification preferences"},
		{"user_003", "LOGOUT", "Session terminated"},
	}

	for _, act := range activities {
		if err := logger.LogActivity(act.userID, act.action, act.details); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}

	fmt.Println("Activity logging completed. Check user_activity.jsonl for results.")
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
		"User activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}