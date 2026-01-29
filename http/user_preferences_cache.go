package cache

import (
	"encoding/json"
	"errors"
	"time"
)

type UserPreferences struct {
	UserID      string                 `json:"user_id"`
	Theme       string                 `json:"theme"`
	Language    string                 `json:"language"`
	Timezone    string                 `json:"timezone"`
	Preferences map[string]interface{} `json:"preferences"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type PreferencesCache struct {
	store    map[string]UserPreferences
	duration time.Duration
}

func NewPreferencesCache(duration time.Duration) *PreferencesCache {
	return &PreferencesCache{
		store:    make(map[string]UserPreferences),
		duration: duration,
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, error) {
	prefs, exists := c.store[userID]
	if !exists {
		return UserPreferences{}, errors.New("preferences not found in cache")
	}

	if time.Since(prefs.UpdatedAt) > c.duration {
		delete(c.store, userID)
		return UserPreferences{}, errors.New("cache expired")
	}

	return prefs, nil
}

func (c *PreferencesCache) Set(prefs UserPreferences) error {
	prefs.UpdatedAt = time.Now()
	c.store[prefs.UserID] = prefs
	return nil
}

func (c *PreferencesCache) Invalidate(userID string) {
	delete(c.store, userID)
}

func (c *PreferencesCache) Serialize(prefs UserPreferences) ([]byte, error) {
	return json.Marshal(prefs)
}

func (c *PreferencesCache) Deserialize(data []byte) (UserPreferences, error) {
	var prefs UserPreferences
	err := json.Unmarshal(data, &prefs)
	return prefs, err
}