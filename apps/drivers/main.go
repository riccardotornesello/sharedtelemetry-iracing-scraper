package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/internal/app"
)

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type MessageData struct {
	CarClass string `json:"carClass"`
}

func init() {
	functions.CloudEvent("Handler", handler)
}

func handler(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	var data MessageData
	if err := json.Unmarshal(msg.Message.Data, &data); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	err := app.Run(ctx, data.CarClass)
	return err
}
