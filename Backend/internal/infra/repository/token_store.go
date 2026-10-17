package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/cache"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/redis/go-redis/v9"
)

type redisTokenStore struct {
	client *cache.Client
}

var _ repository.TokenStore = (*redisTokenStore)(nil)

func NewTokenStore(client *cache.Client) repository.TokenStore {
	return &redisTokenStore{client: client}
}

func tokenKey(hash string) string      { return "refresh:tok:" + hash }
func userIndexKey(userID int64) string { return fmt.Sprintf("refresh:usr:%d", userID) }

func (s *redisTokenStore) Save(ctx context.Context, userID int64, tokenHash string, ttl time.Duration) error {
	pipe := s.client.TxPipeline()

	pipe.Set(ctx, tokenKey(tokenHash), strconv.FormatInt(userID, 10), ttl)
	pipe.SAdd(ctx, userIndexKey(userID), tokenHash)
	pipe.Expire(ctx, userIndexKey(userID), ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return apperror.Wrap(err, apperror.CodeUnavailable,
			"saving the session failed, please try again.")
	}
	return nil
}

func (s *redisTokenStore) Consume(ctx context.Context, tokenHash string) (int64, bool, error) {
	val, err := s.client.GetDel(ctx, tokenKey(tokenHash)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, false, nil
		}
		return 0, false, apperror.Wrap(err, apperror.CodeUnavailable,
			"verifying the session failed, please try again.")
	}

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false, nil
	}

	if err := s.client.SRem(ctx, userIndexKey(userID), tokenHash).Err(); err != nil {
		_ = err
	}

	return userID, true, nil
}

func (s *redisTokenStore) RevokeAll(ctx context.Context, userID int64) error {
	indexKey := userIndexKey(userID)

	hashes, err := s.client.SMembers(ctx, indexKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return apperror.Wrap(err, apperror.CodeUnavailable, "revoking sessions failed.")
	}

	keys := make([]string, 0, len(hashes)+1)
	for _, h := range hashes {
		keys = append(keys, tokenKey(h))
	}
	keys = append(keys, indexKey)

	if err := s.client.Del(ctx, keys...).Err(); err != nil {
		return apperror.Wrap(err, apperror.CodeUnavailable, "revoking sessions failed.")
	}
	return nil
}
