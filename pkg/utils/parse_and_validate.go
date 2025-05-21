package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jhonVitor-rs/url-shortener/internal/api/hooks"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
)

var validate = validator.New()

func ParseAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Failed to decode JSON", "error", err)
		hooks.WriteJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error: &models.ErrorData{
				Message: "Invalid JSON",
				Details: err.Error(),
			},
		})
		return nil, false
	}

	if err := validate.Struct(body); err != nil {
		slog.Error("Failed to validate struct", "error", err)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errs []models.ValidationError
			for _, verr := range validationErrors {
				errs = append(errs, models.ValidationError{
					Field: verr.Field(),
					Tag:   verr.Tag(),
					Value: fmt.Sprintf("%v", verr.Value()),
				})
			}
			hooks.WriteJSON(w, http.StatusBadRequest, models.Response{
				Success: false,
				Error: &models.ErrorData{
					Message: "Invalid input",
					Errors:  errs,
				},
			})
			return nil, false
		}

		hooks.WriteJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Error: &models.ErrorData{
				Message: "Invalid input",
				Details: err.Error(),
			},
		})
		return nil, false
	}

	return &body, true
}
