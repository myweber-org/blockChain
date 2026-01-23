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
}package main

import (
    "log"
    "time"
)

type Session struct {
    ID        string
    UserID    string
    ExpiresAt time.Time
}

type SessionStore interface {
    GetExpiredSessions() ([]Session, error)
    DeleteSession(id string) error
}

type SessionCleaner struct {
    store SessionStore
}

func NewSessionCleaner(store SessionStore) *SessionCleaner {
    return &SessionCleaner{store: store}
}

func (sc *SessionCleaner) CleanExpiredSessions() error {
    expiredSessions, err := sc.store.GetExpiredSessions()
    if err != nil {
        return err
    }

    for _, session := range expiredSessions {
        err := sc.store.DeleteSession(session.ID)
        if err != nil {
            log.Printf("Failed to delete session %s: %v", session.ID, err)
            continue
        }
        log.Printf("Deleted expired session: %s", session.ID)
    }

    return nil
}

func (sc *SessionCleaner) StartDailyCleanup() {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        err := sc.CleanExpiredSessions()
        if err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}

func main() {
    // Implementation would provide actual SessionStore
    var store SessionStore
    cleaner := NewSessionCleaner(store)
    cleaner.StartDailyCleanup()
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

type SessionStore interface {
	GetExpiredSessions() ([]Session, error)
	DeleteSession(id string) error
}

func cleanupExpiredSessions(store SessionStore) error {
	expiredSessions, err := store.GetExpiredSessions()
	if err != nil {
		return err
	}

	for _, session := range expiredSessions {
		err := store.DeleteSession(session.ID)
		if err != nil {
			log.Printf("Failed to delete session %s: %v", session.ID, err)
		} else {
			log.Printf("Deleted expired session %s for user %s", session.ID, session.UserID)
		}
	}

	log.Printf("Cleanup completed. Removed %d expired sessions.", len(expiredSessions))
	return nil
}

func scheduleCleanup(store SessionStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := cleanupExpiredSessions(store)
		if err != nil {
			log.Printf("Session cleanup failed: %v", err)
		}
	}
}

func main() {
	// Implementation would provide actual SessionStore
	var store SessionStore
	scheduleCleanup(store, 24*time.Hour)
}
package main

import (
	"context"
	"log"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionStore interface {
	DeleteExpiredSessions(ctx context.Context) error
}

type SessionCleanupJob struct {
	store     SessionStore
	interval  time.Duration
}

func NewSessionCleanupJob(store SessionStore, interval time.Duration) *SessionCleanupJob {
	return &SessionCleanupJob{
		store:    store,
		interval: interval,
	}
}

func (j *SessionCleanupJob) Run(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Session cleanup job stopped")
			return
		case <-ticker.C:
			if err := j.store.DeleteExpiredSessions(ctx); err != nil {
				log.Printf("Failed to delete expired sessions: %v", err)
			} else {
				log.Println("Expired sessions cleaned up successfully")
			}
		}
	}
}

func main() {
	ctx := context.Background()
	store := &mockSessionStore{}
	job := NewSessionCleanupJob(store, 24*time.Hour)

	log.Println("Starting session cleanup job...")
	job.Run(ctx)
}

type mockSessionStore struct{}

func (m *mockSessionStore) DeleteExpiredSessions(ctx context.Context) error {
	return nil
}