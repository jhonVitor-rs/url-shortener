package shorturltest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jhonVitor-rs/url-shortener/cmd/test"
	api "github.com/jhonVitor-rs/url-shortener/internal/api/server"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, rdb := test.Config(ctx)

	if pool == nil || rdb == nil {
		os.Exit(1)
	}

	defer func() {
		if pool != nil {
			pool.Close()
		}
		if rdb != nil {
			rdb.Close()
		}
	}()

	test.SetHandler(api.NewApiHandler(pgstore.New(pool), rdb))

	exitCode := m.Run()

	if pool != nil {
		_, err := pool.Exec(ctx, `
			DELETE FROM short_urls;
			DELETE FROM users;
		`)
		if err != nil {
			os.Exit(exitCode)
		}
	}

	os.Exit(exitCode)
}

func setupTestShortUrl(t *testing.T) (string, *models.ShortUrl) {
	email := fmt.Sprintf("jhon.doe+%d@email.com", time.Now().UnixNano())

	input := models.CreateUserInput{Name: "Jhon", Email: email}
	payload, err := json.Marshal(input)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	test.Handler().ServeHTTP(recorder, req)
	if recorder.Code != http.StatusCreated && recorder.Code != http.StatusConflict {
		t.Fatalf("unexpected status code when creating user: %d", recorder.Code)
	}

	loginInput := models.GetUserByEmailInput{Email: input.Email}
	loginPayload, err := json.Marshal(loginInput)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()

	test.Handler().ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)

	var responseToken models.Response
	err = json.NewDecoder(recorder.Body).Decode(&responseToken)
	require.NoError(t, err)

	tokenData, ok := responseToken.Data.(map[string]interface{})
	require.True(t, ok, "Response data should be a map")

	jwtToken, exists := tokenData["jwt"]
	require.True(t, exists, "JWT field should exist in response")
	token := jwtToken.(string)
	require.NotEmpty(t, token, "JWT token shoud not be empty")

	shortUrlInput := models.CreateShortUrlInput{
		OriginalUrl: "https://www.youtube.com/watch?v=-Ka4YKW7RwM&t=537s",
	}
	payload, err = json.Marshal(shortUrlInput)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	recorder = httptest.NewRecorder()

	test.Handler().ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code, "Failed to create short URL")

	var responseUrl models.Response
	err = json.NewDecoder(recorder.Body).Decode(&responseUrl)
	require.NoError(t, err)

	shortUrl := models.ShortUrl{}
	shortUrlData, ok := responseUrl.Data.(map[string]interface{})
	require.True(t, ok, "Response data should be a map")

	if id, exists := shortUrlData["id"]; exists {
		shortUrl.ID = id.(string)
	}
	if slug, exists := shortUrlData["slug"]; exists {
		shortUrl.Slug = slug.(string)
	}
	if originalUrl, exists := shortUrlData["original_url"]; exists {
		shortUrl.OriginalUrl = originalUrl.(string)
	}
	if userId, exists := shortUrlData["user_id"]; exists {
		shortUrl.UserID = userId.(string)
	}
	if createdAt, exists := shortUrlData["created_at"]; exists {
		createdAtStr := createdAt.(string)
		parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
		require.NoError(t, err, "Failed to parse created_at time")
		shortUrl.CreatedAt = parsedTime
	}

	require.NotEmpty(t, shortUrl.ID, "Short URL ID should not be empty")
	require.NotEmpty(t, shortUrl.Slug, "Short URL slug should not be empty")

	return token, &shortUrl

}
