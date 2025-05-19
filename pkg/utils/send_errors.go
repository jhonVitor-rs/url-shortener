package utils

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

var appErr *wraperrors.AppError

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
	Error   *error            `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func SendErrors(w http.ResponseWriter, err error) {
	if errors.As(err, &appErr) {
		slog.Error(appErr.Message, "error", appErr.Err)
		WriteJSON(w, appErr.Code, ErrorResponse{
			Message: appErr.Message, Error: &appErr.Err,
		})
		return
	}
	WriteJSON(w, 500, ErrorResponse{
		Message: "Unexpected error", Error: &err,
	})
}
