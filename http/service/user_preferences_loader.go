package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme      string   `json:"theme"`
    Language   string   `json:"language"`
    Notifications struct {
        Email bool `json:"email"`
        Push  bool `json:"push"`
    } `json:"notifications"`
    Timezone string `json:"timezone"`
}

type PreferencesCache struct {
    mu       sync.RWMutex
    store    map[string]UserPreferences
    lastLoad map[string]time.Time
    ttl      time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        store:    make(map[string]UserPreferences),
        lastLoad: make(map[string]time.Time),
        ttl:      ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    prefs, exists := c.store[userID]
    lastLoad, loaded := c.lastLoad[userID]
    c.mu.RUnlock()

    if !exists || !loaded || time.Since(lastLoad) > c.ttl {
        return UserPreferences{}, false
    }
    return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.store[userID] = prefs
    c.lastLoad[userID] = time.Now()
}

func validatePreferences(prefs UserPreferences) error {
    validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
    if !validThemes[prefs.Theme] {
        return errors.New("invalid theme value")
    }

    if prefs.Language == "" {
        return errors.New("language cannot be empty")
    }

    if prefs.Timezone == "" {
        return errors.New("timezone cannot be empty")
    }
    return nil
}

func loadPreferencesFromFile(userID string) (UserPreferences, error) {
    filename := fmt.Sprintf("prefs_%s.json", userID)
    data, err := os.ReadFile(filename)
    if err != nil {
        return UserPreferences{}, fmt.Errorf("failed to read preferences file: %w", err)
    }

    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        return UserPreferences{}, fmt.Errorf("failed to parse preferences: %w", err)
    }

    if err := validatePreferences(prefs); err != nil {
        return UserPreferences{}, fmt.Errorf("invalid preferences: %w", err)
    }

    return prefs, nil
}

func LoadUserPreferences(userID string, cache *PreferencesCache) (UserPreferences, error) {
    if cached, ok := cache.Get(userID); ok {
        return cached, nil
    }

    prefs, err := loadPreferencesFromFile(userID)
    if err != nil {
        return UserPreferences{}, err
    }

    cache.Set(userID, prefs)
    return prefs, nil
}

func main() {
    cache := NewPreferencesCache(5 * time.Minute)
    
    prefs, err := LoadUserPreferences("user123", cache)
    if err != nil {
        fmt.Printf("Error loading preferences: %v\n", err)
        return
    }
    
    fmt.Printf("Loaded preferences: %+v\n", prefs)
}package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type UserPreferences struct {
    Theme      string `json:"theme"`
    Language   string `json:"language"`
    Notify     bool   `json:"notifications_enabled"`
    ItemsPerPage int  `json:"items_per_page"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }

    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        return nil, fmt.Errorf("invalid JSON format: %w", err)
    }

    if prefs.Theme == "" {
        return nil, fmt.Errorf("required field 'theme' is missing")
    }
    if prefs.Language == "" {
        return nil, fmt.Errorf("required field 'language' is missing")
    }
    if prefs.ItemsPerPage <= 0 {
        return nil, fmt.Errorf("items_per_page must be positive")
    }

    return &prefs, nil
}

func main() {
    prefs, err := LoadPreferences("config.json")
    if err != nil {
        fmt.Printf("Error loading preferences: %v\n", err)
        return
    }
    fmt.Printf("Loaded preferences: %+v\n", prefs)
}