package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var validate = validator.New()

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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ParseAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid JSON" + err.Error(),
		})
		return nil, false
	}

	if err := validate.Struct(body); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errs []ValidationError
			for _, verr := range validationErrors {
				errs = append(errs, ValidationError{
					Field: verr.Field(),
					Tag:   verr.Tag(),
					Value: fmt.Sprintf("%v", verr.Value()),
				})
			}
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid input",
				Errors:  errs,
			})
			return nil, false
		}

		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid input: " + err.Error(),
		})
		return nil, false
	}

	return &body, true
}
