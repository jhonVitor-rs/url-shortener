package ports

import (
	"context"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type CreateShortUrlInput struct {
	OriginalUrl string `json:"original_url" validate:"required"`
}

type UpdateShortUrlInput struct {
	OriginalUrl *string `json:"original_url"`
}

type ShortUrlUseCase interface {
	CreateShortUrl(ctx context.Context, userId string, input *CreateShortUrlInput) (string, error)
	GetShortUrl(ctx context.Context, id string) (*models.ShortUrl, error)
	GetShortUrlBySlug(ctx context.Context, slug string) (*models.ShortUrl, error)
	UpdateShortUrl(ctx context.Context, id string, input *UpdateShortUrlInput) (string, error)
	DeleteShortUrl(ctx context.Context, id string) error
	ListShortUrl(ctx context.Context, userId string) ([]*models.ShortUrl, error)
}
