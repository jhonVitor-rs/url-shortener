package infra

import (
	"context"
	"log/slog"
	"time"

	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
	"github.com/redis/go-redis/v9"
)

const (
	accessKeyPrefix = "access:"
	defaultReties   = 3
	retryDelay      = 100 * time.Millisecond
)

type AccessCounter struct {
	client *redis.Client
	logger *slog.Logger
}

func NewAccessCounter(client *redis.Client) *AccessCounter {
	return &AccessCounter{
		client: client,
		logger: slog.Default().With("component", "access_counter"),
	}
}

func (ac *AccessCounter) IncrementAccess(ctx context.Context, slug string) (int64, error) {
	if slug == "" {
		ac.logger.Warn("attempted to increment access count with empty slug")
		return 0, wraperrors.ValidationErr("Slug cannot be empty")
	}

	key := accessKeyPrefix + slug

	var count int64
	var err error

	for attempt := 0; attempt < defaultReties; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(retryDelay):
			}
		}

		count, err = ac.client.Incr(ctx, key).Result()
		if err == nil {
			return count, nil
		}

		ac.logger.Warn("failed to increment access count, retrying", "slug", slug, "error", err, "attempt", attempt+1)
	}

	ac.logger.Warn("failed to increment access count after retries", "slug", slug, "error", err)
	return 0, wraperrors.InternalErr("Failed to incremente access counter", err)
}

func (ac *AccessCounter) GetAllAccessKeys(ctx context.Context) ([]string, error) {
	var keys []string
	var cursor uint64

	for {
		var scanKeys []string
		var err error

		scanKeys, cursor, err = ac.client.Scan(ctx, cursor, accessKeyPrefix+"*", 0).Result()
		if err != nil {
			ac.logger.Error("failed to scan Redis for access keys", "error", err)
			return nil, wraperrors.InternalErr("Failed to retrieve access keys", err)
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (ac *AccessCounter) GetAndDeleteCounter(ctx context.Context, key string) (int, error) {
	count, err := ac.client.GetDel(ctx, key).Int()
	if err != nil {
		ac.logger.Error("failed to get and delete Redis counter", "key", key, "error", err)
		return 0, wraperrors.InternalErr("Failed to proccess counter", err)
	}

	return count, nil
}

func (ac *AccessCounter) RestoreCounter(ctx context.Context, key string, count int64) error {
	_, err := ac.client.IncrBy(ctx, key, count).Result()
	if err != nil {
		ac.logger.Error("failed to restore counter is Redis", "key", key, "error", err)
		return wraperrors.InternalErr("Failed to restore access counter", err)
	}

	return nil
}
