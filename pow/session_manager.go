package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (sm *SessionManager) CreateSession(userID string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sm.ttl),
	}

	sm.mu.Lock()
	sm.sessions[token] = session
	sm.mu.Unlock()

	return token, nil
}

func (sm *SessionManager) ValidateSession(token string) (string, error) {
	sm.mu.RLock()
	session, exists := sm.sessions[token]
	sm.mu.RUnlock()

	if !exists {
		return "", errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		sm.mu.Lock()
		delete(sm.sessions, token)
		sm.mu.Unlock()
		return "", errors.New("session expired")
	}

	return session.UserID, nil
}

func (sm *SessionManager) InvalidateSession(token string) {
	sm.mu.Lock()
	delete(sm.sessions, token)
	sm.mu.Unlock()
}

func (sm *SessionManager) CleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for token, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, token)
		}
	}
}