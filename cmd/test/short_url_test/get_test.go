package shorturltest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhonVitor-rs/url-shortener/cmd/test"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationGetShorUrl(t *testing.T) {
	t.Run("Get all short urls by user with success", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, "/api/short_url/list", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response []*models.ShortUrl
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response)
	})

	t.Run("Get short url by id with success", func(t *testing.T) {
		token, shortUrlId := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/short_url/%s", shortUrlId), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response models.ShortUrl
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.ID)
	})

	t.Run("Redirec with url", func(t *testing.T) {
		token, shortUrlId := setupTestShortUrl(t)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/short_url/%s", shortUrlId), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		var shortUrl models.ShortUrl
		err := json.NewDecoder(recorder.Body).Decode(&shortUrl)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/redirect/%s", shortUrl.Slug), nil)
		req.Header.Set("Content_type", "application/json")
		recorder = httptest.NewRecorder()
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

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		var response utils.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "something went wrong", response.Message)
	})

	t.Run("Erro to get a short url with invalid slug", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/redirect/invalid-slug", nil)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)

		var response utils.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "short URL with slug 'invalid-slug' not found", response.Message)
	})
}
