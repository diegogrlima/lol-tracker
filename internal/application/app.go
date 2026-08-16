package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router      http.Handler
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *App {
	app := &App{
		router:      loadRoutes(),
		redisClient: redisClient,
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {

		return fmt.Errorf("Erro ao iniciar o servidor: %w", err)
	}

	return nil
}
