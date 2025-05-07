package services

import (
	"context"
	"errors"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/repositories"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
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
	if err == nil && existingUser != nil {
		return "", errors.New("email already in use")
	}

	userId, err := s.userRepo.Create(ctx, input)
	if err != nil {
		return "", err
	}

	return userId.String(), nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, input *ports.UpdateUserInput) (string, error) {
	user, err := s.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return "", err
	}

	if input.Email != nil && *input.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, *input.Email)
		if err == nil && existingUser != nil && existingUser.ID != input.ID {
			return "", errors.New("email already in use by another user")
		}
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}

	if _, err := s.userRepo.Update(ctx, &ports.UpdateUserInput{
		ID:    input.ID,
		Name:  &user.Name,
		Email: &user.Email,
	}); err != nil {
		return "", err
	}

	return input.ID, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.List(ctx)
}
