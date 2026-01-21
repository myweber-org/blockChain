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
}