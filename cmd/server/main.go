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
	if err := godotenv.Load(); err != err {
		slog.Warn("Failed to get environments varaibles")
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
		slog.Warn("Failed to pool.Ping!")
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe("0.0.0.0:8080", handler); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Warn("Failed to get environments variables")
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
