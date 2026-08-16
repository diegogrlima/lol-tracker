package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegogrlima/lol-tracker/internal/application"
	"github.com/diegogrlima/lol-tracker/internal/config"
	"github.com/diegogrlima/lol-tracker/internal/database"
	"github.com/diegogrlima/lol-tracker/internal/riot"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application stopped with error: %v", err)
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

	redisClient, err := database.NewRedis(ctx, cfg.RedisAddress)
	if err != nil {
		return fmt.Errorf("initialize Redis: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("close Redis connection: %v", err)
		}
	}()

	riotClient, err := riot.NewClient(cfg.RiotAPIKey, cfg.RiotRegion)
	if err != nil {
		return fmt.Errorf("initialize Riot client: %w", err)
	}

	app := application.New(
		riotClient,
		cfg.ServerAddress,
	)

	log.Printf("API started on %s", cfg.ServerAddress)

	if err := app.Start(ctx); err != nil {
		return fmt.Errorf("run application: %w", err)
	}

	log.Println("API stopped gracefully")

	return nil
}
