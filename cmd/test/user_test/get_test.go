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

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		users, ok := response.Data.([]interface{})
		assert.True(t, ok, "Response data should be a slice")
		assert.NotEmpty(t, users, "Users list should not be empty")
	})

	t.Run("Fail to get user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Missing or invalid token", response.Error.Message)
	})

	t.Run("Login user", func(t *testing.T) {
		setupTestUser(t)

		input := models.GetUserByEmailInput{
			Email: "jhon.doe@email.com",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		token, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		jwt, exists := token["jwt"]
		assert.True(t, exists, "JWT field should exist in response data")
		assert.NotEmpty(t, jwt, "JWT token should not be empty")
	})

	t.Run("Failed to login with invalid email", func(t *testing.T) {
		input := models.GetUserByEmailInput{
			Email: "jhon.due@email.com",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, "User not found", response.Error.Message)
	})

	t.Run("Failed to login with invalid input", func(t *testing.T) {
		input := models.GetUserByEmailInput{}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid input", response.Error.Message)
	})

	t.Run("Get user with success", func(t *testing.T) {
		token := setupTestUser(t)

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		recorder := httptest.NewRecorder()

		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		userData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		_, exists := userData["email"]
		assert.True(t, exists, "Email field should exist in response data")
	})
}
