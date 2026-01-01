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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
    Metadata  string    `json:"metadata,omitempty"`
}

func logActivity(userID, eventType, metadata string) ActivityEvent {
    event := ActivityEvent{
        UserID:    userID,
        EventType: eventType,
        Timestamp: time.Now().UTC(),
        Metadata:  metadata,
    }

    logEntry, err := json.Marshal(event)
    if err != nil {
        fmt.Printf("Failed to marshal event: %v\n", err)
        return event
    }

    fmt.Fprintf(os.Stdout, "%s\n", logEntry)
    return event
}

func main() {
    logActivity("user_123", "login", "from_ip:192.168.1.1")
    logActivity("user_456", "purchase", "item_id:789,amount:29.99")
    logActivity("user_123", "logout", "")
}