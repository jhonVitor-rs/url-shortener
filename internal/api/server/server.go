package server

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/services"
	"github.com/jhonVitor-rs/url-shortener/internal/data/db/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/data/infra"
	"github.com/redis/go-redis/v9"
)

type apiHandler struct {
	r           *chi.Mux
	mu          *sync.Mutex
	user        ports.UserUseCase
	shortUrl    ports.ShortUrlUseCase
	cache       *infra.URLCache
	accessCount *infra.AccessCounter
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewApiHandler(q *pgstore.Queries, rdb *redis.Client) http.Handler {
	a := apiHandler{
		r:           chi.NewRouter(),
		mu:          &sync.Mutex{},
		user:        services.NewUserService(q),
		shortUrl:    services.NewShortUrlService(q),
		cache:       infra.NewURLCache(rdb),
		accessCount: infra.NewAccessCounter(rdb),
	}
	a.registerRoutes()

	return a
}
