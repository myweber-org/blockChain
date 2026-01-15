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

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityRecorder struct {
	ResponseWriter http.ResponseWriter
	StatusCode     int
}

func (ar *ActivityRecorder) WriteHeader(code int) {
	ar.StatusCode = code
	ar.ResponseWriter.WriteHeader(code)
}

func (ar *ActivityRecorder) Header() http.Header {
	return ar.ResponseWriter.Header()
}

func (ar *ActivityRecorder) Write(b []byte) (int, error) {
	return ar.ResponseWriter.Write(b)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &ActivityRecorder{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			recorder.StatusCode,
			duration,
			r.RemoteAddr,
		)
	})
}