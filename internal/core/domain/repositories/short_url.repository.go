package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type ShortUrlUseRepository interface {
	Create(ctx context.Context, input *pgstore.CreateShortUrlParams) (string, error)
	GetShortUrl(ctx context.Context, id uuid.UUID) (*models.ShortUrl, error)
	GetBySlug(ctx context.Context, slug string) (*models.ShortUrl, error)
	Update(ctx context.Context, input *pgstore.UpdateShortUrlParams) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userId uuid.UUID) ([]*models.ShortUrl, error)
}
