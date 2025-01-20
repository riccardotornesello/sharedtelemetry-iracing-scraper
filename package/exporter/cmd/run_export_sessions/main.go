package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
)

var db *gorm.DB

func main() {
	var err error

	// Get configuration
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	// Initialize database
	db, err = database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, nil, 20, 2)
	if err != nil {
		log.Fatalf("database.Connect: %v", err)
	}

	// Start the HTTP server
	http.HandleFunc("/", PubSubHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func PubSubHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get filters from the message

	w.WriteHeader(http.StatusOK)
	return
}
