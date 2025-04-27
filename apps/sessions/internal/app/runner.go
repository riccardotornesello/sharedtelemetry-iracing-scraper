package app

import (
	"context"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/sessions/config"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/internal/iracing"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/internal/processor"
)

func Run(ctx context.Context, subsessionId int, launchAt string) error {
	cfg := config.Load()

	fsClient, err := firestore.NewClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer fsClient.Close()

	iracingClient, err := iracing.NewClient(cfg)
	if err != nil {
		return err
	}

	alreadyExists, err := fsClient.CheckIfSessionExists(ctx, subsessionId)
	if err != nil {
		return err
	}
	if alreadyExists {
		log.Printf("Session with subsessionId %d already exists in Firestore. Exiting.", subsessionId)
		return nil
	}

	sessionResponse, err := iracingClient.FetchResults(subsessionId)
	if err != nil {
		return err
	}

	session, err := processor.Process(iracingClient, sessionResponse, launchAt)
	if err != nil {
		return err
	}

	err = fsClient.AddSession(ctx, subsessionId, session)
	if err != nil {
		return err
	}

	return nil
}
