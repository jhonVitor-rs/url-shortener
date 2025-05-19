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
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/primary/api"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
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

func setupTestShortUrl(t *testing.T) (string, string) {
	email := fmt.Sprintf("jhon.doe+%d@email.com", time.Now().UnixNano())
	input := ports.CreateUserInput{Name: "Jhon", Email: email}
	payload, err := json.Marshal(input)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	test.Handler().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated && recorder.Code != http.StatusConflict {
		t.Fatalf("unexpected status code when creating user: %d", recorder.Code)
	}

	loginInput := ports.GetUserByEmailInput{Email: input.Email}
	loginPayload, err := json.Marshal(loginInput)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	test.Handler().ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	var token models.Token
	err = json.NewDecoder(recorder.Body).Decode(&token)
	require.NoError(t, err)

	shortUrlInput := ports.CreateShortUrlInput{
		OriginalUrl: "https://www.youtube.com/watch?v=-Ka4YKW7RwM&t=537s",
	}
	payload, err = json.Marshal(shortUrlInput)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/api/short_url", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.JWT))
	recorder = httptest.NewRecorder()
	test.Handler().ServeHTTP(recorder, req)

	var shortUrlId models.Response
	err = json.NewDecoder(recorder.Body).Decode(&shortUrlId)
	require.NoError(t, err)

	return token.JWT, shortUrlId.ID
}
