package test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/rdstore"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func setupDatabseConnectionTests(ctx context.Context) *pgxpool.Pool {
	connectionSring := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s_test sslmode=disable",
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

func setupRedisConnectionTests(ctx context.Context) *redis.Client {
	rdb := rdstore.NewRedisClient()

	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := rdb.HealthCheck(checkCtx); err != nil {
		slog.Warn("Redis connection check failed - cache will be unavailable", "error", err)
	} else {
		slog.Info("Redis connection with successfully")
	}

	return rdb.Client
}

func Config(ctx context.Context) (*pgxpool.Pool, *redis.Client) {
	if err := godotenv.Load("../../../.env"); err != nil {
		slog.Warn("Failed to load environment variables")
		panic(err)
	}

	pool := setupDatabseConnectionTests(ctx)
	rdb := setupRedisConnectionTests(ctx)

	_, err := pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		CREATE TABLE IF NOT EXISTS users (
			"id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			"name" VARCHAR(50) NOT NULL,
			"email" VARCHAR(100) UNIQUE NOT NULL,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS short_urls (
			"id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			"user_id" uuid NOT NULL,
			"slug" TEXT UNIQUE NOT NULL,
			"original_url" TEXT NOT NULL,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"expires_at" TIMESTAMP WITH TIME ZONE,
			"access_count" INTEGER DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id) ON
			DELETE
				CASCADE
		);
	`)
	if err != nil {
		panic(err)
	}

	return pool, rdb
}
