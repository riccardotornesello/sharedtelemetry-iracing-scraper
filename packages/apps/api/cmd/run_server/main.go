package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/api/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
)

func main() {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	eventsDbUser := os.Getenv("EVENTS_DB_USER")
	eventsDbPass := os.Getenv("EVENTS_DB_PASS")
	eventsDbName := os.Getenv("EVENTS_DB_NAME")
	eventsDbPort := os.Getenv("EVENTS_DB_PORT")
	eventsDbHost := os.Getenv("EVENTS_DB_HOST")

	carsDbUser := os.Getenv("CARS_DB_USER")
	carsDbPass := os.Getenv("CARS_DB_PASS")
	carsDbName := os.Getenv("CARS_DB_NAME")
	carsDbPort := os.Getenv("CARS_DB_PORT")
	carsDbHost := os.Getenv("CARS_DB_HOST")

	// Initialize database
	eventsDb, err := database.Connect(eventsDbUser, eventsDbPass, eventsDbHost, eventsDbPort, eventsDbName, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	carsDb, err := database.Connect(carsDbUser, carsDbPass, carsDbHost, carsDbPort, carsDbName, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// Handlers
	r.GET("/competitions/:id/ranking", func(c *gin.Context) {
		handlers.CompetitionRankingHandler(c, eventsDb, carsDb)
	})

	r.GET("/competitions/:id/csv", func(c *gin.Context) {
		handlers.CompetitionCsvHandler(c, eventsDb)
	})

	r.Run()
}
