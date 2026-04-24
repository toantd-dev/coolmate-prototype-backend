package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheManager struct {
	client *redis.Client
}

func NewCacheManager(client *redis.Client) *CacheManager {
	return &CacheManager{client: client}
}

func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if cm.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return cm.client.Set(ctx, key, data, expiration).Err()
}

func (cm *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	if cm.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	val, err := cm.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (cm *CacheManager) Delete(ctx context.Context, keys ...string) error {
	if cm.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	return cm.client.Del(ctx, keys...).Err()
}

func (cm *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	if cm.client == nil {
		return false, fmt.Errorf("cache client is nil")
	}
	val, err := cm.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (cm *CacheManager) SetString(ctx context.Context, key, value string, expiration time.Duration) error {
	if cm.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	return cm.client.Set(ctx, key, value, expiration).Err()
}

func (cm *CacheManager) GetString(ctx context.Context, key string) (string, error) {
	if cm.client == nil {
		return "", fmt.Errorf("cache client is nil")
	}
	return cm.client.Get(ctx, key).Result()
}

func (cm *CacheManager) SetHash(ctx context.Context, key string, values map[string]interface{}) error {
	if cm.client == nil {
		return fmt.Errorf("cache client is nil")
	}
	return cm.client.HSet(ctx, key, values).Err()
}

func (cm *CacheManager) GetHashAll(ctx context.Context, key string) (map[string]string, error) {
	if cm.client == nil {
		return nil, fmt.Errorf("cache client is nil")
	}
	return cm.client.HGetAll(ctx, key).Result()
}

func (cm *CacheManager) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	if cm.client == nil {
		return 0, fmt.Errorf("cache client is nil")
	}
	return cm.client.IncrBy(ctx, key, increment).Result()
}
