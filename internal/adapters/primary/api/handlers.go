package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	userId := chi.URLParam(r, "user_id")

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

	utils.SendJSON(w, user)
}

func (h apiHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")
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
	userId := chi.URLParam(r, "user_id")
	if err := h.user.DeleteUser(r.Context(), userId); err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, "User deleted with success")
}

// Short URL handlers -----------------------------------------------
func (h apiHandler) handleCreateShortUrl(w http.ResponseWriter, r *http.Request) {
	body, ok := utils.ParseAndValidate[ports.CreateShortUrlInput](w, r)
	if !ok {
		return
	}

	if ok = utils.ValidateUrl(w, body.OriginalUrl); !ok {
		return
	}

	shortUrlId, err := h.shortUrl.CreateShortUrl(r.Context(), body)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, models.Response{ID: shortUrlId})
}

func (h apiHandler) handleListShortUrlsByUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")

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

	shortUrl, err := h.shortUrl.GetShortUrlBySlug(r.Context(), slug)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

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
