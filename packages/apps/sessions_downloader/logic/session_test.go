package logic

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/events_models"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func TestParseSession(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	irClient, err := irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		t.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}

	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, events_models.AllModels, 20, 2)
	if err != nil {
		log.Fatalf("database.Connect: %v", err)
	}

	ParseSession(irClient, 32057183, time.Now(), db, 3)
}
