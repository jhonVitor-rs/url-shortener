package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func parseAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return nil, false
	}

	if err := validate.Struct(body); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return nil, false
	}

	return &body, true
}

func validateUrl(w http.ResponseWriter, urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return false
	}
	return true
}

func sendJSON(w http.ResponseWriter, rawData any) {
	data, err := json.Marshal(rawData)
	if err != nil {
		http.Error(w, "failed to serialize json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
