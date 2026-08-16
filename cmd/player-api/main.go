package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegogrlima/lol-tracker/internal/config"
	"github.com/diegogrlima/lol-tracker/internal/httpapi"
	redisadapter "github.com/diegogrlima/lol-tracker/internal/platform/redis"
	riotadapter "github.com/diegogrlima/lol-tracker/internal/platform/riot"
	"github.com/diegogrlima/lol-tracker/internal/player"
	"github.com/diegogrlima/lol-tracker/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Printf("player service stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
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
			log.Printf("close Redis connection: %v", err)
		}
	}()

	riotClient, err := riotadapter.NewClient(cfg.RiotAPIKey, cfg.RiotRegion)
	if err != nil {
		return fmt.Errorf("initialize Riot client: %w", err)
	}

	playerCache := redisadapter.NewPlayerCache(redisClient)
	playerService := player.NewService(riotClient, playerCache, cfg.CacheTTL)
	playerHandler := httpapi.NewPlayerHandler(playerService)
	router := httpapi.NewRouter(playerHandler)
	playerServer := server.New(cfg.ServerAddress, router)

	log.Printf("player service started on %s", cfg.ServerAddress)

	if err := playerServer.Start(ctx); err != nil {
		return fmt.Errorf("run player service: %w", err)
	}

	log.Println("player service stopped gracefully")
	return nil
}
