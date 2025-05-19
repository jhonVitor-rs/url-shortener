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

func TestIntegrationCreateShortUrl(t *testing.T) {
	t.Run("Create short url with success", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		input := ports.CreateShortUrlInput{
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

		assert.NotEmpty(t, response.ID)
	})

	t.Run("Error due to lack of token", func(t *testing.T) {
		input := ports.CreateShortUrlInput{
			OriginalUrl: "https://www.youtube.com/watch?v=-Ka4YKW7RwM&t=537s",
		}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		var response utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Missing or invalid token", response.Message)
	})

	t.Run("Invalid input error", func(t *testing.T) {
		token, _ := setupTestShortUrl(t)

		input := ports.CreateShortUrlInput{}
		payload, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var response utils.ErrorResponse
		err = json.NewDecoder(recorder.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid input", response.Message)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "OriginalUrl", response.Errors[0].Field)
		assert.Equal(t, "required", response.Errors[0].Tag)
	})
}
