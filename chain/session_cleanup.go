package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func cleanupExpiredSessions() error {
	now := time.Now().Unix()
	maxScore := float64(now - 86400)

	keys, err := rdb.ZRangeByScore(ctx, "user_sessions", &redis.ZRangeBy{
		Min: "-inf",
		Max: string(maxScore),
	}).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		_, err = rdb.ZRem(ctx, "user_sessions", keys).Result()
		if err != nil {
			return err
		}
		log.Printf("Removed %d expired sessions", len(keys))
	}
	return nil
}

func main() {
	initRedis()
	
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(); err != nil {
				log.Printf("Cleanup failed: %v", err)
			}
		}
	}
}