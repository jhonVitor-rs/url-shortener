package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhonVitor-rs/url-shortener/cmd/test"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateUser(t *testing.T) {
	t.Run("Create user with success", func(t *testing.T) {
		input := models.CreateUserInput{
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
		assert.Equal(t, true, response.Success)

		token, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		jwt, exists := token["jwt"]
		assert.True(t, exists, "JWT field should exist in response data")
		assert.NotEmpty(t, jwt, "JWT token should not be empty")
	})

	t.Run("Error to create with the same email", func(t *testing.T) {
		input := models.CreateUserInput{
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

		var resp models.Response
		err = json.NewDecoder(recorder.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, false, resp.Success)
		assert.Equal(t, resp.Error.Message, "Email already in use")
	})

	t.Run("Error to validate request body", func(t *testing.T) {
		input := models.CreateUserInput{Name: "Jhon Doe"}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var resp models.Response
		err = json.NewDecoder(recorder.Body).Decode(&resp)
		require.NoError(t, err)

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, resp.Error.Message, "Invalid input")
		assert.Len(t, resp.Error.Errors, 1)
	})
}
