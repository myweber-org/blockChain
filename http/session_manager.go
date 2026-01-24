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

func GenerateToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func CreateSession(userID int, duration time.Duration) (string, error) {
    token, err := GenerateToken()
    if err != nil {
        return "", err
    }

    session := Session{
        ID:        token,
        UserID:    userID,
        ExpiresAt: time.Now().Add(duration),
    }

    sessions[token] = session
    return token, nil
}

func ValidateSession(token string) (Session, error) {
    session, exists := sessions[token]
    if !exists {
        return Session{}, errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        delete(sessions, token)
        return Session{}, errors.New("session expired")
    }

    return session, nil
}

func InvalidateSession(token string) {
    delete(sessions, token)
}

func CleanupExpiredSessions() {
    now := time.Now()
    for token, session := range sessions {
        if now.After(session.ExpiresAt) {
            delete(sessions, token)
        }
    }
}