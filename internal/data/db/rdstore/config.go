package rdstore

import (
	"context"
	"log/slog"
	"os"
	"time"

	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
	"github.com/redis/go-redis/v9"
)

type rdb struct {
	Client *redis.Client
	logger *slog.Logger
}

func initRedisClient() *redis.Client {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")

	// Adicione logs para debugar
	slog.Info("Inicializando cliente Redis",
		"host", redisHost,
		"senha_definida", redisPassword != "",
		"senha_vazia", redisPassword == "")

	return redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
		PoolSize: 10,
	})
}

func NewRedisClient() *rdb {
	return &rdb{
		Client: initRedisClient(),
		logger: slog.Default().With("component", "url_cache"),
	}
}

// HealthCheck verifica a saúde da conexão com o Redis
func (r *rdb) HealthCheck(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := r.Client.Ping(timeoutCtx).Err(); err != nil {
		r.logger.Error("redis connection check failed", "error", err)
		return wraperrors.InternalErr("cache connection error", err)
	}

	return nil
}
