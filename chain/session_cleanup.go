
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

func main() {
	ctx := context.Background()
	database := db.GetDB()

	// Run cleanup daily at 2 AM
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(ctx, database)
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *db.Database) {
	cutoff := time.Now().Add(-24 * time.Hour)

	result := db.WithContext(ctx).
		Where("last_activity < ?", cutoff).
		Delete(&models.Session{})

	if result.Error != nil {
		log.Printf("Error cleaning up sessions: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions", result.RowsAffected)
	}
}package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
)

func cleanupExpiredSessions(db *sql.DB) {
    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    cutoff := time.Now().Add(-sessionTTL)

    result, err := db.Exec(query, cutoff)
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned %d expired sessions", rows)
}

func startSessionCleaner(ctx context.Context, db *sql.DB) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            cleanupExpiredSessions(db)
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    startSessionCleaner(ctx, db)
}