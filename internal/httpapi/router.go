package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(playerHandler *PlayerHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(15 * time.Second))

	router.Get("/health", getHealth)
	router.Get("/health/", getHealth)
	router.Get("/players/{gameName}/{tagLine}", playerHandler.GetByRiotID)

	return router
}

func getHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"health": map[string]any{
			"message": "API funcionando",
			"status":  true,
		},
	})
}
