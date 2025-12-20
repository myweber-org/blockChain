package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type UserPreferences struct {
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

type PreferencesCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPreferencesCache(addr string, ttlMinutes int) *PreferencesCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	
	return &PreferencesCache{
		client: rdb,
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}
}

func (c *PreferencesCache) Get(ctx context.Context, userID string) (*UserPreferences, error) {
	key := "preferences:" + userID
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var prefs UserPreferences
	err = json.Unmarshal([]byte(data), &prefs)
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func (c *PreferencesCache) Set(ctx context.Context, userID string, prefs *UserPreferences) error {
	key := "preferences:" + userID
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *PreferencesCache) Invalidate(ctx context.Context, userID string) error {
	key := "preferences:" + userID
	return c.client.Del(ctx, key).Err()
}