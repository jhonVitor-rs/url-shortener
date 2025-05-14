package infra

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/repositories"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type userRepository struct {
	q *pgstore.Queries
}

func NewUserRepository(q *pgstore.Queries) repositories.UserRepository {
	return &userRepository{
		q: q,
	}
}

func (r *userRepository) Create(ctx context.Context, params *pgstore.CreateUserParams) (string, error) {
	userId, err := r.q.CreateUser(ctx, *params)
	if err != nil {
		slog.Error("error to database", "error", err)
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	return userId.String(), nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	dbUser, err := r.q.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("User not fund")
		}
		slog.Error("error to database", "error", err)
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	return &models.User{
		ID:        dbUser.ID.String(),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("User not fund")
		}
		slog.Error("error to database", "error", err)
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	return &models.User{
		ID:        dbUser.ID.String(),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, params *pgstore.UpdateUserParams) (string, error) {
	if _, err := r.q.UpdateUser(ctx, *params); err != nil {
		slog.Error("error to database", "error", err)
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	return params.ID.String(), nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteUser(ctx, id); err != nil {
		return wraperrors.InternalErr("something went wrong", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*models.User, error) {
	users, err := r.q.GetUsers(ctx)
	if err != nil {
		slog.Error("error to database", "error", err)
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	if users == nil {
		users = []pgstore.User{}
	}

	var userPointers []*models.User
	for _, user := range users {
		userPointers = append(userPointers, &models.User{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
		})
	}

	return userPointers, nil
}
