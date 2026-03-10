package session

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    ExpiresAt time.Time
}

var sessions = make(map[string]Session)

func GenerateSession(userID int) (string, error) {
    tokenBytes := make([]byte, 16)
    if _, err := rand.Read(tokenBytes); err != nil {
        return "", err
    }
    token := hex.EncodeToString(tokenBytes)

    sessions[token] = Session{
        ID:        token,
        UserID:    userID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    return token, nil
}

func ValidateSession(token string) (int, error) {
    session, exists := sessions[token]
    if !exists {
        return 0, errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        delete(sessions, token)
        return 0, errors.New("session expired")
    }

    return session.UserID, nil
}

func InvalidateSession(token string) {
    delete(sessions, token)
}package main

import (
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	Data      map[string]interface{}
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
	go sm.cleanupWorker()
	return sm
}

func (sm *SessionManager) CreateSession(userID int) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(sm.ttl),
	}
	sm.sessions[sessionID] = session
	return session
}

func (sm *SessionManager) GetSession(sessionID string) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

func (sm *SessionManager) RefreshSession(sessionID string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return false
	}
	session.ExpiresAt = time.Now().Add(sm.ttl)
	return true
}

func (sm *SessionManager) cleanupWorker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.cleanupExpiredSessions()
	}
}

func (sm *SessionManager) cleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for id, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, id)
		}
	}
}

func generateSessionID() string {
	return "session_" + time.Now().Format("20060102150405") + "_" + randomString(16)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
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
	ExpiresAt time.Time
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]*Session
	duration time.Duration
}

func NewManager(duration time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		duration: duration,
	}
}

func (m *Manager) Create(userID int, data map[string]interface{}) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.duration),
		Data:      data,
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

	return session, nil
}

func (m *Manager) Invalidate(token string) {
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
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
package session

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
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
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

func (m *Manager) Create(userID string, data map[string]interface{}) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.ttl),
		Data:      data,
	}

	m.mu.Lock()
	m.sessions[token] = session
	m.mu.Unlock()

	return session, nil
}

func (m *Manager) Validate(token string) (*Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[token]
	m.mu.RUnlock()

	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		m.Delete(token)
		return nil, ErrSessionExpired
	}

	return session, nil
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

func (m *Manager) Refresh(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[token]
	if !exists {
		return ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return ErrSessionExpired
	}

	session.ExpiresAt = time.Now().Add(m.ttl)
	return nil
}package main

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    ttl      time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
    sm := &SessionManager{
        sessions: make(map[string]*Session),
        ttl:      ttl,
    }
    go sm.cleanupLoop()
    return sm
}

func (sm *SessionManager) CreateSession(userID int, data map[string]interface{}) string {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        Data:      data,
        ExpiresAt: time.Now().Add(sm.ttl),
    }
    sm.sessions[sessionID] = session
    return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) *Session {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    session, exists := sm.sessions[sessionID]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
}

func (sm *SessionManager) DeleteSession(sessionID string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    delete(sm.sessions, sessionID)
}

func (sm *SessionManager) cleanupLoop() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        sm.cleanupExpired()
    }
}

func (sm *SessionManager) cleanupExpired() {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    now := time.Now()
    for id, session := range sm.sessions {
        if now.After(session.ExpiresAt) {
            delete(sm.sessions, id)
        }
    }
}

func generateSessionID() string {
    return "sess_" + time.Now().Format("20060102150405") + "_" + randomString(16)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}