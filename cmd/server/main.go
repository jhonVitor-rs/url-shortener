package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jhonVitor-rs/url-shortener/internal/api"
	"github.com/jhonVitor-rs/url-shortener/internal/store/pgstore"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Failed to load environment variables")
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	))

	if err != nil {
		slog.Warn("Failed to create pool")
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Warn("Failed to ping database")
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}
}
