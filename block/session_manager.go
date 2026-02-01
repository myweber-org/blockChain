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
    duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
    return &Manager{
        sessions: make(map[string]*Session),
        duration: sessionDuration,
    }
}

func (m *Manager) Create(userID string) (*Session, error) {
    token, err := generateToken()
    if err != nil {
        return nil, err
    }

    session := &Session{
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

func (m *Manager) Validate(token string) (*Session, error) {
    m.mu.RLock()
    session, exists := m.sessions[token]
    m.mu.RUnlock()

    if !exists {
        return nil, errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        m.Delete(token)
        return nil, errors.New("session expired")
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

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}