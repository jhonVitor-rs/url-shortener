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
		writeJSON(w, appErr.Code, ErrorResponse{
			Message: appErr.Message, Error: &appErr.Err,
		})
		return
	}
	writeJSON(w, 500, ErrorResponse{
		Message: "Unexpected error", Error: &err,
	})
}
