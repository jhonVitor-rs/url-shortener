package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func ParseAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Failed to decode JSON", "error", err)
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid JSON" + err.Error(),
		})
		return nil, false
	}

	if err := validate.Struct(body); err != nil {
		slog.Error("Failed to validate struct", "error", err)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errs []ValidationError
			for _, verr := range validationErrors {
				errs = append(errs, ValidationError{
					Field: verr.Field(),
					Tag:   verr.Tag(),
					Value: fmt.Sprintf("%v", verr.Value()),
				})
			}
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid input",
				Errors:  errs,
			})
			return nil, false
		}

		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid input: " + err.Error(),
		})
		return nil, false
	}

	return &body, true
}
