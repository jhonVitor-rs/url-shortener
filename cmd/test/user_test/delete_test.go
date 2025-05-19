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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeleteUser(t *testing.T) {
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

	t.Run("Delete user with success", func(t *testing.T) {
		req = httptest.NewRequest(http.MethodDelete, "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.JWT))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		
	})
}