package session

import (
	"crypto/rand"
	"encoding/base64"
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

func (m *Manager) CreateSession(userID string, duration time.Duration) (Session, error) {
	token, err := generateToken()
	if err != nil {
		return Session{}, err
	}

	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
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
	return base64.URLEncoding.EncodeToString(bytes), nil
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

    session := &Session{
        ID:        generateSessionID(),
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(sm.ttl),
    }
    sm.sessions[session.ID] = session
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
        sm.mu.Lock()
        now := time.Now()
        for id, session := range sm.sessions {
            if now.After(session.ExpiresAt) {
                delete(sm.sessions, id)
            }
        }
        sm.mu.Unlock()
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
    sessions map[string]Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
    return &Manager{
        sessions: make(map[string]Session),
        duration: sessionDuration,
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

func (m *Manager) Create(userID string) (Session, error) {
    token, err := generateToken()
    if err != nil {
        return Session{}, err
    }

    session := Session{
        Token:     token,
        UserID:    userID,
        ExpiresAt: time.Now().Add(m.duration),
        Data:      make(map[string]interface{}),
    }

    m.mu.Lock()
    m.sessions[token] = session
    m.mu.Unlock()

    return session, nil
}

func (m *Manager) Get(token string) (Session, bool) {
    m.mu.RLock()
    session, exists := m.sessions[token]
    m.mu.RUnlock()

    if !exists {
        return Session{}, false
    }

    if time.Now().After(session.ExpiresAt) {
        m.Delete(token)
        return Session{}, false
    }

    return session, true
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

func (m *Manager) Extend(token string, duration time.Duration) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    session, exists := m.sessions[token]
    if !exists {
        return false
    }

    session.ExpiresAt = time.Now().Add(duration)
    m.sessions[token] = session
    return true
}