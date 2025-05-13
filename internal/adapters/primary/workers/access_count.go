package workers

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/redis/go-redis/v9"
)

func StartHourlyAccessWorker(db *pgstore.Queries, rdb *redis.Client) {
	if db == nil {
		slog.Error("cannot start access worker with nil database")
		return
	}

	go func() {
		slog.Info("starting hourly access counter worker")
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		processAccessCount(context.Background(), db, rdb)

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			processAccessCount(ctx, db, rdb)
			cancel()
		}
	}()
}

func processAccessCount(ctx context.Context, db *pgstore.Queries, rdb *redis.Client) {
	start := time.Now()
	slog.Info("processing access counst from Redis to database")

	var procesed, failed int
	var mu sync.Mutex

	keys, err := getAllAccessKeys(ctx, rdb)
	if err != nil {
		slog.Error("failed to scan Redis for access keys", "error", err)
		return
	}

	if len(keys) == 0 {
		slog.Info("no access counters to process")
		return
	}

	const maxConcurrent = 10
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, key := range keys {
		if ctx.Err() != nil {
			slog.Warn("context canceled while processing access counts", "error", ctx.Err())
			break
		}

		slug := key[len("access:"):]
		wg.Add(1)
		sem <- struct{}{}

		go func(slug, key string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			count, err := rdb.GetDel(ctx, key).Int()
			if err != nil {
				slog.Error("failed to get and delete Redis key", "slug", slug, "error", err)

				mu.Lock()
				failed++
				mu.Unlock()
				return
			}

			if count > 0 {
				err = db.IncrementAccessCount(ctx, pgstore.IncrementAccessCountParams{
					Slug:        slug,
					AccessCount: pgtype.Int4{Int32: int32(count), Valid: true},
				})
				if err != nil {
					slog.Error("failed to update access count in database", "slug", slug, "error", err)

					_, restoreErr := rdb.IncrBy(ctx, key, int64(count)).Result()
					if restoreErr != nil {
						slog.Error("failed to restore counter in Redis after database failure", "slug", slug, "error", restoreErr)
					}

					mu.Lock()
					failed++
					mu.Unlock()
					return
				}
			}

			mu.Lock()
			procesed++
			mu.Unlock()
		}(slug, key)

		wg.Wait()

		duration := time.Since(start)
		slog.Info("access count processing completed", "processed", procesed, "failed", failed, "duration", duration)
	}
}

func getAllAccessKeys(ctx context.Context, rdb *redis.Client) ([]string, error) {
	var keys []string
	var cursor uint64

	for {
		var scanKeys []string
		var err error

		scanKeys, cursor, err = rdb.Scan(ctx, cursor, "access:*", 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
