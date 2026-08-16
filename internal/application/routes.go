package application

import (
	"github.com/diegogrlima/lol-tracker/internal/handler"
	"github.com/diegogrlima/lol-tracker/internal/riot"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes(riotClient *riot.Client) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	loadHealthRoutes(router)
	loadPlayerRoutes(router, riotClient)

	return router
}

func loadHealthRoutes(router chi.Router) {
	healthHandler := &handler.Health{}

	router.Route("/health", func(r chi.Router) {
		r.Get("/", healthHandler.GetHealth)
	})
}

func loadPlayerRoutes(router chi.Router, riotClient *riot.Client) {
	playerHandler := handler.NewPlayer(riotClient)

	router.Route("/players", func(router chi.Router) {
		router.Get("/{gameName}/{tagLine}", playerHandler.GetPlayer)
	})
}
