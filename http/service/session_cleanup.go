
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

const cleanupInterval = 24 * time.Hour

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, db); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Cleaned up %d expired sessions", rowsAffected)
	return nil
}
package main

import (
    "context"
    "database/sql"
    "log"
    "time"

    _ "github.com/lib/pq"
)

const (
    dbConnectionString = "host=localhost port=5432 user=app dbname=appdb password=secret sslmode=disable"
    cleanupInterval    = 1 * time.Hour
    sessionTTL         = 24 * time.Hour
)

func main() {
    db, err := sql.Open("postgres", dbConnectionString)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    log.Printf("Session cleanup service started. Cleanup interval: %v", cleanupInterval)

    for range ticker.C {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Cleanup failed: %v", err)
        }
    }
}

func cleanupExpiredSessions(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cutoffTime := time.Now().Add(-sessionTTL)

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    result, err := db.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected > 0 {
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
    }

    return nil
}