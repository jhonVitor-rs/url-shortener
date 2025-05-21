package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type userService struct {
	db *pgstore.Queries
}

func NewUserService(queries *pgstore.Queries) ports.UserUseCase {
	return &userService{
		db: queries,
	}
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.User, error) {
	dbUsers, err := s.db.GetUsers(ctx)
	if err != nil {
		return nil, wraperrors.InternalErr("Failed to list users", err)
	}

	users := make([]*models.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, &models.User{
			ID:        dbUser.ID.String(),
			Name:      dbUser.Name,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt.Time,
		})
	}

	return users, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.ValidationErr("Invalid user ID format")
	}

	dbUser, err := s.db.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("User not found")
		}
		return nil, wraperrors.InternalErr("Failed to get user", err)
	}

	return &models.User{
		ID:        dbUser.ID.String(),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	dbUser, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("User not found")
		}
		return nil, wraperrors.InternalErr("Failed to get user by email", err)
	}

	return &models.User{
		ID:        dbUser.ID.String(),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	existingUser, err := s.GetUserByEmail(ctx, input.Email)
	if err != nil && !wraperrors.IsNotFoundError(err) {
		return nil, wraperrors.InternalErr("Could not check email", err)
	}
	if existingUser != nil {
		return nil, wraperrors.AlreadyExistsErr("Email already in use")
	}

	pgUser := input.ToPgCreateUser()
	user, err := s.db.CreateUser(ctx, *pgUser)
	if err != nil {
		return nil, wraperrors.InternalErr("Failed to create user", err)
	}

	return &models.User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, input *models.UpdateUserInput) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.ValidationErr("Invalid user ID format")
	}

	user, err := s.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	input.ApplyTo(user)

	if input.Email != nil && *input.Email != user.Email {
		existingUser, err := s.GetUserByEmail(ctx, *input.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return nil, wraperrors.AlreadyExistsErr("Email already in use by another user")
		}
	}

	newUser, err := s.db.UpdateUser(ctx, pgstore.UpdateUserParams{
		ID:    userId,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, wraperrors.InternalErr("Failed to update user", err)
	}

	return &models.User{
		ID:        newUser.ID.String(),
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Time,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userId, err := uuid.Parse(id)
	if err != nil {
		return wraperrors.ValidationErr("Invalid user ID format")
	}

	err = s.db.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return wraperrors.NotFoundErr("User not found")
		}
		return wraperrors.InternalErr("Failed to delete user", err)
	}

	return nil
}
