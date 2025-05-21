package hooks

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

func SendResponse(w http.ResponseWriter, status int, data interface{}, err error) {
	response := models.Response{
		Success: err == nil,
		Data:    data,
	}

	if err == nil {
		WriteJSON(w, status, response)
		return
	}

	response.Data = nil
	response.Error = processError(err)

	errorStatus := determineErrorStatus(err)

	slog.Error(response.Error.Message, "status", errorStatus, "details", response.Error.Details)

	WriteJSON(w, errorStatus, response)
}

func processError(err error) *models.ErrorData {
	errorData := &models.ErrorData{
		Message: "Unexpected error ocurred",
	}

	var appErr *wraperrors.AppError
	if errors.As(err, &appErr) {
		errorData.Message = appErr.Message
		if appErr.Err != nil {
			errorData.Details = appErr.Err.Error()
		}
	} else if err != nil {
		errorData.Details = err.Error()
	}

	return errorData
}

func determineErrorStatus(err error) int {
	var appErr *wraperrors.AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}

	switch {
	case wraperrors.IsNotFoundError(err):
		return http.StatusNotFound
	case wraperrors.IsValidationError(err):
		return http.StatusBadRequest
	case wraperrors.IsUnauthorizedError(err):
		return http.StatusUnauthorized
	case wraperrors.IsAlreadyExistsError(err) || wraperrors.IsUniqueViolation(err):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
