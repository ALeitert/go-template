package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/risingwavelabs/eris"

	"template/internal/services"
)

func main() {
	fmt.Println("Template for Go programs with multiple services.")

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	//
	// Run services.

	err := services.Run(ctx, []services.Service{
		// List services here.
	})
	if err != nil {
		fmt.Println("Error while running services:", eris.ToString(err, true))
		os.Exit(1)
	}
}
