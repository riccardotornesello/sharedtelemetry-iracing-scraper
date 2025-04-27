package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/internal/app"
)

func main() {
	// Get configuration from environment variables
	godotenv.Load()

	// Get the car class from the command line arguments
	if len(os.Args) < 2 {
		log.Fatal("You must provide a car class as an argument")
	}
	carClass := os.Args[1]

	// Run the application
	if err := app.Run(context.Background(), carClass); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
