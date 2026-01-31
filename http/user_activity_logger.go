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
	log.Printf("Activity: %s %s | IP: %s | Agent: %s | Duration: %v",
		r.Method,
		r.URL.Path,
		ipAddress,
		userAgent,
		duration,
	)
}package middleware

import (
    "context"
    "net/http"
    "time"

    "github.com/google/uuid"
    "go.uber.org/zap"
)

type activityKey struct{}

type UserActivity struct {
    RequestID    string
    UserID       string
    Action       string
    Resource     string
    Timestamp    time.Time
    Status       int
    Duration     time.Duration
    ClientIP     string
    UserAgent    string
}

func ActivityLogger(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            requestID := uuid.New().String()
            
            activity := UserActivity{
                RequestID: requestID,
                UserID:    extractUserID(r),
                Action:    r.Method,
                Resource:  r.URL.Path,
                Timestamp: start,
                ClientIP:  r.RemoteAddr,
                UserAgent: r.UserAgent(),
            }

            ctx := context.WithValue(r.Context(), activityKey{}, &activity)
            r = r.WithContext(ctx)

            rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            next.ServeHTTP(rw, r)

            activity.Status = rw.statusCode
            activity.Duration = time.Since(start)

            logActivity(logger, activity)
        })
    }
}

func GetActivity(ctx context.Context) *UserActivity {
    if activity, ok := ctx.Value(activityKey{}).(*UserActivity); ok {
        return activity
    }
    return nil
}

func UpdateActivity(ctx context.Context, updates func(*UserActivity)) {
    if activity := GetActivity(ctx); activity != nil {
        updates(activity)
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func extractUserID(r *http.Request) string {
    if user := r.Header.Get("X-User-ID"); user != "" {
        return user
    }
    return "anonymous"
}

func logActivity(logger *zap.Logger, activity UserActivity) {
    fields := []zap.Field{
        zap.String("request_id", activity.RequestID),
        zap.String("user_id", activity.UserID),
        zap.String("action", activity.Action),
        zap.String("resource", activity.Resource),
        zap.Time("timestamp", activity.Timestamp),
        zap.Int("status", activity.Status),
        zap.Duration("duration", activity.Duration),
        zap.String("client_ip", activity.ClientIP),
        zap.String("user_agent", activity.UserAgent),
    }

    if activity.Status >= 400 {
        logger.Warn("User activity completed with error", fields...)
    } else {
        logger.Info("User activity completed successfully", fields...)
    }
}