
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
}package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()
		maxScore := float64(now - 86400) // 24 hours ago

		// Remove expired sessions from sorted set
		removed, err := rdb.ZRemRangeByScore(ctx, "user_sessions", "0", string(maxScore)).Result()
		if err != nil {
			log.Printf("Failed to clean sessions: %v", err)
			continue
		}

		if removed > 0 {
			log.Printf("Cleaned %d expired sessions", removed)
		}
	}
}package main

import (
	"context"
	"log"
	"time"

	"github.com/yourproject/internal/database"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	cleanupInterval := 24 * time.Hour

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := cleanupExpiredSessions(ctx, db)
			if err != nil {
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