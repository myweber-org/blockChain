package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Path      string    `json:"path"`
    IPAddress string    `json:"ip_address"`
}

type LoggerConfig struct {
    OutputFile string
    LogToStdout bool
}

type ActivityLogger struct {
    config LoggerConfig
    file   *os.File
}

func NewActivityLogger(config LoggerConfig) (*ActivityLogger, error) {
    logger := &ActivityLogger{config: config}
    
    if config.OutputFile != "" {
        file, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return nil, err
        }
        logger.file = file
    }
    
    return logger, nil
}

func (al *ActivityLogger) LogActivity(userID, action, path, ipAddress string) {
    activity := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Path:      path,
        IPAddress: ipAddress,
    }
    
    logData, err := json.Marshal(activity)
    if err != nil {
        log.Printf("Failed to marshal activity log: %v", err)
        return
    }
    
    if al.config.LogToStdout {
        fmt.Printf("ACTIVITY: %s\n", string(logData))
    }
    
    if al.file != nil {
        if _, err := al.file.Write(append(logData, '\n')); err != nil {
            log.Printf("Failed to write to log file: %v", err)
        }
    }
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }
        
        al.LogActivity(
            userID,
            r.Method,
            r.URL.Path,
            r.RemoteAddr,
        )
        
        next.ServeHTTP(w, r)
    })
}

func (al *ActivityLogger) Close() error {
    if al.file != nil {
        return al.file.Close()
    }
    return nil
}

func main() {
    config := LoggerConfig{
        OutputFile: "activity.log",
        LogToStdout: true,
    }
    
    logger, err := NewActivityLogger(config)
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Activity logging example"))
    })
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: logger.Middleware(mux),
    }
    
    log.Println("Server starting on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}