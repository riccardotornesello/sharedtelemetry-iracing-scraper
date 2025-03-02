package logic

import (
	"context"
	"log"
	"testing"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
)

func TestGetFirstSessionLaunchAt(t *testing.T) {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

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

	// Test
	firstLaunchAt, err := GetFirstSessionLaunchAt(4403, 0, firestoreClient, firestoreContext)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("First launch at: %s", firstLaunchAt)
}
