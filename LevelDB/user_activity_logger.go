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
}