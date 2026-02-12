package main

import (
    "context"
    "log"
    "time"
)

type SessionStore interface {
    DeleteExpiredSessions(ctx context.Context) error
}

type CleanupJob struct {
    store SessionStore
}

func NewCleanupJob(store SessionStore) *CleanupJob {
    return &CleanupJob{store: store}
}

func (j *CleanupJob) Run() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    if err := j.store.DeleteExpiredSessions(ctx); err != nil {
        log.Printf("Failed to delete expired sessions: %v", err)
    } else {
        log.Println("Successfully cleaned up expired sessions")
    }
}

func ScheduleCleanup(job *CleanupJob, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        job.Run()
    }
}

func main() {
    store := NewMemorySessionStore()
    job := NewCleanupJob(store)
    
    go ScheduleCleanup(job, 24*time.Hour)
    
    select {}
}