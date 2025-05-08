package services

import (
	"context"
	"errors"

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

func (s *shortUrlService) CreateShortUrl(ctx context.Context, input *ports.CreateShortUrlInput) (string, error) {
	const maxAttempts = 5

	for range maxAttempts {
		slug, err := utils.CreateHashSlug()
		if err != nil {
			return "", err
		}

		_, err = s.shortUrlRepo.GetBySlug(ctx, slug)
		if err != nil {
			if errors.Is(err, wraperrors.NotFoundErr("")) {
				shortUrlId, err := s.shortUrlRepo.Create(ctx, &pgstore.CreateShortUrlParams{
					UserID:      input.UserID,
					Slug:        slug,
					OriginalUrl: input.OriginalUrl,
					ExpiresAt:   pgtype.Timestamptz{Time: *input.ExpiresAt, Valid: input.ExpiresAt != nil},
				})
				if err != nil && wraperrors.IsUniqueViolation(err) {
					continue
				}
				return shortUrlId, err
			}
			return "", err
		}
	}
	return "", wraperrors.InternalErr("failed to generated unique slug after several attemps", nil)
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

	return s.shortUrlRepo.Update(ctx, &pgstore.UpdateShortUrlParams{
		ID:          shortUrlId,
		Slug:        shortUrl.Slug,
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
