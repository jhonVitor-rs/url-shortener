package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jhonVitor-rs/url-shortener/internal/adapters/primary/middlewares"
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
			r.Post("/login", h.handleGetUserByEmail)
			r.Get("/all", h.handleListUsers)

			r.Post("/", h.handleCreateUser)

			r.Group(func(r chi.Router) {
				r.Use(middlewares.JWTAuth)

				r.Get("/", h.handleGetUser)

				r.Patch("/", h.handleUpdateUser)
				r.Delete("/", h.handleDeleteUser)
			})
		})

		r.Route("/short_url", func(r chi.Router) {
			r.Use(middlewares.JWTAuth)

			r.Get("/{short_url_id}", h.handleGetShortUrl)
			r.Get("/list", h.handleListShortUrlsByUser)

			r.Post("/", h.handleCreateShortUrl)
			r.Patch("/{short_url_id}", h.handleUpdateShortUrl)
			r.Delete("/{short_url_id}", h.handleDeleteShortUrl)
		})
	})

	h.r.Get("/redirect/{slug}", h.handleRedirect)
}
