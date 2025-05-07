package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	short "github.com/TomerHeber/go-short-url"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhonVitor-rs/url-shortener/internal/store/pgstore"
)

// Users Handlers ------------------------------------------------
func (h apiHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	body, ok := parseAndValidate[request](w, r)
	if !ok {
		return
	}

	userId, err := h.q.InsertUser(r.Context(), pgstore.InsertUserParams{
		Name: body.Name, Email: body.Email,
	})
	if err != nil {
		slog.Error("failed to insert user", slog.Any("err", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, response{ID: userId.String()})
}

func (h apiHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.q.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to get users", "error", err)
		return
	}

	if users == nil {
		users = []pgstore.User{}
	}
	sendJSON(w, users)
}

func (h apiHandler) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	rawUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.q.GetUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to get user", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, user)
}

func (h apiHandler) handleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"string" validate:"required, email"`
	}

	body, ok := parseAndValidate[request](w, r)
	if !ok {
		return
	}

	user, err := h.q.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to get user", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, user)
}

func (h apiHandler) handelUpdateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	rawUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	body, ok := parseAndValidate[request](w, r)
	if !ok {
		return
	}

	existing, err := h.q.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	name := existing.Name
	if body.Name != nil {
		name = *body.Name
	}

	email := existing.Email
	if body.Email != nil {
		email = *body.Email
	}

	userId, err = h.q.UpdateUser(r.Context(), pgstore.UpdateUserParams{ID: userId, Name: name, Email: email})
	if err != nil {
		slog.Error("failed to update user", slog.Any("err", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, response{ID: userId.String()})
}

func (h apiHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	rawUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.q.DeleteUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "user not found!", http.StatusBadRequest)
			return
		}
		slog.Error("Failed to delete a user", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, "User deleted with success!")
}

// Short URLs handlers-----------------------
func (h apiHandler) handleCreateShortUrl(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserId      uuid.UUID  `json:"user_id" validate:"required"`
		OriginalUrl string     `json:"original_url" validate:"required"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	body, ok := parseAndValidate[request](w, r)
	if !ok {
		return
	}

	if ok = validateUrl(w, body.OriginalUrl); !ok {
		return
	}

	var expires pgtype.Timestamp
	if body.ExpiresAt != nil && !body.ExpiresAt.IsZero() {
		expires = pgtype.Timestamp{Time: *body.ExpiresAt, Valid: true}
	} else {
		expires = pgtype.Timestamp{Valid: false}
	}

	h.saveNewShortUrl(w, r, body.UserId, body.OriginalUrl, expires)
}

func (h apiHandler) saveNewShortUrl(w http.ResponseWriter, r *http.Request, userId uuid.UUID, originalUrl string, expiresAt pgtype.Timestamp) {
	for range 5 {
		slug, ok := createHashSlug(w)
		if !ok {
			return
		}

		h.mu.Lock()
		shortUrlId, err := h.q.CreateShortUrl(r.Context(), pgstore.CreateShortUrlParams{
			UserID:      userId,
			Slug:        slug,
			OriginalUrl: originalUrl,
			ExpiresAt:   expiresAt,
		})
		h.mu.Unlock()

		if err == nil {
			sendJSON(w, response{ID: shortUrlId.String()})
		}

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "short_urls_slug_key" {
			continue
		}

		slog.Error("Failed to create short url", "error", err)
		http.Error(w, "someting went wrong", http.StatusInternalServerError)
		return
	}
	http.Error(w, "could not generate unique slug", http.StatusInternalServerError)
}

func (h apiHandler) handleGetShortUrlsByUser(w http.ResponseWriter, r *http.Request) {
	rawUserId := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(rawUserId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	shortUrls, err := h.q.GetShortUrlsByUserId(r.Context(), userId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if shortUrls == nil {
		shortUrls = []pgstore.ShortUrl{}
	}

	sendJSON(w, shortUrls)
}

func (h apiHandler) handleGetShorUrl(w http.ResponseWriter, r *http.Request) {
	rawShortUrlId := chi.URLParam(r, "short_url_id")
	shortUrlId, err := uuid.Parse(rawShortUrlId)
	if err != nil {
		http.Error(w, "invalid short url id", http.StatusBadRequest)
		return
	}

	shortUrl, err := h.q.GetShortUrlById(r.Context(), shortUrlId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "short url not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get a short url", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, shortUrl)
}

func (h apiHandler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortUrlSlug := chi.URLParam(r, "slug")

	shortUrl, err := h.q.GetShortUrlBySlug(r.Context(), shortUrlSlug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "short url not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get a short url", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	if shortUrl.ExpiresAt.Valid && time.Now().After(shortUrl.ExpiresAt.Time) {
		http.Error(w, "short url expired", http.StatusGone)
		return
	}

	http.Redirect(w, r, shortUrl.OriginalUrl, http.StatusFound)
}

func (h apiHandler) handleUpdateShortUrl(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserId      *uuid.UUID `json:"user_id"`
		OriginalUrl *string    `json:"original_url"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	rawShortUrlId := chi.URLParam(r, "short_url_id")
	shortUrlId, err := uuid.Parse(rawShortUrlId)
	if err != nil {
		http.Error(w, "invalid short url id", http.StatusBadRequest)
		return
	}

	body, ok := parseAndValidate[request](w, r)
	if !ok {
		return
	}

	existing, err := h.q.GetShortUrlById(r.Context(), shortUrlId)
	if err != nil {
		http.Error(w, "short url not found", http.StatusNotFound)
		return
	}

	// Atualizações condicionais
	originalUrl := existing.OriginalUrl
	if body.OriginalUrl != nil {
		if !validateUrl(w, *body.OriginalUrl) {
			return
		}
		originalUrl = *body.OriginalUrl
	}

	expiresAt := existing.ExpiresAt
	if body.ExpiresAt != nil {
		expiresAt = pgtype.Timestamp{Time: *body.ExpiresAt, Valid: true}
	}

	// Novo slug se URL mudar
	slug := existing.Slug
	if body.OriginalUrl != nil {
		s, _ := short.NewShortener()
		surl, err := s.CreateShortenedUrl(r.Context(), *body.OriginalUrl)
		if err != nil {
			http.Error(w, "failed to generate slug", http.StatusInternalServerError)
			return
		}
		slug = surl.GetUrl()
	}

	// Atualiza no banco
	_, err = h.q.UpdateShortUrl(r.Context(), pgstore.UpdateShortUrlParams{
		ID:          shortUrlId,
		OriginalUrl: originalUrl,
		Slug:        slug,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		slog.Error("failed to update short url", slog.Any("err", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, response{ID: shortUrlId.String()})
}

func (h apiHandler) handleDeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	rawShortUrlId := chi.URLParam(r, "short_url_id")
	shortUrlId, err := uuid.Parse(rawShortUrlId)
	if err != nil {
		http.Error(w, "invalid short url id", http.StatusBadRequest)
		return
	}

	err = h.q.DeleteShortUrl(r.Context(), shortUrlId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "short url not found!", http.StatusNotFound)
			return
		}
		slog.Error("Failed to delete a short url", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, "Short URL has deleted with success!")
}
