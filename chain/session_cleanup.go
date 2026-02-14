package main

import (
    "context"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

const (
    dbURL          = "postgresql://localhost/sessiondb"
    cleanupTimeout = 10 * time.Second
    retentionDays  = 30
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
    defer cancel()

    pool, err := pgxpool.New(ctx, dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer pool.Close()

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    cutoff := time.Now().AddDate(0, 0, -retentionDays)

    result, err := pool.Exec(ctx, query, cutoff)
    if err != nil {
        log.Fatalf("Failed to delete expired sessions: %v\n", err)
    }

    rowsAffected := result.RowsAffected()
    log.Printf("Cleaned up %d expired session records\n", rowsAffected)
}