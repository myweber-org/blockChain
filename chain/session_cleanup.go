package main

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

const (
    sessionKeyPattern = "session:*"
    sessionTTL        = 24 * time.Hour
)

func cleanupExpiredSessions(client *redis.Client) error {
    ctx := context.Background()
    iter := client.Scan(ctx, 0, sessionKeyPattern, 0).Iterator()
    
    var deletedCount int
    for iter.Next(ctx) {
        key := iter.Val()
        ttl, err := client.TTL(ctx, key).Result()
        if err != nil {
            continue
        }
        
        if ttl < 0 {
            if err := client.Del(ctx, key).Err(); err == nil {
                deletedCount++
            }
        }
    }
    
    if err := iter.Err(); err != nil {
        return err
    }
    
    fmt.Printf("Cleaned up %d expired sessions\n", deletedCount)
    return nil
}

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(client); err != nil {
            fmt.Printf("Cleanup error: %v\n", err)
        }
    }
}