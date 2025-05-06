package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (h apiHandler) registerRoutes() {
	h.r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	h.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	h.r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/login", h.handleGetUserByEmail)

			r.Post("/", h.handleCreateUser)
			r.Get("/", h.handleGetUsers)
			r.Get("/{user_id}", h.handleGetUserById)

			r.Patch("/{user_id}", h.handelUpdateUser)
			r.Delete("/{user_id}", h.handleDeleteUser)
		})

		r.Route("/short", func(r chi.Router) {
			r.Get("/{user_id}", h.handleGetShortUrlsByUser)
			r.Get("/{short_url_id}", h.handleGetShorUrl)
			r.Post("/", h.handleCreateShortUrl)
			r.Patch("/{short_url_id}", h.handleUpdateShortUrl)
			r.Delete("/{short_url_id}", h.handleDeleteShortUrl)
		})
	})

	h.r.Get("/{slug}", h.handleRedirect)
}
