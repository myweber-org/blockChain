package main

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("[%s] %s %s - %v", r.RemoteAddr, r.Method, r.URL.Path, duration)
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", apiHandler)

	wrappedMux := NewActivityLogger(mux)

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", wrappedMux)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

type ActivityLogger struct {
	store chan ActivityLog
}

func NewActivityLogger(bufferSize int) *ActivityLogger {
	al := &ActivityLogger{
		store: make(chan ActivityLog, bufferSize),
	}
	go al.processLogs()
	return al
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		activity := ActivityLog{
			UserID:    userID,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: time.Now(),
			IPAddress: r.RemoteAddr,
		}

		select {
		case al.store <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.store {
		log.Printf("ACTIVITY: User=%s %s %s from %s at %v",
			activity.UserID,
			activity.Method,
			activity.Endpoint,
			activity.IPAddress,
			activity.Timestamp.Format(time.RFC3339))
	}
}

func (al *ActivityLogger) Close() {
	close(al.store)
}