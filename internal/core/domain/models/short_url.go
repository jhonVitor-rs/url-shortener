package models

import "time"

type ShortUrl struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	OriginalUrl string     `json:"original_url"`
	UserID      string     `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
