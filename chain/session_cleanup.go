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
}package main

import (
    "context"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
)

func main() {
    dbURL := "postgresql://user:pass@localhost:5432/dbname"
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer pool.Close()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    log.Println("Session cleanup job started")
    for range ticker.C {
        cleanupExpiredSessions(pool)
    }
}

func cleanupExpiredSessions(pool *pgxpool.Pool) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cutoff := time.Now().Add(-sessionTTL)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`

    result, err := pool.Exec(ctx, query, cutoff)
    if err != nil {
        log.Printf("Failed to cleanup sessions: %v", err)
        return
    }

    rowsAffected := result.RowsAffected()
    if rowsAffected > 0 {
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
    }
}