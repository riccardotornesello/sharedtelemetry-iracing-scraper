package handlers

import (
	"context"
	"testing"

	firebase "firebase.google.com/go"
)

func TestGetGroupSessions(t *testing.T) {
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

	err = getGroupSessions(345, "2020-04-27", firestoreClient, firestoreContext)
	if err != nil {
		t.Fatal(err)
	}
}
