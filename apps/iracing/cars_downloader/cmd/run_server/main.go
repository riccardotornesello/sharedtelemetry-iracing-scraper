package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/cars_downloader/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

const projectID = "sharedtelemetryapp" // TODO: move to env

func main() {
	// Get configuration
	godotenv.Load()

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	// Initialize database
	log.Println("Connecting to database")
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
	log.Println("Connected to database")

	// Initialize iRacing client
	log.Println("Initializing iRacing client")
	irClient, err := irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		log.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}
	log.Println("iRacing client initialized")

	// Start the job
	log.Println("Starting job")

	cars, carClasses, err := logic.FetchCars(irClient)
	if err != nil {
		log.Fatal(err)
	}

	err = logic.StoreCars(cars, firestoreClient, firestoreContext)
	if err != nil {
		log.Fatal(err)
	}

	err = logic.StoreCarClasses(carClasses, firestoreClient, firestoreContext)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Job completed")
}
