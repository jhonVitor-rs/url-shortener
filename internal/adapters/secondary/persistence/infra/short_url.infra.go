package infra

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/repositories"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type shortUrlRepository struct {
	q *pgstore.Queries
}

func NewSHortUrlRepository(q *pgstore.Queries) repositories.ShortUrlUseRepository {
	return &shortUrlRepository{
		q: q,
	}
}

func (r *shortUrlRepository) Create(ctx context.Context, params *pgstore.CreateShortUrlParams) (string, error) {
	shortUrlId, err := r.q.CreateShortUrl(ctx, *params)
	if err != nil {
		if wraperrors.IsUniqueViolation(err) {
			return "", err
		}
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	return shortUrlId.String(), nil
}

func (r *shortUrlRepository) GetShortUrl(ctx context.Context, id uuid.UUID) (*models.ShortUrl, error) {
	dbShortUrl, err := r.q.GetShortUrlById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("Short URL not fund")
		}
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	shortUrl := &models.ShortUrl{
		ID:          dbShortUrl.ID.String(),
		Slug:        dbShortUrl.Slug,
		OriginalUrl: dbShortUrl.OriginalUrl,
		UserID:      dbShortUrl.UserID.String(),
		CreatedAt:   dbShortUrl.CreatedAt.Time,
		AccessCount: int(dbShortUrl.AccessCount.Int32),
	}

	if dbShortUrl.ExpiresAt.Valid {
		shortUrl.ExpiresAt = &dbShortUrl.ExpiresAt.Time
	} else {
		shortUrl.ExpiresAt = nil
	}

	return shortUrl, nil
}

func (r *shortUrlRepository) GetBySlug(ctx context.Context, slug string) (*models.ShortUrl, error) {
	dbShortUrl, err := r.q.GetShortUrlBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("err no rouws", "error", err)
			return nil, wraperrors.NotFoundErr(fmt.Sprintf("short URL with slug '%s' not found", slug))
		}
		slog.Error("e", "error", err)
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	if dbShortUrl.ExpiresAt.Valid && time.Now().After(dbShortUrl.ExpiresAt.Time) {
		return nil, wraperrors.UnauthorizedErr("short url expired")
	}

	shortUrl := &models.ShortUrl{
		ID:          dbShortUrl.ID.String(),
		Slug:        dbShortUrl.Slug,
		OriginalUrl: dbShortUrl.OriginalUrl,
		UserID:      dbShortUrl.UserID.String(),
		CreatedAt:   dbShortUrl.CreatedAt.Time,
	}

	if dbShortUrl.ExpiresAt.Valid {
		shortUrl.ExpiresAt = &dbShortUrl.ExpiresAt.Time
	} else {
		shortUrl.ExpiresAt = nil
	}

	return shortUrl, nil
}

func (r *shortUrlRepository) Update(ctx context.Context, params *pgstore.UpdateShortUrlParams) (string, error) {
	if _, err := r.q.UpdateShortUrl(ctx, *params); err != nil {
		return "", err
	}

	return params.ID.String(), nil
}

func (r *shortUrlRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteShortUrl(ctx, id); err != nil {
		return wraperrors.InternalErr("something went wrong", err)
	}
	return nil
}

func (r *shortUrlRepository) List(ctx context.Context, userId uuid.UUID) ([]*models.ShortUrl, error) {
	shortUrls, err := r.q.GetShortUrlsByUserId(ctx, userId)
	if err != nil {
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	var shortUrlsPointers []*models.ShortUrl
	for _, shortUrl := range shortUrls {
		shortUrlsPointers = append(shortUrlsPointers, &models.ShortUrl{
			ID:          shortUrl.ID.String(),
			Slug:        shortUrl.Slug,
			OriginalUrl: shortUrl.OriginalUrl,
			UserID:      shortUrl.UserID.String(),
			CreatedAt:   shortUrl.CreatedAt.Time,
			ExpiresAt: func() *time.Time {
				if shortUrl.ExpiresAt.Valid {
					return &shortUrl.ExpiresAt.Time
				}
				return nil
			}(),
			AccessCount: int(shortUrl.AccessCount.Int32),
		})
	}

	return shortUrlsPointers, nil
}
