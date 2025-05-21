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

func TestIntegrationUpdateUser(t *testing.T) {
	t.Run("Update user with success", func(t *testing.T) {
		token := setupTestUser(t)

		input := models.UpdateUserInput{
			Name:  ptr("Joao"),
			Email: ptr("joao.siben@email.com"),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Failed to update user, status code: %d", recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		userData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Respose data should be a map")
		userName, exists := userData["name"]
		assert.True(t, exists, "Name field should exist in response data")
		assert.Equal(t, "Joao", userName)
	})

	t.Run("Failed to update user", func(t *testing.T) {
		input := models.UpdateUserInput{
			Name:  ptr("Joao"),
			Email: ptr("joao.siben@email.com"),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Missing or invalid token", response.Error.Message)
	})
}

func ptr(s string) *string {
	return &s
}
