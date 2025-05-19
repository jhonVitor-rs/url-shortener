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

func TestIntegrationUpdateUser(t *testing.T) {
	createInput := ports.CreateUserInput{
		Name:  "Jhon Doe",
		Email: "jhon.doe@email.com",
	}
	payload, err := json.Marshal(createInput)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	test.Handler().ServeHTTP(recorder, req)

	loginInput := ports.GetUserByEmailInput{
		Email: "jhon.doe@email.com",
	}
	payload, err = json.Marshal(loginInput)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	test.Handler().ServeHTTP(recorder, req)

	var token models.Token
	err = json.NewDecoder(recorder.Body).Decode(&token)
	require.NoError(t, err)

	t.Run("Update user with success", func(t *testing.T) {
		input := ports.UpdateUserInput{
			Name:  ptr("Joao"),
			Email: ptr("joao.siben@email.com"),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.JWT))
		recorder = httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, "Failed to update user, status code: %d", recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("Failed to update user", func(t *testing.T) {
		input := ports.UpdateUserInput{
			Name:  ptr("Joao"),
			Email: ptr("joao.siben@email.com"),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder = httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var response utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "Missing or invalid token", response.Message)
	})
}

func ptr(s string) *string {
	return &s
}
