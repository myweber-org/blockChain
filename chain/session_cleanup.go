
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
}