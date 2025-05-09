package ports

import (
	"context"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email" validate:"email"`
}

type GetUserByEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type UserUseCase interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (string, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, input *UpdateUserInput) (string, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]*models.User, error)
}
