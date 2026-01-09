
package main

import (
	"log"
	"time"
)

type SessionStore struct {
	sessions map[string]time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]time.Time),
	}
}

func (s *SessionStore) AddSession(sessionID string) {
	s.sessions[sessionID] = time.Now()
}

func (s *SessionStore) IsValidSession(sessionID string) bool {
	created, exists := s.sessions[sessionID]
	if !exists {
		return false
	}
	return time.Since(created) < 24*time.Hour
}

func (s *SessionStore) CleanupExpiredSessions() {
	expirationTime := 24 * time.Hour
	now := time.Now()
	count := 0

	for sessionID, created := range s.sessions {
		if now.Sub(created) > expirationTime {
			delete(s.sessions, sessionID)
			count++
		}
	}

	log.Printf("Cleaned up %d expired sessions", count)
}

func startCleanupCron(store *SessionStore) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			store.CleanupExpiredSessions()
		}
	}
}

func main() {
	sessionStore := NewSessionStore()

	sessionStore.AddSession("abc123")
	sessionStore.AddSession("def456")

	log.Println("Session 'abc123' valid:", sessionStore.IsValidSession("abc123"))
	log.Println("Session 'xyz789' valid:", sessionStore.IsValidSession("xyz789"))

	go startCleanupCron(sessionStore)

	select {}
}package main

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions(ctx, rdb)
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, rdb *redis.Client) {
    pattern := "session:*"
    var cursor uint64
    var keys []string
    var err error

    for {
        keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 100).Result()
        if err != nil {
            log.Printf("Error scanning keys: %v", err)
            return
        }

        for _, key := range keys {
            ttl, err := rdb.TTL(ctx, key).Result()
            if err != nil {
                log.Printf("Error getting TTL for key %s: %v", key, err)
                continue
            }
            if ttl == -2 || ttl == -1 {
                if err := rdb.Del(ctx, key).Err(); err != nil {
                    log.Printf("Error deleting expired session %s: %v", key, err)
                } else {
                    log.Printf("Cleaned up expired session: %s", key)
                }
            }
        }

        if cursor == 0 {
            break
        }
    }
}