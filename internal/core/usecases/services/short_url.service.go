package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/repositories"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type shortUrlService struct {
	shortUrlRepo repositories.ShortUrlUseRepository
}

func NewShortUrlService(shortUrlRepo repositories.ShortUrlUseRepository) ports.ShortUrlUseCase {
	return &shortUrlService{
		shortUrlRepo: shortUrlRepo,
	}
}

func (s *shortUrlService) CreateShortUrl(ctx context.Context, rawUserId string, input *ports.CreateShortUrlInput) (string, error) {
	slug, err := s.genSlug(ctx, 1)
	if err != nil {
		return "", err
	}

	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		slog.Error("error to convert user id", "error", err)
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	return s.shortUrlRepo.Create(ctx, &pgstore.CreateShortUrlParams{
		UserID:      userId,
		Slug:        slug,
		OriginalUrl: input.OriginalUrl,
		ExpiresAt:   pgtype.Timestamptz{Time: *input.ExpiresAt, Valid: input.ExpiresAt != nil},
	})
}

func (s *shortUrlService) GetShortUrl(ctx context.Context, id string) (*models.ShortUrl, error) {
	shortUrlId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	return s.shortUrlRepo.GetShortUrl(ctx, shortUrlId)
}

func (s *shortUrlService) GetShortUrlBySlug(ctx context.Context, slug string) (*models.ShortUrl, error) {
	return s.shortUrlRepo.GetBySlug(ctx, slug)
}

func (s *shortUrlService) UpdateShortUrl(ctx context.Context, id string, input *ports.UpdateShortUrlInput) (string, error) {
	shortUrlId, err := uuid.Parse(id)
	if err != nil {
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	shortUrl, err := s.shortUrlRepo.GetShortUrl(ctx, shortUrlId)
	if err != nil {
		return "", wraperrors.NotFoundErr("short err not found")
	}

	if input.OriginalUrl != nil {
		shortUrl.OriginalUrl = *input.OriginalUrl
	}
	if input.ExpiresAt != nil {
		shortUrl.ExpiresAt = input.ExpiresAt
	}

	newSlug, err := s.genSlug(ctx, 1)
	if err != nil {
		return "", err
	}

	return s.shortUrlRepo.Update(ctx, &pgstore.UpdateShortUrlParams{
		ID:          shortUrlId,
		Slug:        newSlug,
		OriginalUrl: shortUrl.OriginalUrl,
		ExpiresAt:   pgtype.Timestamptz{Time: *shortUrl.ExpiresAt, Valid: shortUrl.ExpiresAt != nil},
	})
}

func (s *shortUrlService) DeleteShortUrl(ctx context.Context, id string) error {
	shortUrlId, err := uuid.Parse(id)
	if err != nil {
		return wraperrors.InternalErr("something went wrong", err)
	}

	return s.shortUrlRepo.Delete(ctx, shortUrlId)
}

func (s *shortUrlService) ListShortUrl(ctx context.Context, rawUserId string) ([]*models.ShortUrl, error) {
	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	return s.shortUrlRepo.List(ctx, userId)
}

func (s *shortUrlService) genSlug(ctx context.Context, count int) (string, error) {
	if count > 5 {
		return "", wraperrors.InternalErr("failed to generate unique slug after several attempts", nil)
	}

	slug, err := utils.CreateHashSlug()
	if err != nil {
		return "", wraperrors.InternalErr("failed to create hash slug", err)
	}

	shortUrl, err := s.shortUrlRepo.GetBySlug(ctx, slug)

	var appErr *wraperrors.AppError
	if errors.As(err, &appErr) && appErr.Is(wraperrors.ErrNotFound) {
		return slug, nil
	}

	if err == nil && shortUrl != nil {
		return s.genSlug(ctx, count+1)
	}

	return "", wraperrors.InternalErr("error checking slug existence", err)
}
