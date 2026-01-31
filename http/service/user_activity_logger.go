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
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("Activity: %s %s | User-Agent: %s | IP: %s | Duration: %v",
		r.Method, r.URL.Path, userAgent, ipAddress, duration)
}package middleware

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	limiter *rate.Limiter
	logger  *log.Logger
}

func NewActivityLogger(rps int, logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{
		limiter: rate.NewLimiter(rate.Limit(rps), rps),
		logger:  logger,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !al.limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)
			al.logger.Printf(
				"Method=%s Path=%s Status=%d Duration=%s RemoteAddr=%s UserAgent=%s",
				r.Method,
				r.URL.Path,
				recorder.statusCode,
				duration,
				r.RemoteAddr,
				r.UserAgent(),
			)
		}()

		next.ServeHTTP(recorder, r)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}