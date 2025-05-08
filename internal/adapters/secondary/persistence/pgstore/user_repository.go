package pgstore

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/repositories"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

type userRepository struct {
	q *Queries
}

func NewUserRepository(q *Queries) repositories.UserRepository {
	return &userRepository{
		q: q,
	}
}

func (r *userRepository) Create(ctx context.Context, user *ports.CreateUserInput) (string, error) {
	params := CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
	}

	userId, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	return userId.String(), nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	dbUser, err := r.q.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, wraperrors.NotFoundErr("User not fund")
		}
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
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	return &models.User{
		ID:        dbUser.ID.String(),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *ports.UpdateUserInput) (string, error) {
	userId, err := uuid.Parse(user.ID)
	if err != nil {
		return "", wraperrors.InternalErr("something went wrong", err)
	}

	params := UpdateUserParams{
		ID:    userId,
		Name:  *user.Name,
		Email: *user.Email,
	}
	if _, err := r.q.UpdateUser(ctx, params); err != nil {
		return "", err
	}

	return userId.String(), nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	userId, err := uuid.Parse(id)
	if err != nil {
		return wraperrors.ValidationErr("user not found")
	}

	return r.q.DeleteUser(ctx, userId)
}

func (r *userRepository) List(ctx context.Context) ([]*models.User, error) {
	users, err := r.q.GetUsers(ctx)
	if err != nil {
		return nil, wraperrors.InternalErr("something went wrong", err)
	}

	if users == nil {
		users = []User{}
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
