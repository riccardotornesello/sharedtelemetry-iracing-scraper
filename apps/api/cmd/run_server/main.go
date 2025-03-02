package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/api/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
)

const projectID = "sharedtelemetryapp" // TODO: move to env

func main() {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	carsDbUser := os.Getenv("CARS_DB_USER")
	carsDbPass := os.Getenv("CARS_DB_PASS")
	carsDbName := os.Getenv("CARS_DB_NAME")
	carsDbPort := os.Getenv("CARS_DB_PORT")
	carsDbHost := os.Getenv("CARS_DB_HOST")

	// Initialize database
	firestoreContext := context.Background()
	firebaseConf := &firebase.Config{ProjectID: projectID}
	firebaseApp, err := firebase.NewApp(firestoreContext, firebaseConf)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := firebaseApp.Firestore(firestoreContext)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()

	carsDb, err := database.Connect(carsDbUser, carsDbPass, carsDbHost, carsDbPort, carsDbName, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// Handlers
	r.GET("/competitions/:id/ranking", func(c *gin.Context) {
		handlers.CompetitionRankingHandler(c, eventsDb, carsDb, firestoreClient, firestoreContext)
	})

	r.GET("/competitions/:id/csv", func(c *gin.Context) {
		handlers.CompetitionCsvHandler(c, eventsDb)
	})

	r.Run()
}
