package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisStorage(prefix string, c *redis.Client) StorageInterface {
	return &RedisStorage{
		prefix: prefix,
		client: c,
	}
}

func (it *RedisStorage) Get(ctx context.Context, key Key) (Value, error) {
	fullKey := it.getKey(key)
	var ret Value
	val, err := it.client.Get(ctx, fullKey).Result()
	if err != nil && err != redis.Nil {
		return ret, err
	}
	if err == redis.Nil {
		return ret, nil
	}
	ret = Value(val)
	return ret, nil
}

func (it *RedisStorage) Set(ctx context.Context, key Key, val Value) error {
	fullKey := it.getKey(key)
	return it.client.Set(ctx, fullKey, string(val), time.Hour).Err()
}

func (it *RedisStorage) Delete(ctx context.Context, key Key) error {
	fullKey := it.getKey(key)
	return it.client.Del(ctx, fullKey).Err()
}
func (it *RedisStorage) getKey(key Key) string {
	return it.prefix + ":" + string(key)
}

type RedisStorage struct {
	prefix string
	client *redis.Client
}
