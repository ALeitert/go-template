package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/risingwavelabs/eris"

	"template/internal/api"
	"template/internal/config"
	"template/internal/database"
	"template/internal/metrics"
	"template/internal/services"
)

func main() {
	fmt.Println("Template for Go programs with multiple services.")

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := run(ctx)
	if err != nil {
		fmt.Println(eris.ToString(err, true))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	//
	// Load and print config.

	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to a file with configurations.")
	flag.Parse()

	err := config.C.Load(configPath)
	if err != nil {
		return eris.Wrap(err, "error while loading config")
	}
	config.C.Print()

	//
	// Connect to database.

	err = database.Connect(ctx)
	if err != nil {
		return eris.Wrap(err, "Failed to connect to database")
	}

	//
	// Run services.

	err = services.Run(ctx, []services.Service{
		// List services here.
		&api.Server{},
		&metrics.Server{},
	})
	if err != nil {
		return eris.Wrap(err, "error while running services")
	}

	return nil
}
