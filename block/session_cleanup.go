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

func (s *SessionStore) CleanExpiredSessions(maxAge time.Duration) {
	now := time.Now()
	for sessionID, createdAt := range s.sessions {
		if now.Sub(createdAt) > maxAge {
			delete(s.sessions, sessionID)
			log.Printf("Removed expired session: %s", sessionID)
		}
	}
}

func sessionCleanupCron(store *SessionStore, interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		store.CleanExpiredSessions(maxAge)
	}
}

func main() {
	sessionStore := NewSessionStore()
	
	sessionStore.sessions["abc123"] = time.Now().Add(-48 * time.Hour)
	sessionStore.sessions["def456"] = time.Now().Add(-1 * time.Hour)
	sessionStore.sessions["ghi789"] = time.Now()

	go sessionCleanupCron(sessionStore, 24*time.Hour, 24*time.Hour)

	select {}
}package main

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

	ctx := context.Background()
	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, db); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
	return nil
}