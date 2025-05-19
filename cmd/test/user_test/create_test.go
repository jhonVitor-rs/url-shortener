package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhonVitor-rs/url-shortener/cmd/test"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateUser(t *testing.T) {
	t.Run("Create user with success", func(t *testing.T) {
		input := ports.CreateUserInput{
			Name:  "Jhon Doe",
			Email: "jhon.doe@email.com",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("Error to create with the same email", func(t *testing.T) {
		input := ports.CreateUserInput{
			Name:  "Jhon Doe",
			Email: "jhon.doe@email.com",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusConflict, recorder.Code)

		var resp utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "email already in use", resp.Message)
	})

	t.Run("Error to validate request body", func(t *testing.T) {
		input := ports.CreateUserInput{Name: "Jhon Doe"}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var resp utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, "Invalid input", resp.Message)
		assert.Len(t, resp.Errors, 1)
		assert.Equal(t, "Email", resp.Errors[0].Field)
		assert.Equal(t, "required", resp.Errors[0].Tag)
	})
}
