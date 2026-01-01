package session

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
	go m.cleanupLoop()
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

func (m *Manager) Get(id string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[id]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

func (m *Manager) cleanupLoop() {
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
	duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		duration: sessionDuration,
	}
}

func (m *Manager) Create(userID int, initialData map[string]interface{}) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.duration),
		Data:      initialData,
	}

	m.sessions[token] = session
	return token, nil
}

func (m *Manager) Validate(token string) (*Session, error) {
	session, exists := m.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return nil, errors.New("session expired")
	}

	session.ExpiresAt = time.Now().Add(m.duration)
	return session, nil
}

func (m *Manager) UpdateData(token string, key string, value interface{}) error {
	session, err := m.Validate(token)
	if err != nil {
		return err
	}

	session.Data[key] = value
	return nil
}

func (m *Manager) Destroy(token string) {
	delete(m.sessions, token)
}

func (m *Manager) Cleanup() {
	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}