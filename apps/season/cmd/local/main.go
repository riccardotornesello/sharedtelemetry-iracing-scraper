package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/app"
)

func main() {
	// Get configuration from environment variables
	godotenv.Load()

	// Get the league id and season id from the command line arguments
	if len(os.Args) < 3 {
		log.Fatal("You must provide the league id and season id as arguments")
	}
	leagueId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid league id: %v", err)
	}
	seasonId, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid season id: %v", err)
	}

	// Run the application
	if err := app.Run(context.Background(), leagueId, seasonId); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
