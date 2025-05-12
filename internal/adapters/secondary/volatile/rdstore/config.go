package rdstore

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func initRedisClient() *redis.Client {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Failed to load environment variables")
		panic(err)
	}

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
var rdb = initRedisClient()
