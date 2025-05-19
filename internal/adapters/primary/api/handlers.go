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

// =============================================================================
// User Handlers
// =============================================================================

// handleCreateUser creates a new user
// @Summary Create a new user
// @Description Creates a new user with the provided information
// @Accept json
// @Produce json
// @Param user body ports.CreateUserInput true "User information"
// @Success 201 {object} models.Response "User created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 409 {object} utils.ErrorResponse "Email already in use"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/users [post]
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

	w.WriteHeader(http.StatusCreated)
	utils.SendJSON(w, models.Response{ID: userId})
}

// handleListUsers returns a list of all users
// @Summary List all users
// @Description Returns a list of all users in the system
// @Produce json
// @Success 200 {array} models.User "List of users"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/users/all [get]
func (h apiHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.user.ListUsers(r.Context())
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, users)
}

// handleGetUser returns details of the authenticated user
// @Summary Get user details
// @Description Returns details of the currently authenticated user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "User details"
// @Failure 401 {string} string "Invalid user ID in token"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/users [get]
func (h apiHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	user, err := h.user.GetUser(r.Context(), userId)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, user)
}

// handleGetUserByEmail authenticates a user by email and returns a JWT token
// @Summary Login user
// @Description Authenticates a user by email and returns a JWT token
// @Accept json
// @Produce json
// @Param credentials body ports.GetUserByEmailInput true "User email"
// @Success 200 {object} models.Token "JWT token"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/users/login [post]
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
		slog.Error("Failed to generate token", "error", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, models.Token{JWT: token})
}

// handleUpdateUser updates the authenticated user's information
// @Summary Update user information
// @Description Updates the authenticated user's information
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body ports.UpdateUserInput true "User updated information"
// @Success 200 {object} models.Response "User updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 401 {string} string "Invalid user ID in token"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/users [patch]
func (h apiHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid user ID in token")
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

// handleDeleteUser deletes the authenticated user
// @Summary Delete user
// @Description Deletes the authenticated user's account
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "User deleted with success"
// @Failure 401 {string} string "Invalid user ID in token"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/users [delete]
func (h apiHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	if err := h.user.DeleteUser(r.Context(), userId); err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, "User deleted with success")
}

// =============================================================================
// Short URL Handlers
// =============================================================================

// handleCreateShortUrl creates a new short URL
// @Summary Create short URL
// @Description Creates a new short URL for the authenticated user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param url body ports.CreateShortUrlInput true "Original URL information"
// @Success 201 {object} models.Response "Short URL created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data or URL format"
// @Failure 401 {string} string "Invalid user ID in token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/short_url [post]
func (h apiHandler) handleCreateShortUrl(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid user ID in token")
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

	w.WriteHeader(http.StatusCreated)
	utils.SendJSON(w, models.Response{ID: shortUrlId})
}

// handleListShortUrlsByUser returns all short URLs created by the authenticated user
// @Summary List user's short URLs
// @Description Returns all short URLs created by the authenticated user
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ShortUrl "List of short URLs"
// @Failure 401 {string} string "Invalid user ID in token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/short_url [get]
func (h apiHandler) handleListShortUrlsByUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserIdFromContex(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	shortUrls, err := h.shortUrl.ListShortUrl(r.Context(), userId)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, shortUrls)
}

// handleGetShortUrl returns details of a specific short URL
// @Summary Get short URL details
// @Description Returns details of a specific short URL by ID
// @Produce json
// @Param short_url_id path string true "Short URL ID"
// @Success 200 {object} models.ShortUrl "Short URL details"
// @Failure 404 {object} utils.ErrorResponse "Short URL not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/short_url/{short_url_id} [get]
func (h apiHandler) handleGetShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")

	shortUrl, err := h.shortUrl.GetShortUrl(r.Context(), shortUrlId)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, shortUrl)
}

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

	// Try to get URL from Redis cache first
	url, err := rdstore.GetUrl(r.Context(), h.rdb, slug)
	if err == nil {
		// URL found in cache, increment access counter and redirect
		rdstore.IncrementAccess(r.Context(), h.rdb, slug)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	// Not found in cache, fetch from database
	shortUrl, err := h.shortUrl.GetShortUrlBySlug(r.Context(), slug)
	if err != nil {
		utils.SendErrors(w, err)
		return
	}

	// Cache the result and update access statistics
	rdstore.LogRecentAccess(r.Context(), h.rdb, shortUrl)
	rdstore.IncrementAccess(r.Context(), h.rdb, slug)
	http.Redirect(w, r, shortUrl.OriginalUrl, http.StatusFound)
}

// handleUpdateShortUrl updates a specific short URL
// @Summary Update short URL
// @Description Updates a specific short URL by ID
// @Accept json
// @Produce json
// @Param short_url_id path string true "Short URL ID"
// @Param url body ports.UpdateShortUrlInput true "Updated URL information"
// @Success 200 {object} models.Response "Short URL updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 404 {object} utils.ErrorResponse "Short URL not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/short_url/{short_url_id} [patch]
func (h apiHandler) handleUpdateShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")
	body, ok := utils.ParseAndValidate[ports.UpdateShortUrlInput](w, r)
	if !ok {
		return
	}

	shortUrlId, err := h.shortUrl.UpdateShortUrl(r.Context(), shortUrlId, body)
	if err != nil {
		slog.Error("error to update short url", "error", err)
		utils.SendErrors(w, err)
		return
	}

	slog.Info(shortUrlId)

	utils.SendJSON(w, models.Response{ID: shortUrlId})
}

// handleDeleteShortUrl deletes a specific short URL
// @Summary Delete short URL
// @Description Deletes a specific short URL by ID
// @Produce json
// @Param short_url_id path string true "Short URL ID"
// @Success 200 {string} string "Short URL has deleted with success"
// @Failure 404 {object} utils.ErrorResponse "Short URL not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/short_url/{short_url_id} [delete]
func (h apiHandler) handleDeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	shortUrlId := chi.URLParam(r, "short_url_id")
	if err := h.shortUrl.DeleteShortUrl(r.Context(), shortUrlId); err != nil {
		utils.SendErrors(w, err)
		return
	}

	utils.SendJSON(w, "Short URL has deleted with success")
}
