package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func ParseAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
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
