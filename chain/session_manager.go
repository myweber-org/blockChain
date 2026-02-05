package main

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
    timeout  time.Duration
}

func NewSessionManager(timeout time.Duration) *SessionManager {
    sm := &SessionManager{
        sessions: make(map[string]*Session),
        timeout:  timeout,
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
        ExpiresAt: time.Now().Add(sm.timeout),
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
    session.ExpiresAt = time.Now().Add(sm.timeout)
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
    return "session_" + time.Now().Format("20060102150405") + "_" + randomString(8)
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

func (m *Manager) CreateSession() *Session {
    m.mu.Lock()
    defer m.mu.Unlock()

    session := &Session{
        ID:        generateID(),
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.ttl),
    }
    m.sessions[session.ID] = session
    return session
}

func (m *Manager) GetSession(id string) *Session {
    m.mu.RLock()
    defer m.mu.RUnlock()

    session, exists := m.sessions[id]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
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

func generateID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(b)
}package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

var (
    ErrSessionNotFound = errors.New("session not found")
    ErrInvalidToken    = errors.New("invalid session token")
)

type Session struct {
    UserID    string                 `json:"user_id"`
    Data      map[string]interface{} `json:"data"`
    CreatedAt time.Time              `json:"created_at"`
    ExpiresAt time.Time              `json:"expires_at"`
}

type Manager struct {
    client    *redis.Client
    prefix    string
    expiry    time.Duration
}

func NewManager(client *redis.Client, prefix string, expiry time.Duration) *Manager {
    return &Manager{
        client: client,
        prefix: prefix,
        expiry: expiry,
    }
}

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) Create(userID string, data map[string]interface{}) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    session := Session{
        UserID:    userID,
        Data:      data,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(m.expiry),
    }

    key := m.prefix + token
    ctx := context.Background()
    
    if err := m.client.Set(ctx, key, session, m.expiry).Err(); err != nil {
        return "", err
    }

    return token, nil
}

func (m *Manager) Get(token string) (*Session, error) {
    key := m.prefix + token
    ctx := context.Background()

    var session Session
    if err := m.client.Get(ctx, key).Scan(&session); err != nil {
        if err == redis.Nil {
            return nil, ErrSessionNotFound
        }
        return nil, err
    }

    if time.Now().After(session.ExpiresAt) {
        m.Delete(token)
        return nil, ErrSessionNotFound
    }

    return &session, nil
}

func (m *Manager) Update(token string, data map[string]interface{}) error {
    session, err := m.Get(token)
    if err != nil {
        return err
    }

    session.Data = data
    key := m.prefix + token
    ctx := context.Background()

    remaining := time.Until(session.ExpiresAt)
    if remaining <= 0 {
        return ErrSessionNotFound
    }

    return m.client.Set(ctx, key, session, remaining).Err()
}

func (m *Manager) Delete(token string) error {
    key := m.prefix + token
    ctx := context.Background()
    return m.client.Del(ctx, key).Err()
}

func (m *Manager) Extend(token string, duration time.Duration) error {
    session, err := m.Get(token)
    if err != nil {
        return err
    }

    session.ExpiresAt = session.ExpiresAt.Add(duration)
    key := m.prefix + token
    ctx := context.Background()

    return m.client.Set(ctx, key, session, m.expiry+duration).Err()
}

func (m *Manager) Cleanup() error {
    ctx := context.Background()
    keys, err := m.client.Keys(ctx, m.prefix+"*").Result()
    if err != nil {
        return err
    }

    for _, key := range keys {
        var session Session
        if err := m.client.Get(ctx, key).Scan(&session); err != nil {
            continue
        }
        if time.Now().After(session.ExpiresAt) {
            m.client.Del(ctx, key)
        }
    }

    return nil
}