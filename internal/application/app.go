package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/diegogrlima/lol-tracker/internal/riot"
)

type App struct {
	router        http.Handler
	serverAddress string
}

func New(
	riotClient *riot.Client,
	serverAddress string,
) *App {
	return &App{
		router:        loadRoutes(riotClient),
		serverAddress: serverAddress,
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:              a.serverAddress,
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
