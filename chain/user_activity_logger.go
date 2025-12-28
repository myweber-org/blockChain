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
    
    log.Printf("Activity: %s %s from %s took %v", r.Method, r.URL.Path, r.RemoteAddr, duration)
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
    return &ActivityLogger{handler: handler}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Request processed successfully"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", mainHandler)
    
    wrappedMux := NewActivityLogger(mux)
    
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", wrappedMux)
}