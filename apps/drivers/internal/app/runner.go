package app

import (
	"context"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/drivers/config"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/internal/iracing"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/internal/processor"
)

func Run(ctx context.Context, carClass string) error {
	cfg := config.Load()

	fsClient, err := firestore.NewClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer fsClient.Close()

	data, err := iracing.FetchData(cfg, carClass)
	if err != nil {
		return err
	}

	processed, err := processor.Process(data)
	if err != nil {
		return err
	}

	if err := fsClient.SaveDrivers(ctx, processed, carClass); err != nil {
		return err
	}

	log.Println("Done")
	return nil
}
