package main

import (
	"context"
	"fmt"
	"os"

	"github.com/diegogrlima/lol-tracker/internal/application"
	"github.com/diegogrlima/lol-tracker/internal/database"
)

func main() {
	ctx := context.Background()
	redisAddress := os.Getenv("REDIS_ADDRESS")

	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}

	redisClient, err := database.NewRedis(ctx, redisAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer redisClient.Close()
	fmt.Println("Conectado ao Redis")

	app := application.New(redisClient)

	if err := app.Start(ctx); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
