package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
)

type ShortUrl struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	OriginalUrl string     `json:"original_url"`
	UserID      string     `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	AccessCount int        `json:"access_count"`
}

type CreateShortUrlInput struct {
	OriginalUrl string  `json:"original_url" validate:"required"`
	ExpiresAt   *string `json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateShortUrlInput struct {
	OriginalUrl *string `json:"original_url,omitempty"`
	ExpiresAt   *string `json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

func (i *CreateShortUrlInput) ToPgCreateShortUrl(userId uuid.UUID, slug string) *pgstore.CreateShortUrlParams {
	var expiresAt pgtype.Timestamptz
	if i.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *i.ExpiresAt)
		if err == nil {
			expiresAt.Time = t
			expiresAt.Valid = true
		}
	} else {
		expiresAt.Valid = false
	}

	return &pgstore.CreateShortUrlParams{
		Slug:        slug,
		UserID:      userId,
		OriginalUrl: i.OriginalUrl,
		ExpiresAt:   expiresAt,
	}
}

func (i *UpdateShortUrlInput) ApplyTo(shortUrl *ShortUrl) {
	if i.OriginalUrl != nil {
		shortUrl.OriginalUrl = *i.OriginalUrl
	}
	if i.ExpiresAt != nil {
		t, _ := time.Parse(time.RFC3339, *i.ExpiresAt)
		shortUrl.ExpiresAt = &t
	}
}

func (i *UpdateShortUrlInput) ToPgUpdateShortUrl(id uuid.UUID, slug string) *pgstore.UpdateShortUrlParams {
	var expiresAt pgtype.Timestamptz
	if i.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *i.ExpiresAt)
		if err == nil {
			expiresAt.Time = t
			expiresAt.Valid = true
		}
	} else {
		expiresAt.Valid = false
	}

	return &pgstore.UpdateShortUrlParams{
		ID:          id,
		Slug:        slug,
		OriginalUrl: *i.OriginalUrl,
		ExpiresAt:   expiresAt,
	}
}
