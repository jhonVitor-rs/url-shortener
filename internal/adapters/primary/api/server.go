package api

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/infra"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/services"
	"github.com/redis/go-redis/v9"
)

type apiHandler struct {
	r  *chi.Mux
	mu *sync.Mutex

	user     ports.UserUseCase
	shortUrl ports.ShortUrlUseCase

	rdb *redis.Client
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewApiHandler(q *pgstore.Queries, rdb *redis.Client) http.Handler {
	a := apiHandler{
		r:        chi.NewRouter(),
		mu:       &sync.Mutex{},
		user:     services.NewUserService(infra.NewUserRepository(q)),
		shortUrl: services.NewShortUrlService(infra.NewSHortUrlRepository(q)),
		rdb:      rdb,
	}
	a.registerRoutes()

	return a
}
