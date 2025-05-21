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

func TestIntegrationCreateShortUrl(t *testing.T) {
	t.Run("Create short url with success", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		input := models.CreateShortUrlInput{
			OriginalUrl: "https://www.youtube.com/watch?v=-Ka4YKW7RwM&t=537s",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, true, response.Success)

		shortUrl, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		slug, exists := shortUrl["slug"]
		assert.True(t, exists, "Slug field should exist in response data")
		assert.NotEmpty(t, slug, "Slug should not empty")
	})

	t.Run("Error due to lack of token", func(t *testing.T) {
		input := models.CreateShortUrlInput{
			OriginalUrl: "https://www.youtube.com/watch?v=-Ka4YKW7RwM&t=537s",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
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

	t.Run("Invalid input error", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		input := models.CreateShortUrlInput{}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response models.Response
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, false, response.Success)
		assert.Equal(t, response.Error.Message, "Invalid input")
	})
}
