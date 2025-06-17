package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	client       *redis.Client
	ctx          = context.Background()
	ErrCacheMiss = redis.Nil
)

func InitRedis(addr string) error {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return client.Ping(ctx).Err()
}

func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Set(key string, value string, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}
