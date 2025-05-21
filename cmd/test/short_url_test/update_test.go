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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationUpdateShortUrl(t *testing.T) {
	t.Run("Update short url with success", func(t *testing.T) {
		token, shortUrl := setupTestShortUrl(t)

		originalUrl := "https://www.youtube.com/watch?v=g5ZUG1gKZpE"
		input := models.UpdateShortUrlInput{
			OriginalUrl: ptr(originalUrl),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/short_url/%s", shortUrl.ID), bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		shortUrlData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		newUrl, exists := shortUrlData["original_url"]
		assert.True(t, exists, "Original URL should exist in response data")
		assert.Equal(t, originalUrl, newUrl)
	})

	t.Run("Failed to update to missing token", func(t *testing.T) {
		_, shortUrl := setupTestShortUrl(t)

		originalUrl := "https://www.youtube.com/watch?v=g5ZUG1gKZpE"
		input := models.UpdateShortUrlInput{
			OriginalUrl: ptr(originalUrl),
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/short_url/%s", shortUrl.ID), bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, response.Error.Message, "Missing or invalid token")
	})
}

func ptr(s string) *string {
	return &s
}
