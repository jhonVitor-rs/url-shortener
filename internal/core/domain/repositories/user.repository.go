package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *pgstore.CreateUserParams) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *pgstore.UpdateUserParams) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.User, error)
}
