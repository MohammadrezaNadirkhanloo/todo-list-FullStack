package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache: key not found")

type Client struct {
	*redis.Client
}

func Connect(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("cache: Redis at %s did not respond: %w", cfg.Addr(), err)
	}

	// log.Info(ctx, "connected to Redis",
	// 	logging.Cat(logging.CategoryRedis),
	// 	logging.Sub(logging.SubConnection),
	// 	logging.F("addr", cfg.Addr()),
	// 	logging.F("db", cfg.DB),
	// )
	fmt.Println("connected to Redis")
	return &Client{Client: rdb}, nil
}

func (c *Client) HealthCheck(ctx context.Context) error {
	return c.Ping(ctx).Err()
}

func Set[T any](ctx context.Context, c *Client, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: failed to marshal value to JSON for key %q: %w", key, err)
	}
	if err := c.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache: failed to store key %q: %w", key, err)
	}
	return nil
}

func Get[T any](ctx context.Context, c *Client, key string) (T, error) {
	var out T

	data, err := c.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return out, ErrCacheMiss
		}
		return out, fmt.Errorf("cache: failed to read key %q: %w", key, err)
	}

	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("cache: failed to unmarshal JSON for key %q: %w", key, err)
	}
	return out, nil
}

func Delete(ctx context.Context, c *Client, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := c.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("cache: failed to delete keys: %w", err)
	}
	return nil
}

func SetNX[T any](ctx context.Context, c *Client, key string, value T, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("cache: failed to marshal value to JSON for key %q: %w", key, err)
	}

	ok, err := c.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("cache: SetNX operation on key %q failed: %w", key, err)
	}
	return ok, nil
}
