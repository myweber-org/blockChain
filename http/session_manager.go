
package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]Session),
	}
}

func (m *Manager) CreateSession(userID string) (Session, error) {
	token, err := generateToken()
	if err != nil {
		return Session{}, err
	}

	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	m.sessions[token] = session
	return session, nil
}

func (m *Manager) ValidateSession(token string) (Session, error) {
	session, exists := m.sessions[token]
	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) InvalidateSession(token string) {
	delete(m.sessions, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}