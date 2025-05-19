package shorturltest_test

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

func TestIntegrationUpdateShortUrl(t *testing.T) {
	t.Run("Update short url with success", func(t *testing.T) {
		token, shortUrlId := setupTestShortUrl(t)

		originalUrl := "https://www.youtube.com/watch?v=g5ZUG1gKZpE"
		input := ports.UpdateShortUrlInput{
			OriginalUrl: ptr(originalUrl),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/short_url/%s", shortUrlId), bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.ID)
	})

	t.Run("Failed to update to missing token", func(t *testing.T) {
		_, shortUrlId := setupTestShortUrl(t)

		originalUrl := "https://www.youtube.com/watch?v=g5ZUG1gKZpE"
		input := ports.UpdateShortUrlInput{
			OriginalUrl: ptr(originalUrl),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/short_url/%s", shortUrlId), bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
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
