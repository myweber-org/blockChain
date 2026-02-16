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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s | UserAgent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
		r.UserAgent(),
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
    "log"
    "net/http"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    Method    string    `json:"method"`
    Path      string    `json:"path"`
    UserAgent string    `json:"user_agent"`
    IPAddress string    `json:"ip_address"`
}

func activityLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        logEntry := ActivityLog{
            Timestamp: start,
            Method:    r.Method,
            Path:      r.URL.Path,
            UserAgent: r.UserAgent(),
            IPAddress: r.RemoteAddr,
        }
        
        logData, err := json.Marshal(logEntry)
        if err != nil {
            log.Printf("Failed to marshal activity log: %v", err)
        } else {
            log.Printf("Activity: %s", string(logData))
        }
        
        next.ServeHTTP(w, r)
    })
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{"status": "ok", "message": "request processed"}
    json.NewEncoder(w).Encode(response)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", mainHandler)
    
    wrappedMux := activityLogger(mux)
    
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
        log.Fatal(err)
    }
}