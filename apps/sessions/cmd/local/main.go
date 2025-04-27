package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/internal/app"
)

func main() {
	// Get configuration from environment variables
	godotenv.Load()

	// Get the subsession id and launch date from the command line arguments
	if len(os.Args) < 3 {
		log.Fatal("You must provide the subsession id and launch date as arguments")
	}
	subsessionId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid subsession id: %v", err)
	}
	launchAt := os.Args[2]

	// Run the application
	if err := app.Run(context.Background(), subsessionId, launchAt); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
