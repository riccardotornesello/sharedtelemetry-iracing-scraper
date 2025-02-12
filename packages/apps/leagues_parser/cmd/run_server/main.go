package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues_parser/logic"
)

func main() {
	// Get configuration
	godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	pubSubProjectId := os.Getenv("PUBSUB_PROJECT")
	pubSubTopicId := os.Getenv("PUBSUB_TOPIC")

	// Initialize database
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, 1, 0)
	if err != nil {
		log.Fatalf("database.Connect: %v", err)
		return
	}

	// Initialize pubsub
	pubSubCtx := context.Background()
	client, err := pubsub.NewClient(pubSubCtx, pubSubProjectId)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	pubSubTopic := client.Topic(pubSubTopicId)

	// Start the job
	log.Println("Starting job")

	err = logic.FeedLeagues(db, pubSubTopic, pubSubCtx)
	if err != nil {
		log.Fatalf("logic.FeedLeagues: %v", err)
		return
	}

	log.Println("Job completed")
}
