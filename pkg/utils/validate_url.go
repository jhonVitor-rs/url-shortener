package utils

import (
	"net/http"
	"net/url"

	"github.com/jhonVitor-rs/url-shortener/internal/api/hooks"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

func ValidateUrl(w http.ResponseWriter, urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		hooks.SendResponse(w, http.StatusBadRequest, nil, wraperrors.ValidationErr("Invalid URL"))
		return false
	}
	return true
}
