package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vijayshahwal/jobScheduling/interfaces"
)

type RedisCacheRepository struct {
	client *redis.Client
}

var (
	redisInstance *RedisCacheRepository
	redisOnce     sync.Once
)

func NewRedisCacheRepository(addr, password string, db int) interfaces.CacheRepository {
	redisOnce.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		ctx := context.Background()
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}

		redisInstance = &RedisCacheRepository{client: rdb}
	})
	return redisInstance
}

func (rcr *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return rcr.client.Set(ctx, key, string(jsonData), expiration).Err()
}

func (rcr *RedisCacheRepository) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := rcr.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s does not exist", key)
	}
	return val, err
}

func (rcr *RedisCacheRepository) Scan(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	var cursor uint64

	for {
		scanKeys, newCursor, err := rcr.client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scanKeys...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (rcr *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	return rcr.client.Del(ctx, key).Err()
}
