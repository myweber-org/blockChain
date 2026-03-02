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
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

const cleanupInterval = 24 * time.Hour

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(db)
		}
	}
}

func cleanupExpiredSessions(db *database.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to cleanup sessions: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rowsAffected)
}package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/database"
)

const cleanupInterval = 24 * time.Hour
const sessionTTL = 30 * 24 * time.Hour

func cleanupExpiredSessions(ctx context.Context) error {
    db := database.GetConnection()
    cutoff := time.Now().Add(-sessionTTL)

    result, err := db.ExecContext(ctx,
        "DELETE FROM user_sessions WHERE last_activity < ?",
        cutoff)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rows)
    return nil
}

func startCleanupScheduler() {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            if err := cleanupExpiredSessions(ctx); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
            cancel()
        }
    }
}

func main() {
    log.Println("Starting session cleanup scheduler...")
    startCleanupScheduler()
}
package main

import (
	"log"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionStore struct {
	sessions map[string]Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]Session),
	}
}

func (s *SessionStore) CleanExpiredSessions() {
	now := time.Now()
	count := 0
	for id, session := range s.sessions {
		if session.ExpiresAt.Before(now) {
			delete(s.sessions, id)
			count++
		}
	}
	log.Printf("Cleaned %d expired sessions", count)
}

func startCleanupJob(store *SessionStore) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		store.CleanExpiredSessions()
	}
}

func main() {
	store := NewSessionStore()
	go startCleanupJob(store)

	select {}
}package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionCleaner struct {
	client *redis.Client
	prefix string
}

func NewSessionCleaner(addr string, prefix string) *SessionCleaner {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &SessionCleaner{
		client: rdb,
		prefix: prefix,
	}
}

func (sc *SessionCleaner) CleanExpiredSessions(ctx context.Context) error {
	iter := sc.client.Scan(ctx, 0, sc.prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		ttl, err := sc.client.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get TTL for key %s: %v", key, err)
			continue
		}
		if ttl < 0 {
			if err := sc.client.Del(ctx, key).Err(); err != nil {
				log.Printf("Failed to delete expired session %s: %v", key, err)
			} else {
				log.Printf("Removed expired session: %s", key)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("iteration error: %w", err)
	}
	return nil
}

func (sc *SessionCleaner) RunCleanupJob(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sc.CleanExpiredSessions(ctx); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	cleaner := NewSessionCleaner("localhost:6379", "session:")
	ctx := context.Background()

	go cleaner.RunCleanupJob(ctx, 5*time.Minute)

	<-ctx.Done()
}