package ports

import (
	"context"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type UserUseCase interface {
	ListUsers(ctx context.Context) ([]*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error)
	UpdateUser(ctx context.Context, id string, input *models.UpdateUserInput) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
