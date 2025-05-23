package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jhonVitor-rs/url-shortener/docs"
	api "github.com/jhonVitor-rs/url-shortener/internal/api/server"
	"github.com/jhonVitor-rs/url-shortener/internal/api/worker"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/rdstore"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

//	@title			URL Shortener API
//	@version		1.0
//	@description	API for creating and managing short URLs and users
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Jhon Vitor
//	@contact.email	seuemail@exemplo.com

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Failed to load environment variables")
		panic(err)
	}

	ctx := context.Background()

	pool := setupDatabseConnection(ctx)
	defer pool.Close()

	rdb := setupRedisConnection(ctx)
	defer rdb.Close()

	handler := api.NewApiHandler(pgstore.New(pool), rdb)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: handler,
	}

	go func() {
		slog.Info("Starting server on 0.0.0.0:8080")
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Server error", "error", err)
				panic(err)
			}
			slog.Info("Server closed")
		}
	}()

	worker.StartHourlyAccessSyncWorker(pgstore.New(pool), rdb)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}
}

func setupDatabseConnection(ctx context.Context) *pgxpool.Pool {
	connectionSring := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	config, err := pgxpool.ParseConfig(connectionSring)
	if err != nil {
		slog.Warn("Failed to create pool")
		panic(err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("Failed to create database pool", "error", err)
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Failed to ping database", "error", err)
		panic(err)
	}

	slog.Info("Database connection established successfully")
	return pool
}

func setupRedisConnection(ctx context.Context) *redis.Client {
	rdb := rdstore.NewRedisClient()
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := rdb.HealthCheck(checkCtx); err != nil {
		slog.Warn("Redis connection check failed - cache will be unavailable", "error", err)
		panic(err)
	} else {
		slog.Info("Redis connection with successfully")
	}

	return rdb.Client
}
