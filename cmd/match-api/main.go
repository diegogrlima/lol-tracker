package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegogrlima/lol-tracker/internal/config"
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

	address, err := config.LoadMatchServerAddress()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	matchServer := server.New(address, server.NewRouter(nil))

	logger.Info("match service scaffold started", "address", address)
	if err := matchServer.Start(ctx); err != nil {
		return fmt.Errorf("run match service: %w", err)
	}

	logger.Info("match service stopped gracefully")
	return nil
}
