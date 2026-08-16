package main

import (
	"context"
	"fmt"

	"github.com/diegogrlima/lol-tracker/internal/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
