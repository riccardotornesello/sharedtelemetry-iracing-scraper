package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers_downloader/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers_models"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

const BATCH_SIZE = 500

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

	carClass := os.Getenv("CAR_CLASS")

	// Initialize database
	log.Println("Connecting to database")
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, drivers_models.AllModels, 20, 2)
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
	log.Println("Starting job for car class", carClass)
	err = logic.UpdateDriverStatsByCategory(db, irClient, carClass, BATCH_SIZE)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Job completed")
}
