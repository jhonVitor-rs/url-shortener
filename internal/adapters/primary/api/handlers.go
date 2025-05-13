package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/primary/middlewares"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/volatile/rdstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
)

// User handlrs -----------------------------------------
func (h apiHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, ok := utils.ParseAndValidate[ports.CreateUserInput](w, r)
	if !ok {
		return
	}

	userId, err := h.user.CreateUser(r.Context(), body)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, models.Response{ID: userId})
}

func (h apiHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.user.ListUsers(r.Context())
	if err != nil {
		utils.SendErrors(w, err)
	}

	utils.SendJSON(w, users)
}

func (h apiHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	user, err := h.user.GetUser(r.Context(), userId)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, user)
}

func (h apiHandler) handleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	body, ok := utils.ParseAndValidate[ports.GetUserByEmailInput](w, r)
	if !ok {
		return
	}

	user, err := h.user.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	token, err := middlewares.GenerateJWT(user.ID)
	if err != nil {
		slog.Error("Err to generate toke", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, models.Token{JWT: token})
}

func (h apiHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	body, ok := utils.ParseAndValidate[ports.UpdateUserInput](w, r)
	if !ok {
		return
	}

	userId, err := h.user.UpdateUser(r.Context(), userId, body)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, models.Response{ID: userId})
}

func (h apiHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	if err := h.user.DeleteUser(r.Context(), userId); err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, "User deleted with success")
}

// Short URL handlers -----------------------------------------------
func (h apiHandler) handleCreateShortUrl(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	body, ok := utils.ParseAndValidate[ports.CreateShortUrlInput](w, r)
	if !ok {
		return
	}

	if ok = utils.ValidateUrl(w, body.OriginalUrl); !ok {
		return
	}

	shortUrlId, err := h.shortUrl.CreateShortUrl(r.Context(), userId, body)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, models.Response{ID: shortUrlId})
}

func (h apiHandler) handleListShortUrlsByUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	shortUrls, err := h.shortUrl.ListShortUrl(r.Context(), userId)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, shortUrls)
}

func (h apiHandler) handleGetShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")

	shortUrl, err := h.shortUrl.GetShortUrl(r.Context(), shortUrlId)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, shortUrl)
}

func (h apiHandler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	url, err := rdstore.GetUrl(r.Context(), h.rdb, slug)
	if err == nil {
		rdstore.IncrementAccess(r.Context(), h.rdb, slug)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	shortUrl, err := h.shortUrl.GetShortUrlBySlug(r.Context(), slug)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	rdstore.LogRecentAccess(r.Context(), h.rdb, shortUrl)
	rdstore.IncrementAccess(r.Context(), h.rdb, slug)
	http.Redirect(w, r, shortUrl.OriginalUrl, http.StatusFound)
}

func (h apiHandler) handleUpdateShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")
	body, ok := utils.ParseAndValidate[ports.UpdateShortUrlInput](w, r)
	if !ok {
		return
	}

	shortUrlId, err := h.shortUrl.UpdateShortUrl(r.Context(), shortUrlId, body)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, models.Response{ID: shortUrlId})
}

func (h apiHandler) handleDeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")
	if err := h.shortUrl.DeleteShortUrl(r.Context(), shortUrlId); err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, "Short URL has deleted with success")
}
