package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/repositories"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) ports.UserUseCase {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *ports.CreateUserInput) (string, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, wraperrors.NotFoundErr("")) {
		return "", wraperrors.InternalErr("could not check email", err)
	}
	if existingUser != nil {
		return "", wraperrors.AlreadyExistsErr("email already in use")
	}

	return s.userRepo.Create(ctx, &pgstore.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
	})
}

func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	return s.userRepo.GetByID(ctx, userId)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, id string, input *ports.UpdateUserInput) (string, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return "", wraperrors.NotFoundErr("user not found")
	}

	if input.Email != nil && *input.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, *input.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return "", wraperrors.AlreadyExistsErr("email already in use by another user")
		}
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}

	return s.userRepo.Update(ctx, &pgstore.UpdateUserParams{
		ID:    userId,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userId, err := uuid.Parse(id)
	if err != nil {
		return wraperrors.InternalErr("something went wrong", err)
	}

	return s.userRepo.Delete(ctx, userId)
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.List(ctx)
}
