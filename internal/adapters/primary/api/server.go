package api

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/pgstore"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/secondary/persistence/repositories"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/ports"
	"github.com/jhonVitor-rs/url-shortener/internal/core/usecases/services"
)

type apiHandler struct {
	r  *chi.Mux
	mu *sync.Mutex

	user     ports.UserUseCase
	shortUrl ports.ShortUrlUseCase
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewApiHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		r:        chi.NewRouter(),
		mu:       &sync.Mutex{},
		user:     services.NewUserService(repositories.NewUserRepository(q)),
		shortUrl: services.NewShortUrlService(repositories.NewSHortUrlRepository(q)),
	}
	a.registerRoutes()

	return a
}
