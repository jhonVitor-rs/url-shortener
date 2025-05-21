package ports

import (
	"context"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type ShortUrlUseCase interface {
	ListShortUrl(ctx context.Context, userId string) ([]*models.ShortUrl, error)
	GetShortUrl(ctx context.Context, id string) (*models.ShortUrl, error)
	GetShortUrlBySlug(ctx context.Context, slug string) (*models.ShortUrl, error)
	CreateShortUrl(ctx context.Context, userId string, input *models.CreateShortUrlInput) (*models.ShortUrl, error)
	UpdateShortUrl(ctx context.Context, id string, input *models.UpdateShortUrlInput) (*models.ShortUrl, error)
	DeleteShortUrl(ctx context.Context, id string) error
}
