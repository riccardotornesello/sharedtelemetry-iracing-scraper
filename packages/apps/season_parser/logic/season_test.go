package logic

import (
	"context"
	"log"
	"os"
	"testing"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func TestGetLeagueSeasonSessionsInfo(t *testing.T) {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	const projectID = "sharedtelemetryapp" // TODO: move to env

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

	// Initialize iRacing client
	irClient, err := irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		log.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}

	// Test
	firstLaunchAt, err := GetFirstSessionLaunchAt(4403, 0, firestoreClient, firestoreContext)
	if err != nil {
		t.Fatal(err)
	}

	seasonSessionsInfo, err := GetLeagueSeasonSessionsInfo(4403, 0, firstLaunchAt, irClient)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Season sessions info: %d", len(seasonSessionsInfo))
}
