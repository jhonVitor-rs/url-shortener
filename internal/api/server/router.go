package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	my_middleware "github.com/jhonVitor-rs/url-shortener/internal/api/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *apiHandler) registerRoutes() {
	h.r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	h.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	h.r.Get("/swagger/*", httpSwagger.Handler())

	h.r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/all", h.handleListUsers)
			r.Post("/login", h.handleGetUserByEmail)
			r.Post("/", h.handleCreateUser)

			r.Group(func(r chi.Router) {
				r.Use(my_middleware.JWTAuth)

				r.Get("/", h.handleGetUser)
				r.Patch("/", h.handleUpdateUser)
				r.Delete("/", h.handleDeleteUser)
			})
		})

		r.Route("/short_url", func(r chi.Router) {
			r.Use(my_middleware.JWTAuth)

			r.Get("/", h.handleListShortUrlsByUser)
			r.Get("/{short_url_id}", h.handleGetShortUrl)
			r.Post("/", h.handleCreateShortUrl)
			r.Patch("/{short_url_id}", h.handleUpdateShortUrl)
			r.Delete("/{short_url_id}", h.handleDeleteShortUrl)
		})
	})

	h.r.Get("/{slug}", h.handleRedirect)
}
