package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	*redis.Client
}

func NewClient(addr string, db int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	return &Cache{client}
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

// Set usage:
// var origin Test
// err := cache.Set(ctx, "test", origin, 0)
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, key, marshaled, ttl).Err()
}

// Get usage:
// var dest Test
// err := cache.Get(context.Background(), "test", &dest)
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	cmdValStr, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return ErrNotFound
		}
		return err
	}

	err = json.Unmarshal([]byte(cmdValStr), &dest)
	if err != nil {
		return err
	}
	return nil
}
