package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidToken    = errors.New("invalid session token")
)

type Session struct {
	UserID    string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	client     *redis.Client
	expiration time.Duration
	prefix     string
}

func NewManager(client *redis.Client, expiration time.Duration) *Manager {
	return &Manager{
		client:     client,
		expiration: expiration,
		prefix:     "session:",
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) Create(userID, username string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.expiration),
	}

	ctx := context.Background()
	key := m.prefix + token

	err = m.client.Set(ctx, key, session, m.expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return token, nil
}

func (m *Manager) Get(token string) (*Session, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	ctx := context.Background()
	key := m.prefix + token

	var session Session
	err := m.client.Get(ctx, key).Scan(&session)
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	return &session, nil
}

func (m *Manager) Delete(token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	ctx := context.Background()
	key := m.prefix + token

	err := m.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (m *Manager) Refresh(token string) error {
	session, err := m.Get(token)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(m.expiration)

	ctx := context.Background()
	key := m.prefix + token

	err = m.client.Set(ctx, key, session, m.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

func (m *Manager) Cleanup() error {
	ctx := context.Background()
	iter := m.client.Scan(ctx, 0, m.prefix+"*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		var session Session
		err := m.client.Get(ctx, key).Scan(&session)
		if err == nil && time.Now().After(session.ExpiresAt) {
			m.client.Del(ctx, key)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed during cleanup: %w", err)
	}

	return nil
}