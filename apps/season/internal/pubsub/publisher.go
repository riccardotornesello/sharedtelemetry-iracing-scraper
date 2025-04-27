package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"riccardotornesello.it/sharedtelemetry/iracing/season/config"
	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/processor"
)

type Publisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPublisher(ctx context.Context, cfg config.Config) (*Publisher, error) {
	client, err := pubsub.NewClient(ctx, cfg.PubSubProjectID)
	if err != nil {
		return nil, fmt.Errorf("errore creazione client PubSub: %w", err)
	}

	topic := client.Topic(cfg.PubSubTopicID)
	return &Publisher{client: client, topic: topic}, nil
}

func (p *Publisher) SendSessions(ctx context.Context, processedData []processor.ProcessedData) error {
	var wg sync.WaitGroup
	var totalErrors uint64

	log.Println("Publishing data to Pub/Sub")

	for _, session := range processedData {
		data, err := json.Marshal(session)
		if err != nil {
			return err
		}

		result := p.topic.Publish(ctx, &pubsub.Message{
			Data: data,
		})

		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := res.Get(ctx)
			if err != nil {
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("Failed to send %d pub/sub messages", totalErrors)
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.client.Close()
}
