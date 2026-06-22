package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Products Routing
	r.Route("/products", func(r chi.Router) {
		// Public
		r.Get("/", h.ListProducts)
		r.Get("/{id}", h.GetProduct)

		// Admin Only
		r.With(AdminOnlyMiddleware).Post("/", h.CreateProduct)
		r.With(AdminOnlyMiddleware).Patch("/{id}", h.UpdateProduct)
		r.With(AdminOnlyMiddleware).Delete("/{id}", h.DeleteProduct)
	})

	// Orders Routing
	r.Route("/orders", func(r chi.Router) {
		// User Auth
		r.With(UserAuthMiddleware).Post("/", h.CreateOrder)
		r.With(UserAuthMiddleware).Get("/{id}", h.GetOrder)
		r.With(UserAuthMiddleware).Delete("/{id}", h.DeleteOrder)

		// Admin Only
		r.With(AdminOnlyMiddleware).Get("/", h.ListOrders)
		r.With(AdminOnlyMiddleware).Patch("/{id}", h.UpdateOrder)
	})

	return r
}
