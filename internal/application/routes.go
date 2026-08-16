package application

import (
	"github.com/diegogrlima/lol-tracker/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	loadHealthRoutes(router)
	loadPlayerRoutes(router)

	return router
}

func loadHealthRoutes(router chi.Router) {
	router.Get("/health", handler.GetHealth)
}

func loadPlayerRoutes(router chi.Router) {
	playerHandler := &handler.Player{}

	router.Route("/players", func(router chi.Router) {
		router.Get("/{riotID}", playerHandler.GetPlayer)
	})
}
