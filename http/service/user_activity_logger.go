package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	redisClient *redis.Client
	limiter     *rate.Limiter
}

func NewActivityLogger(redisAddr string) *ActivityLogger {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	return &ActivityLogger{
		redisClient: rdb,
		limiter:     rate.NewLimiter(rate.Every(time.Minute), 10),
	}
}

func (al *ActivityLogger) LogActivity(ctx context.Context, userID, action string) error {
	if !al.limiter.Allow() {
		return fmt.Errorf("rate limit exceeded for user %s", userID)
	}

	key := fmt.Sprintf("activity:%s:%d", userID, time.Now().Unix())
	data := map[string]interface{}{
		"user_id":    userID,
		"action":     action,
		"timestamp":  time.Now().Format(time.RFC3339),
		"user_agent": ctx.Value("User-Agent"),
		"ip_address": ctx.Value("X-Forwarded-For"),
	}

	err := al.redisClient.HSet(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	expiration := 24 * time.Hour
	al.redisClient.Expire(ctx, key, expiration)

	return nil
}

func (al *ActivityLogger) GetRecentActivities(ctx context.Context, userID string, limit int64) ([]map[string]string, error) {
	pattern := fmt.Sprintf("activity:%s:*", userID)
	keys, err := al.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	if int64(len(keys)) > limit {
		keys = keys[:limit]
	}

	var activities []map[string]string
	for _, key := range keys {
		result, err := al.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}
		activities = append(activities, result)
	}

	return activities, nil
}

func (al *ActivityLogger) Close() error {
	return al.redisClient.Close()
}