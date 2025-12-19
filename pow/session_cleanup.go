
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
}