package api

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/store/pgstore"
)

type apiHandler struct {
	q  *pgstore.Queries
	r  *chi.Mux
	mu *sync.Mutex
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		q:  q,
		r:  chi.NewRouter(),
		mu: &sync.Mutex{},
	}
	a.registerRoutes()
	return a
}
