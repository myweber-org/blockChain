package main

import (
    "encoding/json"
    "errors"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme     string `json:"theme"`
    Language  string `json:"language"`
    Notify    bool   `json:"notifications_enabled"`
    Timezone  string `json:"timezone"`
    UpdatedAt time.Time
}

type PreferencesCache struct {
    mu      sync.RWMutex
    store   map[string]UserPreferences
    ttl     time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        store: make(map[string]UserPreferences),
        ttl:   ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    prefs, exists := c.store[userID]
    if !exists {
        return UserPreferences{}, false
    }
    if time.Since(prefs.UpdatedAt) > c.ttl {
        return UserPreferences{}, false
    }
    return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()
    prefs.UpdatedAt = time.Now()
    c.store[userID] = prefs
}

func LoadUserPreferences(userID string, cache *PreferencesCache) (UserPreferences, error) {
    if cache != nil {
        if cached, ok := cache.Get(userID); ok {
            return cached, nil
        }
    }

    prefs, err := fetchPreferencesFromStorage(userID)
    if err != nil {
        return UserPreferences{}, err
    }

    if err := validatePreferences(prefs); err != nil {
        return UserPreferences{}, err
    }

    if cache != nil {
        cache.Set(userID, prefs)
    }

    return prefs, nil
}

func fetchPreferencesFromStorage(userID string) (UserPreferences, error) {
    if userID == "" {
        return UserPreferences{}, errors.New("invalid user identifier")
    }

    simulatedData := `{"theme":"dark","language":"en","notifications_enabled":true,"timezone":"UTC"}`
    var prefs UserPreferences
    if err := json.Unmarshal([]byte(simulatedData), &prefs); err != nil {
        return UserPreferences{}, err
    }
    return prefs, nil
}

func validatePreferences(prefs UserPreferences) error {
    if prefs.Theme != "light" && prefs.Theme != "dark" && prefs.Theme != "auto" {
        return errors.New("invalid theme selection")
    }
    if prefs.Language == "" {
        return errors.New("language must be specified")
    }
    return nil
}package main

import (
    "encoding/json"
    "errors"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme     string `json:"theme"`
    Language  string `json:"language"`
    Notify    bool   `json:"notifications_enabled"`
    Timezone  string `json:"timezone"`
}

type PreferencesCache struct {
    mu      sync.RWMutex
    store   map[string]UserPreferences
    ttl     time.Duration
    created map[string]time.Time
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        store:   make(map[string]UserPreferences),
        created: make(map[string]time.Time),
        ttl:     ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    prefs, exists := c.store[userID]
    if !exists {
        return UserPreferences{}, false
    }

    if time.Since(c.created[userID]) > c.ttl {
        return UserPreferences{}, false
    }

    return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.store[userID] = prefs
    c.created[userID] = time.Now()
}

func LoadPreferencesFromJSON(data []byte) (UserPreferences, error) {
    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        return UserPreferences{}, err
    }

    if prefs.Theme == "" {
        prefs.Theme = "light"
    }

    if prefs.Language == "" {
        prefs.Language = "en"
    }

    if prefs.Timezone == "" {
        return UserPreferences{}, errors.New("timezone is required")
    }

    return prefs, nil
}

func LoadUserPreferences(userID string, cache *PreferencesCache, fetchFunc func(string) ([]byte, error)) (UserPreferences, error) {
    if cached, ok := cache.Get(userID); ok {
        return cached, nil
    }

    data, err := fetchFunc(userID)
    if err != nil {
        return UserPreferences{}, err
    }

    prefs, err := LoadPreferencesFromJSON(data)
    if err != nil {
        return UserPreferences{}, err
    }

    cache.Set(userID, prefs)
    return prefs, nil
}