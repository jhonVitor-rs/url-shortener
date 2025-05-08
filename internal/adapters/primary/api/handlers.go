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
	body, ok := utils.ParseAndValidate[ports.UpdateUserInput](w, r)
	if !ok {
		return
	}

	userId, err := h.user.UpdateUser(r.Context(), body)
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
