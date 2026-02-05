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
	Timezone  string `json:"timezone"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

type PreferencesCache struct {
	mu      sync.RWMutex
	store   map[string]UserPreferences
	expires map[string]time.Time
	ttl     time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		store:   make(map[string]UserPreferences),
		expires: make(map[string]time.Time),
		ttl:     ttl,
	}
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) error {
	if err := validatePreferences(prefs); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.store[userID] = prefs
	c.expires[userID] = time.Now().Add(c.ttl)
	return nil
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	c.mu.RLock()
	prefs, exists := c.store[userID]
	expiry, expExists := c.expires[userID]
	c.mu.RUnlock()

	if !exists || !expExists || time.Now().After(expiry) {
		if exists {
			go c.evict(userID)
		}
		return UserPreferences{}, false
	}
	return prefs, true
}

func (c *PreferencesCache) evict(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, userID)
	delete(c.expires, userID)
}

func validatePreferences(prefs UserPreferences) error {
	validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
	validLanguages := map[string]bool{"en": true, "es": true, "fr": true, "de": true}

	if !validThemes[prefs.Theme] {
		return errors.New("invalid theme selection")
	}
	if !validLanguages[prefs.Language] {
		return errors.New("invalid language code")
	}
	if prefs.Timezone == "" {
		return errors.New("timezone cannot be empty")
	}
	return nil
}

func SerializePreferences(prefs UserPreferences) ([]byte, error) {
	return json.MarshalIndent(prefs, "", "  ")
}

func DeserializePreferences(data []byte) (UserPreferences, error) {
	var prefs UserPreferences
	err := json.Unmarshal(data, &prefs)
	if err != nil {
		return UserPreferences{}, err
	}
	return prefs, validatePreferences(prefs)
}