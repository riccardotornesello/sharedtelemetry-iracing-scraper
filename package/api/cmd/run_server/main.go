package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/api/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	models "riccardotornesello.it/sharedtelemetry/iracing/db/events_models"
)

func main() {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	// Initialize database
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, models.AllModels, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/competitions/:id", func(c *gin.Context) {
		// logic.GetEventGroupSessions(db)
		// logic.GetCompetitionDrivers(db, 1)
		_, err := logic.GetLaps(db, [][]int{{1, 1}, {2, 2}})
		if err != nil {
			log.Println(err)
		}

		id := c.Param("id")
		c.JSON(200, gin.H{
			"id": id,
		})
		return
	})
	r.Run()
}
