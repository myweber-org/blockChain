
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserPreferences struct {
	UserID    string `json:"user_id"`
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	UpdatedAt int64  `json:"updated_at"`
}

type PreferencesManager struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewPreferencesManager(addr string, ttl time.Duration) *PreferencesManager {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &PreferencesManager{
		redisClient: client,
		ttl:         ttl,
	}
}

func (pm *PreferencesManager) SetPreferences(ctx context.Context, prefs *UserPreferences) error {
	prefs.UpdatedAt = time.Now().Unix()
	data, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("marshal preferences failed: %w", err)
	}

	key := fmt.Sprintf("prefs:%s", prefs.UserID)
	err = pm.redisClient.Set(ctx, key, data, pm.ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}
	return nil
}

func (pm *PreferencesManager) GetPreferences(ctx context.Context, userID string) (*UserPreferences, error) {
	key := fmt.Sprintf("prefs:%s", userID)
	data, err := pm.redisClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get failed: %w", err)
	}

	var prefs UserPreferences
	err = json.Unmarshal(data, &prefs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal preferences failed: %w", err)
	}
	return &prefs, nil
}

func (pm *PreferencesManager) DeletePreferences(ctx context.Context, userID string) error {
	key := fmt.Sprintf("prefs:%s", userID)
	err := pm.redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete failed: %w", err)
	}
	return nil
}

func (pm *PreferencesManager) Close() error {
	return pm.redisClient.Close()
}