package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type shortUrlService struct {
	db *pgstore.Queries
}

func NewShortUrlService(queries *pgstore.Queries) ports.ShortUrlUseCase {
	return &shortUrlService{
		db: queries,
	}
}

func (s *shortUrlService) ListShortUrl(ctx context.Context, rawUserId string) ([]*models.ShortUrl, error) {
	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		return nil, wraperrors.ValidationErr("Invalid user ID format")
	}

	dbShortUrls, err := s.db.GetShortUrlsByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("Short URLs not found for this user ID")
		}
		return nil, wraperrors.InternalErr("Failed to list short URLs", err)
	}

	shortUrls := make([]*models.ShortUrl, 0, len(dbShortUrls))
	for _, dbShortUrl := range dbShortUrls {
		shortUrls = append(shortUrls, &models.ShortUrl{
			ID:          dbShortUrl.ID.String(),
			Slug:        dbShortUrl.Slug,
			OriginalUrl: dbShortUrl.OriginalUrl,
			ExpiresAt:   &dbShortUrl.ExpiresAt.Time,
			CreatedAt:   dbShortUrl.CreatedAt.Time,
		})
	}

	return shortUrls, nil
}

func (s *shortUrlService) GetShortUrl(ctx context.Context, id string) (*models.ShortUrl, error) {
	shortUrlId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.ValidationErr("Invalid short URL ID format")
	}

	dbShortUrl, err := s.db.GetShortUrlById(ctx, shortUrlId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("Short URL not found")
		}
		return nil, wraperrors.InternalErr("Failed to get short URL", err)
	}

	return &models.ShortUrl{
		ID:          dbShortUrl.ID.String(),
		Slug:        dbShortUrl.Slug,
		OriginalUrl: dbShortUrl.OriginalUrl,
		ExpiresAt:   &dbShortUrl.ExpiresAt.Time,
		CreatedAt:   dbShortUrl.CreatedAt.Time,
	}, nil
}

func (s *shortUrlService) GetShortUrlBySlug(ctx context.Context, slug string) (*models.ShortUrl, error) {
	dbShortUrl, err := s.db.GetShortUrlBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("Short URL not found")
		}
		return nil, wraperrors.InternalErr("Failed to get short URL", err)
	}

	if dbShortUrl.ExpiresAt.Valid && dbShortUrl.ExpiresAt.Time.Before(time.Now()) {
		return nil, wraperrors.NotFoundErr("Short URL has expired")
	}

	return &models.ShortUrl{
		ID:          dbShortUrl.ID.String(),
		Slug:        dbShortUrl.Slug,
		OriginalUrl: dbShortUrl.OriginalUrl,
		ExpiresAt:   &dbShortUrl.ExpiresAt.Time,
		CreatedAt:   dbShortUrl.CreatedAt.Time,
	}, nil
}

func (s *shortUrlService) CreateShortUrl(ctx context.Context, rawUserId string, input *models.CreateShortUrlInput) (*models.ShortUrl, error) {
	slug, err := s.genSlug(ctx, 1)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		return nil, wraperrors.ValidationErr("Invalid user ID format")
	}

	pgShortUrl := input.ToPgCreateShortUrl(userId, slug)
	dbShortUrl, err := s.db.CreateShortUrl(ctx, *pgShortUrl)
	if err != nil {
		return nil, wraperrors.InternalErr("Failed to create short URL", err)
	}

	return &models.ShortUrl{
		ID:          dbShortUrl.ID.String(),
		Slug:        dbShortUrl.Slug,
		OriginalUrl: dbShortUrl.OriginalUrl,
		ExpiresAt:   &dbShortUrl.ExpiresAt.Time,
		CreatedAt:   dbShortUrl.CreatedAt.Time,
	}, nil
}

func (s *shortUrlService) UpdateShortUrl(ctx context.Context, id string, input *models.UpdateShortUrlInput) (*models.ShortUrl, error) {
	shortUrlId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.ValidationErr("Invalid short URL ID format")
	}

	shortUrl, err := s.GetShortUrl(ctx, id)
	if err != nil {
		return nil, err
	}

	slug, err := s.genSlug(ctx, 1)
	if err != nil {
		return nil, err
	}

	input.ApplyTo(shortUrl)
	pgShortUrl := input.ToPgUpdateShortUrl(shortUrlId, slug)

	dbShortUrl, err := s.db.UpdateShortUrl(ctx, *pgShortUrl)
	if err != nil {
		return nil, wraperrors.InternalErr("Failed to update short URL", err)
	}

	return &models.ShortUrl{
		ID:          dbShortUrl.ID.String(),
		Slug:        dbShortUrl.Slug,
		OriginalUrl: dbShortUrl.OriginalUrl,
		ExpiresAt:   &dbShortUrl.ExpiresAt.Time,
		CreatedAt:   dbShortUrl.CreatedAt.Time,
	}, nil
}

func (s *shortUrlService) DeleteShortUrl(ctx context.Context, id string) error {
	shortUrlId, err := uuid.Parse(id)
	if err != nil {
		return wraperrors.InternalErr("something went wrong", err)
	}

	err = s.db.DeleteShortUrl(ctx, shortUrlId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return wraperrors.NotFoundErr("Short URL not found")
		}
		return wraperrors.InternalErr("Failed to update short URL", err)
	}

	return nil
}

func (s *shortUrlService) genSlug(ctx context.Context, count int) (string, error) {
	if count > 5 {
		return "", wraperrors.InternalErr("Failed to generate unique slug after several attempts", nil)
	}

	slug, err := utils.CreateHashSlug()
	if err != nil {
		return "", wraperrors.InternalErr("Failed to create hash slug", err)
	}

	shortUrl, err := s.GetShortUrlBySlug(ctx, slug)

	if err != nil && wraperrors.IsNotFoundError(err) {
		return slug, nil
	}

	if err == nil && shortUrl != nil {
		return s.genSlug(ctx, count+1)
	}

	return "", wraperrors.InternalErr("error checking slug existence", err)
}
