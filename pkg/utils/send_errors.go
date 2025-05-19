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
	Error   string            `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func SendErrors(w http.ResponseWriter, err error) {
	if errors.As(err, &appErr) {
		slog.Error(appErr.Message, "error", appErr.Err)
		errResp := ErrorResponse{
			Message: appErr.Message,
		}

		if appErr.Err != nil {
			errStr := appErr.Err.Error()
			errResp.Error = errStr
		}

		WriteJSON(w, appErr.Code, errResp)
		return
	}
	errResp := ErrorResponse{
		Message: "Unexpected error",
	}

	if err != nil {
		errStr := err.Error()
		errResp.Error = errStr
	}

	WriteJSON(w, 500, errResp)
}
