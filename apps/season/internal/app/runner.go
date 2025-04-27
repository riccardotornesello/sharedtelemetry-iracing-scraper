package app

import (
	"context"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/season/config"
	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/iracing"
	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/processor"
	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/pubsub"
)

func Run(ctx context.Context, leagueId int, seasonId int) error {
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

	data, err := iracing.FetchData(cfg, leagueId, seasonId)
	if err != nil {
		return err
	}

	lastSessionLaunchAt, err := fsClient.GetLastSessionLaunchAt(ctx, leagueId, seasonId)
	if err != nil {
		return err
	}

	processed, err := processor.Process(data, lastSessionLaunchAt)
	if err != nil {
		return err
	}

	if err := pubSubClient.SendSessions(ctx, processed); err != nil {
		return err
	}

	log.Println("Done")
	return nil
}
