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
	redisadapter "github.com/diegogrlima/lol-tracker/internal/platform/redis"
	riotadapter "github.com/diegogrlima/lol-tracker/internal/platform/riot"
	"github.com/diegogrlima/lol-tracker/internal/player"
	"github.com/diegogrlima/lol-tracker/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("player service stopped", "error", err)
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

	cfg, err := config.LoadPlayer()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	redisClient, err := redisadapter.NewClient(
		ctx,
		cfg.RedisAddress,
		cfg.RedisPassword,
		cfg.RedisDB,
	)
	if err != nil {
		return fmt.Errorf("initialize Redis: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Warn("failed to close Redis connection", "error", err)
		}
	}()

	riotClient, err := riotadapter.NewClient(cfg.RiotAPIKey, cfg.RiotRegion)
	if err != nil {
		return fmt.Errorf("initialize Riot client: %w", err)
	}

	playerCache := redisadapter.NewPlayerCache(redisClient)
	playerService := player.NewService(riotClient, playerCache, cfg.CacheTTL, logger)
	playerHandler := player.NewHandler(playerService, logger)
	router := server.NewRouter(map[string]http.Handler{
		"/players": playerHandler.Routes(),
	})
	playerServer := server.New(cfg.ServerAddress, router)

	logger.Info("player service started", "address", cfg.ServerAddress)

	if err := playerServer.Run(ctx); err != nil {
		return fmt.Errorf("run player service: %w", err)
	}

	logger.Info("player service stopped gracefully")
	return nil
}
