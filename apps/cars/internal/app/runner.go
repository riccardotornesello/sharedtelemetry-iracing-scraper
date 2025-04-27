package app

import (
	"context"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/cars/config"
	"riccardotornesello.it/sharedtelemetry/iracing/cars/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/cars/internal/iracing"
	"riccardotornesello.it/sharedtelemetry/iracing/cars/internal/processor"
)

func Run(ctx context.Context) error {
	cfg := config.Load()

	fsClient, err := firestore.NewClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer fsClient.Close()

	data, err := iracing.FetchData(cfg)
	if err != nil {
		return err
	}

	processed, err := processor.Process(data)
	if err != nil {
		return err
	}

	if err := fsClient.SaveCars(ctx, processed.Cars); err != nil {
		return err
	}

	if err := fsClient.SaveCarClasses(ctx, processed.CarClasses); err != nil {
		return err
	}

	log.Println("Done")
	return nil
}
