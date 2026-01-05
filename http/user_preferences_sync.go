
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme      string `json:"theme"`
    Language   string `json:"language"`
    Notifications bool `json:"notifications"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type PreferenceManager struct {
    preferences UserPreferences
    filePath    string
    mu          sync.RWMutex
    syncTicker  *time.Ticker
    stopChan    chan bool
}

func NewPreferenceManager(configDir string) *PreferenceManager {
    filePath := filepath.Join(configDir, "preferences.json")
    pm := &PreferenceManager{
        filePath: filePath,
        stopChan: make(chan bool),
    }
    pm.loadPreferences()
    return pm
}

func (pm *PreferenceManager) loadPreferences() {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    data, err := os.ReadFile(pm.filePath)
    if err != nil {
        pm.preferences = UserPreferences{
            Theme:         "light",
            Language:      "en",
            Notifications: true,
            UpdatedAt:     time.Now(),
        }
        pm.savePreferences()
        return
    }

    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        log.Printf("Failed to parse preferences: %v", err)
        pm.preferences = UserPreferences{
            Theme:         "light",
            Language:      "en",
            Notifications: true,
            UpdatedAt:     time.Now(),
        }
    } else {
        pm.preferences = prefs
    }
}

func (pm *PreferenceManager) savePreferences() {
    data, err := json.MarshalIndent(pm.preferences, "", "  ")
    if err != nil {
        log.Printf("Failed to marshal preferences: %v", err)
        return
    }

    if err := os.WriteFile(pm.filePath, data, 0644); err != nil {
        log.Printf("Failed to write preferences: %v", err)
    }
}

func (pm *PreferenceManager) UpdatePreferences(updater func(*UserPreferences)) {
    pm.mu.Lock()
    updater(&pm.preferences)
    pm.preferences.UpdatedAt = time.Now()
    pm.mu.Unlock()

    go pm.savePreferences()
}

func (pm *PreferenceManager) GetPreferences() UserPreferences {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    return pm.preferences
}

func (pm *PreferenceManager) StartSync(interval time.Duration) {
    pm.syncTicker = time.NewTicker(interval)
    go func() {
        for {
            select {
            case <-pm.syncTicker.C:
                pm.mu.Lock()
                if time.Since(pm.preferences.UpdatedAt) > time.Hour {
                    pm.preferences.Theme = "auto"
                }
                pm.mu.Unlock()
                pm.savePreferences()
            case <-pm.stopChan:
                pm.syncTicker.Stop()
                return
            }
        }
    }()
}

func (pm *PreferenceManager) StopSync() {
    close(pm.stopChan)
}

func main() {
    configDir, err := os.UserConfigDir()
    if err != nil {
        configDir = "."
    }

    pm := NewPreferenceManager(configDir)
    defer pm.StopSync()

    pm.StartSync(5 * time.Minute)

    fmt.Printf("Current preferences: %+v\n", pm.GetPreferences())

    pm.UpdatePreferences(func(p *UserPreferences) {
        p.Theme = "dark"
        p.Language = "fr"
    })

    time.Sleep(2 * time.Second)
    fmt.Printf("Updated preferences: %+v\n", pm.GetPreferences())
}