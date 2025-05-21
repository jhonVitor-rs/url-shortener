package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/api/hooks"
	"github.com/jhonVitor-rs/url-shortener/internal/api/middleware"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
)

// =============================================================================
// User Handlers
// =============================================================================

// handleListUsers returns a list of all users
// @Summary List all users
// @Description Returns a list of all users in the system
// @Produce json
// @Success 200 {object} models.Response "List of users"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/users/all [get]
func (h apiHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.user.ListUsers(r.Context())
	hooks.SendResponse(w, http.StatusOK, users, err)
}

// handleGetUser returns details of the authenticated user
// @Summary Get user details
// @Description Returns details of the currently authenticated user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response "User details"
// @Failure 401 {objext} models.Response "Invalid user ID in token"
// @Failure 404 {object} models.Response "User not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/users [get]
func (h apiHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserIdFromContext(r.Context())
	if err != nil {
		hooks.SendResponse(w, http.StatusUnauthorized, nil, err)
		return
	}

	user, err := h.user.GetUser(r.Context(), userId)
	hooks.SendResponse(w, http.StatusOK, user, err)
}

// handleGetUserByEmail authenticates a user by email and returns a JWT token
// @Summary Login user
// @Description Authenticates a user by email and returns a JWT token
// @Accept json
// @Produce json
// @Param credentials body models.GetUserByEmailInput true "User email"
// @Success 200 {object} models.Response "JWT token"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 404 {object} models.Response "User not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/users/login [post]
func (h apiHandler) handleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	body, ok := utils.ParseAndValidate[models.GetUserByEmailInput](w, r)
	if !ok {
		return
	}

	user, err := h.user.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		hooks.SendResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	token, err := middleware.GenerateJWT(user.ID)
	hooks.SendResponse(w, http.StatusCreated, models.Token{JWT: token}, err)
}

// handleCreateUser creates a new user
// @Summary Create a new user
// @Description Creates a new user with the provided information
// @Accept json
// @Produce json
// @Param user body models.CreateUserInput true "User information"
// @Success 200 {object} models.Response "JWT token"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 409 {object} models.Response "Email already in use"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/users [post]
func (h apiHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, ok := utils.ParseAndValidate[models.CreateUserInput](w, r)
	if !ok {
		return
	}

	user, err := h.user.CreateUser(r.Context(), body)
	if err != nil {
		hooks.SendResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	token, err := middleware.GenerateJWT(user.ID)
	hooks.SendResponse(w, http.StatusCreated, models.Token{JWT: token}, err)
}

// handleUpdateUser updates the authenticated user's information
// @Summary Update user information
// @Description Updates the authenticated user's information
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.UpdateUserInput true "User updated information"
// @Success 200 {object} models.Response "User updated successfully"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 401 {object} models.Response "Invalid user ID in token"
// @Failure 404 {object} models.Response "User not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/users [patch]
func (h apiHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserIdFromContext(r.Context())
	if err != nil {
		hooks.SendResponse(w, http.StatusUnauthorized, nil, err)
		return
	}

	body, ok := utils.ParseAndValidate[models.UpdateUserInput](w, r)
	if !ok {
		return
	}

	updatedUser, err := h.user.UpdateUser(r.Context(), userId, body)
	hooks.SendResponse(w, http.StatusOK, updatedUser, err)
}

// handleDeleteUser deletes the authenticated user
// @Summary Delete user
// @Description Deletes the authenticated user's account
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response "User deleted with success"
// @Failure 401 {object} models.Response "Invalid user ID in token"
// @Failure 404 {object} models.Response "User not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/users [delete]
func (h apiHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserIdFromContext(r.Context())
	if err != nil {
		hooks.SendResponse(w, http.StatusUnauthorized, nil, err)
		return
	}

	err = h.user.DeleteUser(r.Context(), userId)
	hooks.SendResponse(w, http.StatusNoContent, nil, err)
}

// =============================================================================
// Short URL Handlers
// =============================================================================

// handleListShortUrlsByUser returns all short URLs created by the authenticated user
// @Summary List user's short URLs
// @Description Returns all short URLs created by the authenticated user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response "List of short URLs"
// @Failure 401 {object} models.Response "Invalid user ID in token"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/short_url [get]
func (h apiHandler) handleListShortUrlsByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserIdFromContext(r.Context())
	if err != nil {
		hooks.SendResponse(w, http.StatusUnauthorized, nil, err)
		return
	}

	shortUrls, err := h.shortUrl.ListShortUrl(r.Context(), userId)
	hooks.SendResponse(w, http.StatusOK, shortUrls, err)
}

