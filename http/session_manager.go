package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}
}

func (m *Manager) Create(userID string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(m.ttl),
	}

	m.mu.Lock()
	m.sessions[token] = session
	m.mu.Unlock()

	return token, nil
}

func (m *Manager) Validate(token string) (string, error) {
	m.mu.RLock()
	session, exists := m.sessions[token]
	m.mu.RUnlock()

	if !exists {
		return "", errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		m.Delete(token)
		return "", errors.New("session expired")
	}

	return session.UserID, nil
}

func (m *Manager) Delete(token string) {
	m.mu.Lock()
	delete(m.sessions, token)
	m.mu.Unlock()
}

func (m *Manager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}package session

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

func NewManager(duration time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]Session),
		duration: duration,
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

func (m *Manager) Validate(token string) (Session, error) {
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

func (m *Manager) Invalidate(token string) {
	delete(m.sessions, token)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}