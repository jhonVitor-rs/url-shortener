package utils

import (
	"errors"
	"log/slog"
	"net/http"

	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

var appErr *wraperrors.AppError

func SendErrors(w http.ResponseWriter, err error) {
	if errors.As(err, &appErr) {
		slog.Error(appErr.Message, "error", appErr.Err)
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	http.Error(w, "unexpected error", http.StatusInternalServerError)
}
