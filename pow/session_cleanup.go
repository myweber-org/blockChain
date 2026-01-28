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

func cleanupExpiredSessions() {
    now := time.Now().Unix()
    cursor := uint64(0)
    pattern := "session:*"

    for {
        var keys []string
        var err error
        keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            log.Printf("Error scanning keys: %v", err)
            return
        }

        for _, key := range keys {
            exp, err := rdb.Get(ctx, key+"_expiry").Int64()
            if err != nil {
                continue
            }
            if exp < now {
                rdb.Del(ctx, key, key+"_expiry")
                log.Printf("Removed expired session: %s", key)
            }
        }

        if cursor == 0 {
            break
        }
    }
}

func main() {
    initRedis()
    cleanupExpiredSessions()
}package main

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

func cleanupExpiredSessions() {
    now := time.Now().Unix()
    cursor := uint64(0)
    pattern := "session:*"

    for {
        var keys []string
        var err error
        keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 10).Result()
        if err != nil {
            log.Printf("Scan error: %v", err)
            break
        }

        for _, key := range keys {
            exp, err := rdb.Get(ctx, key+":expires").Int64()
            if err != nil {
                continue
            }
            if exp < now {
                rdb.Del(ctx, key, key+":expires")
                log.Printf("Removed expired session: %s", key)
            }
        }

        if cursor == 0 {
            break
        }
    }
}

func main() {
    initRedis()
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions()
            log.Println("Session cleanup completed")
        }
    }
}