package rdstore

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func IncrementAccess(ctx context.Context, rdb *redis.Client, slug string) {
	if slug == "" {
		slog.Warn("attempted to increment access count with empty slug")
	}

	key := "access:" + slug
	err := rdb.Incr(ctx, key).Err()
	if err != nil {
		slog.Error("error to increment access count", "slug", slug, "error", err)
	}
}
