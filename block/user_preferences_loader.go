package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type UserPreferences struct {
	Theme      string   `json:"theme"`
	Language   string   `json:"language"`
	Timezone   string   `json:"timezone"`
	EmailAlerts bool    `json:"email_alerts"`
	Categories []string `json:"categories"`
}

type PreferencesCache struct {
	mu          sync.RWMutex
	preferences map[string]UserPreferences
	expiry      map[string]time.Time
	ttl         time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		preferences: make(map[string]UserPreferences),
		expiry:      make(map[string]time.Time),
		ttl:         ttl,
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	prefs, exists := c.preferences[userID]
	if !exists {
		return UserPreferences{}, false
	}

	expiryTime, expiryExists := c.expiry[userID]
	if !expiryExists || time.Now().After(expiryTime) {
		return UserPreferences{}, false
	}

	return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.preferences[userID] = prefs
	c.expiry[userID] = time.Now().Add(c.ttl)
}

func ValidatePreferences(prefs UserPreferences) error {
	if prefs.Theme == "" {
		return errors.New("theme cannot be empty")
	}
	if prefs.Language == "" {
		return errors.New("language cannot be empty")
	}
	if prefs.Timezone == "" {
		return errors.New("timezone cannot be empty")
	}
	if len(prefs.Categories) > 20 {
		return errors.New("too many categories")
	}
	return nil
}

func LoadPreferencesFromSource(userID string) (UserPreferences, error) {
	simulatedData := fmt.Sprintf(`{
		"theme": "dark",
		"language": "en",
		"timezone": "UTC",
		"email_alerts": true,
		"categories": ["tech", "sports", "news"]
	}`)

	var prefs UserPreferences
	err := json.Unmarshal([]byte(simulatedData), &prefs)
	if err != nil {
		return UserPreferences{}, fmt.Errorf("failed to parse preferences: %w", err)
	}

	prefs.Theme = "dark"
	prefs.Language = "en"
	prefs.Timezone = "UTC"
	prefs.EmailAlerts = true
	prefs.Categories = []string{"tech", "sports", "news"}

	return prefs, nil
}

func LoadUserPreferences(userID string, cache *PreferencesCache) (UserPreferences, error) {
	if cachedPrefs, found := cache.Get(userID); found {
		return cachedPrefs, nil
	}

	prefs, err := LoadPreferencesFromSource(userID)
	if err != nil {
		return UserPreferences{}, fmt.Errorf("failed to load preferences: %w", err)
	}

	if err := ValidatePreferences(prefs); err != nil {
		return UserPreferences{}, fmt.Errorf("invalid preferences: %w", err)
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

	cachedPrefs, found := cache.Get("user123")
	if found {
		fmt.Printf("Retrieved from cache: %+v\n", cachedPrefs)
	}
}