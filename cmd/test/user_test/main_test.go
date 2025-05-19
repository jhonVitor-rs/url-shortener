package user_test

import (
	"context"
	"os"
	"testing"

	"github.com/jhonVitor-rs/url-shortener/cmd/test"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/primary/api"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
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
