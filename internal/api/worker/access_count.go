package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/data/infra"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
	"github.com/redis/go-redis/v9"
)

const (
	maxConcurrentProcesses = 10
	workerInterval         = 1 * time.Hour
	workerTimeout          = 5 * time.Minute
	accessKeyPrefix        = "access:"
)

type AccessSyncWorker struct {
	db           *pgstore.Queries
	counter      *infra.AccessCounter
	logger       *slog.Logger
	interval     time.Duration
	concurrency  int
	shotdownChan chan struct{}
	wg           sync.WaitGroup
}

func NewAccessSyncWorker(db pgstore.Queries, rdb redis.Client) *AccessSyncWorker {
	return &AccessSyncWorker{
		db:           &db,
		counter:      infra.NewAccessCounter(&rdb),
		logger:       slog.Default().With("component", "access_sync_worker"),
		interval:     workerInterval,
		concurrency:  maxConcurrentProcesses,
		shotdownChan: make(chan struct{}),
	}
}

func (w *AccessSyncWorker) Start() error {
	if w.db == nil {
		return wraperrors.InternalErr("Cannot start access worker with nil database", nil)
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		w.logger.Info("starting access counter sync worker", "interval", w.interval, "concurrency", w.concurrency)

		ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
		w.processAccessCount(ctx)
		cancel()

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
				w.processAccessCount(ctx)
				cancel()

			case <-w.shotdownChan:
				w.logger.Info("acccess sync worker shtting down")
				return
			}
		}
	}()

	return nil
}

func (w *AccessSyncWorker) Stop() {
	close(w.shotdownChan)
	w.wg.Wait()
	w.logger.Info("access sync worker stopped")
}

func (w *AccessSyncWorker) processAccessCount(ctx context.Context) {
	start := time.Now()
	w.logger.Info("processing accss counts from Redis to database")

	processed := 0
	failed := 0
	var mu sync.Mutex

	keys, err := w.counter.GetAllAccessKeys(ctx)
	if err != nil {
		w.logger.Error("failed to get access keys", "error", err)
		return
	}

	if len(keys) == 0 {
		w.logger.Info("no access counters to process")
		return
	}

	sem := make(chan struct{}, w.concurrency)
	var wg sync.WaitGroup

	for _, key := range keys {
		if ctx.Err() != nil {
			w.logger.Warn("context canceled while processing access counts", "error", ctx.Err())
			break
		}

		slug := key[len(accessKeyPrefix):]

		wg.Add(1)
		sem <- struct{}{}

		go func(slug, key string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			count, err := w.counter.GetAndDeleteCounter(ctx, key)
			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				return
			}

			if count <= 0 {
				return
			}

			err = w.db.IncrementAccessCount(ctx, pgstore.IncrementAccessCountParams{
				Slug: slug,
				AccessCount: pgtype.Int4{
					Int32: int32(count),
					Valid: true,
				},
			})

			if err != nil {
				w.logger.Error("failed to update access count in database", "slug", slug, "count", count, "error", err)

				restoreErr := w.counter.RestoreCounter(ctx, key, int64(count))
				if restoreErr != nil {
					w.logger.Error("critical: failed to restore counter after DB error", "slug", slug, "count", count, "error", restoreErr)
				}

				mu.Lock()
				failed++
				mu.Unlock()
				return
			}

			mu.Lock()
			processed++
			mu.Unlock()
		}(slug, key)
	}

	wg.Wait()

	duration := time.Since(start)
	w.logger.Info("access count processing completed", "processed", processed, "failed", failed, "duration", duration, "keys_count", len(keys))
}

func StartHourlyAccessSyncWorker(db *pgstore.Queries, rdb *redis.Client) {
	if db == nil {
		slog.Error("cannot start access worker with nil database")
		return
	}

	if rdb == nil {
		slog.Error("cannot start access worker with nil redis client")
		return
	}

	worker := NewAccessSyncWorker(*db, *rdb)
	err := worker.Start()
	if err != nil {
		slog.Error("failed to start access sync worker", "error", err)
	}
}
