package rdstore

import (
	"context"
	"log/slog"
	"time"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
	"github.com/redis/go-redis/v9"
)

const (
	listKey    = "url:recent"
	maxLength  = 20
	urlPrefix  = "url:"
	defaultTTL = 24 * time.Hour
)

func LogRecentAccess(shortUrl *models.ShortUrl) error {
	go func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		urlKey := urlPrefix + shortUrl.Slug

		var ttl time.Duration
		if shortUrl.ExpiresAt != nil {
			ttl = time.Until(*shortUrl.ExpiresAt)
			if ttl <= 0 {
				slog.Debug("url already expired, not caching", "slug", shortUrl.Slug)
				return
			}
		} else {
			ttl = defaultTTL
		}

		if err := rdb.Set(timeoutCtx, urlKey, shortUrl.OriginalUrl, ttl).Err(); err != nil {
			slog.Error("faled to save url in caceh", "error", err)
			return
		}

		updateListUrls(timeoutCtx, shortUrl.Slug)
	}()

	return nil
}

func GetUrl(ctx context.Context, slug string) (string, error) {
	urlKey := urlPrefix + slug

	url, err := rdb.Get(ctx, urlKey).Result()
	if err != nil {
		// Independente do erro retornado vou ter que pesquisar no banco de dados
		if err == redis.Nil {
			return "", wraperrors.NotFoundErr("url not found in cache")
		}

		slog.Error("failed to retrieve url from cache",
			"slug", slug,
			"error", err)
		return "", wraperrors.NotFoundErr("url not found in cache")
	}

	go func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		updateListUrls(timeoutCtx, slug)
	}()

	return url, nil
}

func HealthCheck(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := rdb.Ping(timeoutCtx).Err(); err != nil {
		slog.Error("redis connection check failed", "error", err)
		return wraperrors.InternalErr("cache connection error", err)
	}

	return nil
}

func updateListUrls(ctx context.Context, slug string) {
	pipe := rdb.TxPipeline()

	pipe.LRem(ctx, listKey, 0, slug)

	pipe.LPush(ctx, listKey, slug)

	pipe.LLen(ctx, listKey)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.Error("failed to update recent url list", "error", err)
		return
	}

	listLen := cmds[2].(*redis.IntCmd).Val()

	if listLen > maxLength {
		cleanPipe := rdb.Pipeline()

		cleanPipe.LTrim(ctx, listKey, 0, maxLength-1)

		if _, err := cleanPipe.Exec(ctx); err != nil {
			slog.Error("failed to trim recent url list", "error", err)
		}
	}
}
