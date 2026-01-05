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
	ExpiresAt time.Time
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]*Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

func (m *Manager) CreateSession(userID int, ttl time.Duration) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
		Data:      make(map[string]interface{}),
	}

	m.sessions[token] = session
	return session, nil
}

func (m *Manager) ValidateSession(token string) (*Session, error) {
	session, exists := m.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return nil, errors.New("session expired")
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
	return base64.URLEncoding.EncodeToString(bytes), nil
}package session

import (
	"sync"
	"time"
)

type Session struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
	go m.cleanupWorker()
	return m
}

func (m *Manager) Create(id string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		ID:        id,
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(m.ttl),
	}
	m.sessions[id] = session
	return session
}

func (m *Manager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[id]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, false
	}
	return session, true
}

func (m *Manager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, id)
}

func (m *Manager) cleanupWorker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for id, session := range m.sessions {
			if now.After(session.ExpiresAt) {
				delete(m.sessions, id)
			}
		}
		m.mu.Unlock()
	}
}