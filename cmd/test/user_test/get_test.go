package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestIntegrationGetUser(t *testing.T) {
	t.Run("Get users with success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users/all", nil)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response []*models.User
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response)
	})

	t.Run("Fail to get user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var response utils.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Missing or invalid token", response.Message)
	})

	t.Run("Login user", func(t *testing.T) {
		setupTestUser(t)

		input := ports.GetUserByEmailInput{
			Email: "jhon.doe@email.com",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var token models.Token
		err = json.NewDecoder(recorder.Body).Decode(&token)
		require.NoError(t, err)

		assert.NotEmpty(t, token.JWT)
	})

	t.Run("Failed to login with invalid email", func(t *testing.T) {
		input := ports.GetUserByEmailInput{
			Email: "jhon.due@email.com",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)

		var response utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "User not found", response.Message)
	})

	t.Run("Failed to login with invalid input", func(t *testing.T) {
		input := ports.GetUserByEmailInput{}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid input", response.Message)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "Email", response.Errors[0].Field)
		assert.Equal(t, "required", response.Errors[0].Tag)
	})

	t.Run("Get user with success", func(t *testing.T) {
		token := setupTestUser(t)

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()

		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.User
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response)
	})
}
