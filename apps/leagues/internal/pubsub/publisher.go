package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/config"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/internal/firestore"
)

type Publisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

type PubSubMessage struct {
	LeagueID int `json:"leagueId"`
	SeasonID int `json:"seasonId"`
}

func NewPublisher(ctx context.Context, cfg config.Config) (*Publisher, error) {
	client, err := pubsub.NewClient(ctx, cfg.PubSubProjectID)
	if err != nil {
		return nil, fmt.Errorf("errore creazione client PubSub: %w", err)
	}

	topic := client.Topic(cfg.PubSubTopicID)
	return &Publisher{client: client, topic: topic}, nil
}

func (p *Publisher) PublishLeagues(ctx context.Context, leagues []firestore.League) error {
	var wg sync.WaitGroup
	var totalErrors uint64

	log.Println("Publishing leagues to Pub/Sub")

	for _, league := range leagues {
		for _, season := range league.Seasons {
			message := PubSubMessage{
				LeagueID: league.LeagueID,
				SeasonID: season.SeasonID,
			}
			data, err := json.Marshal(message)
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
