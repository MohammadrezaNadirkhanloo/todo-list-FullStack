package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// "github.com/yourorg/goclean/internal/bootstrap"
	// "github.com/yourorg/goclean/internal/config"
)

//

var version = "dev"

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error Running Application: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	configDir := os.Getenv("APP_CONFIG_DIR")
	if configDir == "" {
		configDir = "config"
	}

	cfg, err := config.Load(configDir)
	if err != nil {

		return fmt.Errorf("Error loading configuration (config) : %w", err)
	}

	app, err := bootstrap.New(ctx, cfg, version)
	if err != nil {
		return err
	}
	defer func() {

		if cerr := app.Close(context.Background()); cerr != nil {
			fmt.Fprintf(os.Stderr, "Error closing resources (bootstrap): %v\n", cerr)
		}
	}()

	if err := app.Server.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
