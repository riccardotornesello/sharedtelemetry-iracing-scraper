package logic

import (
	"context"
	"os"
	"testing"
	"time"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func TestParseSession(t *testing.T) {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	// Initialize database
	firestoreContext := context.Background()
	firebaseConf := &firebase.Config{ProjectID: "sharedtelemetryapp"}
	firebaseApp, err := firebase.NewApp(firestoreContext, firebaseConf)
	if err != nil {
		t.Fatal(err)
	}

	firestoreClient, err := firebaseApp.Firestore(firestoreContext)
	if err != nil {
		t.Fatal(err)
	}
	defer firestoreClient.Close()

	// Initialize iRacing client
	irClient, err := irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		t.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}

	// Get the sessions
	sessions, err := irClient.GetLeagueSeasonSessions(4403, 0, true)
	if err != nil {
		t.Fatal(err)
	}

	// Parse the sessions
	for i := len(sessions.Sessions) - 1; i >= 0; i-- {
		t.Log("Parsing session", sessions.Sessions[i].SubsessionId)

		launchAt, err := time.Parse(time.RFC3339, sessions.Sessions[i].LaunchAt)
		if err != nil {
			t.Fatal(err)
		}

		err = ParseSession(irClient, sessions.Sessions[i].SubsessionId, launchAt, firestoreClient, firestoreContext, 10)
		if err != nil {
			t.Fatal(err)
		}
	}
}
