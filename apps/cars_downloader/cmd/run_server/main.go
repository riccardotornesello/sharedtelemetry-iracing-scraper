package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/cars_downloader/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func main() {
	// Get configuration
	godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	// Initialize database
	log.Println("Connecting to database")
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, 20, 2)
	if err != nil {
		log.Fatalf("database.Connect: %v", err)
	}
	log.Println("Connected to database")

	// Initialize iRacing client
	log.Println("Initializing iRacing client")
	irClient, err := irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		log.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}
	log.Println("iRacing client initialized")

	// Start the job
	log.Println("Starting job")
	err = logic.UpdateCarsDb(db, irClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Job completed")
}
