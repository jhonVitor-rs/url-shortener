package models

import (
	"time"

	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

type GetUserByEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (i *CreateUserInput) ToPgCreateUser() *pgstore.CreateUserParams {
	return &pgstore.CreateUserParams{
		Name:  i.Name,
		Email: i.Email,
	}
}

func (i *UpdateUserInput) ApplyTo(user *User) {
	if i.Name != nil {
		user.Name = *i.Name
	}
	if i.Email != nil {
		user.Email = *i.Email
	}
}
