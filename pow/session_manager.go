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

	session := Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.duration),
		Data:      make(map[string]interface{}),
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

	return &session, nil
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
}package main

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    Data      map[string]interface{}
    CreatedAt time.Time
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
    now := time.Now()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        Data:      make(map[string]interface{}),
        CreatedAt: now,
        ExpiresAt: now.Add(sm.ttl),
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

func (sm *SessionManager) cleanupWorker() {
    ticker := time.NewTicker(5 * time.Minute)
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
    return "sess_" + time.Now().Format("20060102150405") + "_" + randomString(16)
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

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type SessionManager struct {
	client     *redis.Client
	expiration time.Duration
}

func NewSessionManager(redisAddr string, expiration time.Duration) (*SessionManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &SessionManager{
		client:     client,
		expiration: expiration,
	}, nil
}

func (sm *SessionManager) GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (sm *SessionManager) CreateSession(userID string) (string, error) {
	token, err := sm.GenerateToken()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	err = sm.client.Set(ctx, token, userID, sm.expiration).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (sm *SessionManager) ValidateSession(token string) (string, error) {
	ctx := context.Background()
	userID, err := sm.client.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("session not found")
		}
		return "", err
	}

	// Refresh expiration on validation
	sm.client.Expire(ctx, token, sm.expiration)

	return userID, nil
}

func (sm *SessionManager) DestroySession(token string) error {
	ctx := context.Background()
	return sm.client.Del(ctx, token).Err()
}

func (sm *SessionManager) CleanupExpiredSessions() error {
	// Redis handles expiration automatically, but this method
	// can be used for additional cleanup if needed
	return nil
}