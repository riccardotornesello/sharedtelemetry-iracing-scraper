package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

func main() {
	// TODO: move to logic and create cloud function

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	irClient, err := irapi.NewIRacingApiClient(os.Getenv("IRACING_EMAIL"), os.Getenv("IRACING_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	csvContent := irClient.GetDriverStatsByCategoryFormulaCar()

	// Connect to S3
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Write to bucket
	bucketName := "iracing-driver-stats-csv" // TODO: variable name
	bucket := client.Bucket(bucketName)
	obj := bucket.Object("test.csv") // TODO: variable name
	w := obj.NewWriter(ctx)
	w.Write(csvContent)
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
}
