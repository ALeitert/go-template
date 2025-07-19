package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/risingwavelabs/eris"

	"template/internal/config"
	"template/internal/database"
	"template/internal/services"
)

func main() {
	fmt.Println("Template for Go programs with multiple services.")

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	//
	// Load and print config.

	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to a file with configurations.")
	flag.Parse()

	err := config.C.Load(configPath)
	if err != nil {
		fmt.Println("Error while loading config:", eris.ToString(err, true))
		os.Exit(1)
	}
	config.C.Print()

	//
	// Connect to database.

	err = database.Connect(ctx)
	if err != nil {
		fmt.Println("Failed to connect to database:", eris.ToString(err, true))
		os.Exit(1)
	}

	//
	// Run services.

	err = services.Run(ctx, []services.Service{
		// List services here.
	})
	if err != nil {
		fmt.Println("Error while running services:", eris.ToString(err, true))
		os.Exit(1)
	}
}
