package shorturltest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jhonVitor-rs/url-shortener/cmd/test"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDeleteShortUrl(t *testing.T) {
	t.Run("Delete short url with success", func(t *testing.T) {
		token, shortUrlId := setupTestShortUrl(t)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/short_url/%s", shortUrlId), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		test.Handler().ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
