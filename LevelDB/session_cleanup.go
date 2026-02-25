package main

import (
    "context"
    "log"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
    deleteBatchSize = 1000
)

type SessionCleaner struct {
    db *sqlx.DB
}

func NewSessionCleaner(db *sqlx.DB) *SessionCleaner {
    return &SessionCleaner{db: db}
}

func (sc *SessionCleaner) Run(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleaner stopped")
            return
        case <-ticker.C:
            sc.cleanupExpiredSessions()
        }
    }
}

func (sc *SessionCleaner) cleanupExpiredSessions() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cutoffTime := time.Now().Add(-sessionTTL)
    deletedCount := 0

    for {
        result, err := sc.db.ExecContext(ctx,
            "DELETE FROM user_sessions WHERE last_activity < $1 LIMIT $2",
            cutoffTime, deleteBatchSize)
        if err != nil {
            log.Printf("Failed to delete expired sessions: %v", err)
            return
        }

        rowsAffected, _ := result.RowsAffected()
        deletedCount += int(rowsAffected)

        if rowsAffected < deleteBatchSize {
            break
        }

        time.Sleep(100 * time.Millisecond)
    }

    if deletedCount > 0 {
        log.Printf("Cleaned up %d expired sessions", deletedCount)
    }
}

func main() {
    db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=app dbname=appdb password=secret sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    cleaner := NewSessionCleaner(db)
    ctx := context.Background()
    cleaner.Run(ctx)
}