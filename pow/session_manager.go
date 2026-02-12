package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

var (
    ErrInvalidToken = errors.New("invalid session token")
    ErrSessionExpired = errors.New("session has expired")
)

type Session struct {
    UserID    string
    Email     string
    CreatedAt time.Time
    ExpiresAt time.Time
}

type Manager struct {
    client *redis.Client
    prefix string
    ttl    time.Duration
}

func NewManager(client *redis.Client, prefix string, ttl time.Duration) *Manager {
    return &Manager{
        client: client,
        prefix: prefix,
        ttl:    ttl,
    }
}

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) Create(userID, email string) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    session := Session{
        UserID:    userID,
        Email:     email,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(m.ttl),
    }

    key := m.prefix + token
    ctx := context.Background()
    
    err = m.client.Set(ctx, key, session, m.ttl).Err()
    if err != nil {
        return "", err
    }

    return token, nil
}

func (m *Manager) Validate(token string) (*Session, error) {
    key := m.prefix + token
    ctx := context.Background()

    var session Session
    err := m.client.Get(ctx, key).Scan(&session)
    if err != nil {
        if err == redis.Nil {
            return nil, ErrInvalidToken
        }
        return nil, err
    }

    if time.Now().After(session.ExpiresAt) {
        m.client.Del(ctx, key)
        return nil, ErrSessionExpired
    }

    return &session, nil
}

func (m *Manager) Refresh(token string) error {
    key := m.prefix + token
    ctx := context.Background()

    exists, err := m.client.Exists(ctx, key).Result()
    if err != nil {
        return err
    }
    if exists == 0 {
        return ErrInvalidToken
    }

    return m.client.Expire(ctx, key, m.ttl).Err()
}

func (m *Manager) Invalidate(token string) error {
    key := m.prefix + token
    ctx := context.Background()
    return m.client.Del(ctx, key).Err()
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