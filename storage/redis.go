package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Prefix    string
	Client    *redis.Client
	CacheTime time.Duration
}

type Option func(*Config)

func SetCacheTime(duration time.Duration) Option {
	return func(c *Config) {
		c.CacheTime = duration
	}
}

func SetPrefix(prefix string) Option {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

func SetClient(client *redis.Client) Option {
	return func(c *Config) {
		c.Client = client
	}
}

func NewRedisStorage(prefix string, c *redis.Client, opts ...Option) StorageInterface {
	config := &Config{
		Prefix:    prefix,
		Client:    c,
		CacheTime: time.Hour,
	}
	for _, f := range opts {
		f(config)
	}
	return &RedisStorage{
		prefix:    config.Prefix,
		client:    config.Client,
		cacheTime: config.CacheTime,
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
	return it.client.Set(ctx, fullKey, string(val), it.cacheTime).Err()
}

func (it *RedisStorage) Delete(ctx context.Context, key Key) error {
	fullKey := it.getKey(key)
	return it.client.Del(ctx, fullKey).Err()
}
func (it *RedisStorage) getKey(key Key) string {
	return it.prefix + ":" + string(key)
}

type RedisStorage struct {
	prefix    string
	client    *redis.Client
	cacheTime time.Duration
}
