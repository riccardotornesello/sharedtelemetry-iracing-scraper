package leagues

import (
	"context"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/internal/app"
)

func init() {
	functions.CloudEvent("Handler", handler)
}

func handler(ctx context.Context, e event.Event) error {
	err := app.Run(ctx)
	return err
}
