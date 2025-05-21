package user_test

import (
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

func TestIntegrationDeleteUser(t *testing.T) {
	t.Run("Delete user with success", func(t *testing.T) {
		token := setupTestUser(t)

		req := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
		assert.Equal(t, true, response.Success)
	})
}
