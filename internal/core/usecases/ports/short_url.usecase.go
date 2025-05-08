package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type CreateShortUrlInput struct {
	UserID      uuid.UUID  `json:"user_id" validate:"required"`
	OriginalUrl string     `json:"original_url" validate:"required"`
	ExpiresAt   *time.Time `json:"epires_at"`
}

type UpdateShortUrlInput struct {
	OriginalUrl *string    `json:"original_url"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type ShortUrlUseCase interface {
	CreateShortUrl(ctx context.Context, input *CreateShortUrlInput) (string, error)
	GetShortUrl(ctx context.Context, id string) (*models.ShortUrl, error)
	GetShortUrlBySlug(ctx context.Context, slug string) (*models.ShortUrl, error)
	UpdateShortUrl(ctx context.Context, id string, input *UpdateShortUrlInput) (string, error)
	DeleteShortUrl(ctx context.Context, id string) error
	ListShortUrl(ctx context.Context, userId string) ([]*models.ShortUrl, error)
}
