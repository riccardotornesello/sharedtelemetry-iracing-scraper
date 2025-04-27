package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/internal/app"
)

func main() {
	// Get configuration from environment variables
	godotenv.Load()

	if err := app.Run(context.Background()); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
