package rdstore

import (
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

type rdb struct {
	Client *redis.Client
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

// E depois use:
func NewRedisClient() *rdb {
	return &rdb{
		Client: initRedisClient(),
	}
}
