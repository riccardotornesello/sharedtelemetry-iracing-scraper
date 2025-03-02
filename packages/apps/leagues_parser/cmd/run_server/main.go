package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues_parser/logic"
)

const projectID = "sharedtelemetryapp" // TODO: move to env

func main() {
	// Get configuration
	godotenv.Load()

	pubSubProjectId := os.Getenv("PUBSUB_PROJECT")
	pubSubTopicId := os.Getenv("PUBSUB_TOPIC")

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

	// Get active leagues
	leagues, err := logic.GetActiveLeagues(firestoreClient, firestoreContext)
	if err != nil {
		log.Fatalf("Failed to get active leagues: %v", err)
		return
	}

	// Send pub/sub messages to parse the season
	// TODO: move logic to specific function
	var wg sync.WaitGroup
	var totalErrors uint64

	for _, league := range leagues {
		for _, season := range league.Seasons {
			result := pubSubTopic.Publish(pubSubCtx, &pubsub.Message{
				Data: []byte("{\"leagueId\":" + strconv.Itoa(league.LeagueID) + ",\"seasonId\":" + strconv.Itoa(season.SeasonID) + "}"),
			})

			wg.Add(1)
			go func(res *pubsub.PublishResult) {
				defer wg.Done()
				_, err := res.Get(pubSubCtx)
				if err != nil {
					atomic.AddUint64(&totalErrors, 1)
					return
				}
			}(result)
		}
	}

	wg.Wait()

	if totalErrors > 0 {
		log.Fatalf("Failed to send %d pub/sub messages", totalErrors)
		return
	}

	log.Println("Job completed")
}
