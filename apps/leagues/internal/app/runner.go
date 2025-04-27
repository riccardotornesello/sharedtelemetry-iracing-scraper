package app

import (
	"context"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/leagues/config"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/internal/pubsub"
)

func Run(ctx context.Context) error {
	cfg := config.Load()

	fsClient, err := firestore.NewClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer fsClient.Close()

	pubSubClient, err := pubsub.NewPublisher(ctx, cfg)
	if err != nil {
		return err
	}
	defer pubSubClient.Close()

	leagues, err := fsClient.GetLeagues(ctx)
	if err != nil {
		return err
	}

	err = pubSubClient.PublishLeagues(ctx, leagues)
	if err != nil {
		return err
	}

	log.Println("Done")
	return nil
}
