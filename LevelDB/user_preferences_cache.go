
package cache

import (
	"sync"
	"time"
)

type UserPreference struct {
	Theme      string
	Language   string
	Timezone   string
	UpdatedAt  time.Time
}

type PreferencesCache struct {
	mu      sync.RWMutex
	store   map[string]UserPreference
	ttl     time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		store: make(map[string]UserPreference),
		ttl:   ttl,
	}
}

func (c *PreferencesCache) Set(userID string, pref UserPreference) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	pref.UpdatedAt = time.Now()
	c.store[userID] = pref
}

func (c *PreferencesCache) Get(userID string) (UserPreference, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	pref, exists := c.store[userID]
	if !exists {
		return UserPreference{}, false
	}
	
	if time.Since(pref.UpdatedAt) > c.ttl {
		return UserPreference{}, false
	}
	
	return pref, true
}

func (c *PreferencesCache) Delete(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, userID)
}

func (c *PreferencesCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now()
	for userID, pref := range c.store {
		if now.Sub(pref.UpdatedAt) > c.ttl {
			delete(c.store, userID)
		}
	}
}