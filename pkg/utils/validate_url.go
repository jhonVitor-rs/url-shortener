package utils

import (
	"net/http"
	"net/url"
)

func ValidateUrl(w http.ResponseWriter, urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return false
	}
	return true
}
