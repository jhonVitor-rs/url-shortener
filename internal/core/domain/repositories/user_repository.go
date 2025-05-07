package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
)

type UserRepository interface {
	Create(ctx context.Context, user *ports.CreateUserInput) (uuid.UUID, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *ports.UpdateUserInput) (uuid.UUID, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*models.User, error)
}
