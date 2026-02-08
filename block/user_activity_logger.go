package middleware

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	limiter *rate.Limiter
	store   ActivityStore
}

type ActivityStore interface {
	LogActivity(ctx context.Context, userID string, action string, metadata map[string]interface{}) error
}

func NewActivityLogger(store ActivityStore, requestsPerSecond float64, burst int) *ActivityLogger {
	return &ActivityLogger{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		store:   store,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if !al.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		userID := extractUserID(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		action := r.Method + " " + r.URL.Path
		metadata := map[string]interface{}{
			"user_agent": r.UserAgent(),
			"ip_address": r.RemoteAddr,
			"timestamp":  time.Now().UTC(),
		}

		go func() {
			if err := al.store.LogActivity(ctx, userID, action, metadata); err != nil && ctx.Err() == nil {
				logError(r.Context(), "failed to log activity", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return parseToken(auth)
	}
	return r.URL.Query().Get("user_id")
}

func parseToken(token string) string {
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}
	return ""
}

func logError(ctx context.Context, msg string, err error) {
	select {
	case <-ctx.Done():
		return
	default:
	}
}