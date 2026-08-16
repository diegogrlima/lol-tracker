package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegogrlima/lol-tracker/internal/config"
	"github.com/diegogrlima/lol-tracker/internal/match"
	riotadapter "github.com/diegogrlima/lol-tracker/internal/platform/riot"
	"github.com/diegogrlima/lol-tracker/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("match service stopped", "error", err)
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

	cfg, err := config.LoadMatch()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	riotClient, err := riotadapter.NewClient(cfg.RiotAPIKey, cfg.RiotRegion)
	if err != nil {
		return fmt.Errorf("initialize Riot client: %w", err)
	}

	matchService := match.NewService(riotClient)
	matchHandler := match.NewHandler(matchService, logger)
	router := server.NewRouter(map[string]http.Handler{
		"/matches": matchHandler.Routes(),
	})
	matchServer := server.New(cfg.ServerAddress, router)

	logger.Info("match service started", "address", cfg.ServerAddress)
	if err := matchServer.Start(ctx); err != nil {
		return fmt.Errorf("run match service: %w", err)
	}

	logger.Info("match service stopped gracefully")
	return nil
}
