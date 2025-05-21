package shorturltest_test

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

func TestIntegrationGetShorUrl(t *testing.T) {
	t.Run("Get all short urls by user with success", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, "/api/short_url", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		shortUrls, ok := response.Data.([]interface{})
		assert.True(t, ok, "Response data should be a slice")
		assert.NotEmpty(t, shortUrls, "Short URLs list should not be empty")
	})

	t.Run("Get short url by id with success", func(t *testing.T) {
		token, shortUrl := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/short_url/%s", shortUrl.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		shortUrlMap, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		slug, exists := shortUrlMap["slug"]
		assert.True(t, exists, "Slug field should exist in response data")
		assert.NotEmpty(t, slug, "Slug should not empty")
	})

	t.Run("Redirec with url", func(t *testing.T) {
		_, shortUrl := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", shortUrl.Slug), nil)
		req.Header.Set("Content_type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusFound, recorder.Code)
	})

	t.Run("Erro to get short url with invalid id", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, "/api/short_url/invalid-id", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, response.Error.Message, "Invalid short URL ID format")
	})

	t.Run("Erro to get a short url with invalid slug", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invalid-slug", nil)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)

		var response models.Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, response.Error.Message, "Short URL not found")
	})
}
