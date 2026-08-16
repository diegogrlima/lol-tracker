package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegogrlima/lol-tracker/internal/application"
	"github.com/diegogrlima/lol-tracker/internal/database"
	"github.com/diegogrlima/lol-tracker/internal/riot"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	redisAddress := os.Getenv("REDIS_ADDRESS")

	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}

	redisClient, err := database.NewRedis(ctx, redisAddress)
	if err != nil {
		log.Fatalf("connect to Redis: %v", err)
	}
	defer redisClient.Close()

	riotAPIKey := os.Getenv("RIOT_API_KEY")
	riotRegion := os.Getenv("RIOT_REGION")
	if riotRegion == "" {
		riotRegion = "americas"
	}

	riotClient, err := riot.NewClient(riotAPIKey, riotRegion)
	if err != nil {
		log.Fatalf("configure Riot client: %v", err)
	}

	app := application.New(redisClient, riotClient)

	log.Println("API started on port 8080")

	if err := app.Start(ctx); err != nil {
		log.Fatalf("run application: %v", err)
	}

	log.Println("API stopped")
}
