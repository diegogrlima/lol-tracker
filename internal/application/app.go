package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/diegogrlima/lol-tracker/internal/riot"
	"github.com/redis/go-redis/v9"
)

type App struct {
	router      http.Handler
	redisClient *redis.Client
}

func New(redisClient *redis.Client, riotClient *riot.Client) *App {
	return &App{
		router:      loadRoutes(riotClient),
		redisClient: redisClient,
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           a.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		serverError <- server.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("start HTTP server: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}

	return nil
}
