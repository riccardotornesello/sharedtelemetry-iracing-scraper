package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	"riccardotornesello.it/sharedtelemetry/iracing/exporter/logic"
)

func main() {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	secretKey := os.Getenv("SECRET_KEY")

	r := gin.Default()
	r.POST("/export", func(c *gin.Context) {
		if secretKey != c.GetHeader("SECRET_KEY") {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// Initialize database
		db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, nil, 1, 1)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		// Get the data
		data, err := logic.GenerateSessionsCsv(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		// Return the data as .csv to be downloaded
		c.Data(200, "text/csv", []byte(data))
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
