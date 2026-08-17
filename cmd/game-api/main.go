package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegogrlima/lol-tracker/internal/champion"
	"github.com/diegogrlima/lol-tracker/internal/config"
	ddragonadapter "github.com/diegogrlima/lol-tracker/internal/platform/ddragon"
	"github.com/diegogrlima/lol-tracker/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("game service stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.LoadGame()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	dataDragonClient, err := ddragonadapter.NewClient(
		cfg.DataDragonBaseURL,
		cfg.DataDragonVersion,
		cfg.DataDragonLocale,
	)
	if err != nil {
		return fmt.Errorf("initialize Data Dragon client: %w", err)
	}

	championService := champion.NewService(dataDragonClient)
	championHandler := champion.NewHandler(championService, logger)
	router := server.NewRouter(map[string]http.Handler{
		"/champions": championHandler.Routes(),
	})
	gameServer := server.New(cfg.ServerAddress, router)

	logger.Info("game service started", "address", cfg.ServerAddress)
	if err := gameServer.Start(ctx); err != nil {
		return fmt.Errorf("run game service: %w", err)
	}

	logger.Info("game service stopped gracefully")
	return nil
}
