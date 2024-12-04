package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"riccardotornesello.it/iracing-average-lap/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	irClient := client.NewIRacingApiClient(os.Getenv("IRACING_EMAIL"), os.Getenv("IRACING_PASSWORD"))

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
