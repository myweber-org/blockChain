package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]Session
	duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]Session),
		duration: sessionDuration,
	}
}

func (m *Manager) Create(userID int) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	now := time.Now()
	session := Session{
		ID:        token,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(m.duration),
	}

	m.sessions[token] = session
	return token, nil
}

func (m *Manager) Validate(token string) (int, error) {
	session, exists := m.sessions[token]
	if !exists {
		return 0, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return 0, errors.New("session expired")
	}

	return session.UserID, nil
}

func (m *Manager) Invalidate(token string) {
	delete(m.sessions, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}