// handleGetShortUrl returns details of a specific short URL
// @Summary Get short URL details
// @Description Returns details of a specific short URL by ID
// @Produce json
// @Param short_url_id path string true "Short URL ID"
// @Success 200 {object} models.Response "Short URL details"
// @Failure 404 {object} models.Response "Short URL not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/short_url/{short_url_id} [get]
func (h apiHandler) handleGetShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")

	shortUrl, err := h.shortUrl.GetShortUrl(r.Context(), shortUrlId)
	hooks.SendResponse(w, http.StatusOK, shortUrl, err)
}

// handleCreateShortUrl creates a new short URL
// @Summary Create short URL
// @Description Creates a new short URL for the authenticated user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param url body ports.CreateShortUrlInput true "Original URL information"
// @Success 201 {object} models.Response "Short URL created successfully"
// @Failure 400 {object} models.Response "Invalid input data or URL format"
// @Failure 401 {object} models.Response "Invalid user ID in token"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/short_url [post]
func (h apiHandler) handleCreateShortUrl(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserIdFromContext(r.Context())
	if err != nil {
		hooks.SendResponse(w, http.StatusUnauthorized, nil, err)
		return
	}

	body, ok := utils.ParseAndValidate[models.CreateShortUrlInput](w, r)
	if !ok {
		return
	}

	if ok = utils.ValidateUrl(w, body.OriginalUrl); !ok {
		return
	}

	shortUrl, err := h.shortUrl.CreateShortUrl(r.Context(), userId, body)
	hooks.SendResponse(w, http.StatusCreated, shortUrl, err)
}

// handleUpdateShortUrl updates a specific short URL
// @Summary Update short URL
// @Description Updates a specific short URL by ID
// @Accept json
// @Produce json
// @Param short_url_id path string true "Short URL ID"
// @Param url body ports.UpdateShortUrlInput true "Updated URL information"
// @Success 200 {object} models.Response "Short URL updated successfully"
// @Failure 400 {object} models.Response "Invalid input data"
// @Failure 404 {object} models.Response "Short URL not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/short_url/{short_url_id} [patch]
func (h apiHandler) handleUpdateShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")
	body, ok := utils.ParseAndValidate[models.UpdateShortUrlInput](w, r)
	if !ok {
		return
	}

	shortUrl, err := h.shortUrl.UpdateShortUrl(r.Context(), shortUrlId, body)
	hooks.SendResponse(w, http.StatusOK, shortUrl, err)
}

// handleDeleteShortUrl deletes a specific short URL
// @Summary Delete short URL
// @Description Deletes a specific short URL by ID
// @Produce json
// @Param short_url_id path string true "Short URL ID"
// @Success 200 {object} models.Response "Short URL has deleted with success"
// @Failure 404 {object} models.Response "Short URL not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /api/short_url/{short_url_id} [delete]
func (h apiHandler) handleDeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")

	err := h.shortUrl.DeleteShortUrl(r.Context(), shortUrlId)
	hooks.SendResponse(w, http.StatusNoContent, nil, err)
}

// =============================================================================
// Redirect URL Handler
// =============================================================================

// handleRedirect redirects to the original URL
// @Summary Redirect to original URL
// @Description Redirects to the original URL associated with the provided slug
// @Param slug path string true "Short URL slug"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 404 {object} utils.ErrorResponse "Short URL not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /{slug} [get]
func (h apiHandler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	url, err := h.cache.GetURL(r.Context(), slug)
	if err == nil {
		h.accessCount.IncrementAccess(r.Context(), slug)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	// Not found in cache, fetch from database
	shortUrl, err := h.shortUrl.GetShortUrlBySlug(r.Context(), slug)
	if err != nil {
		hooks.SendResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	// Cache the result and update access statistics
	h.cache.LogRecentAccess(r.Context(), shortUrl)
	h.accessCount.IncrementAccess(r.Context(), slug)
	http.Redirect(w, r, shortUrl.OriginalUrl, http.StatusFound)
}
