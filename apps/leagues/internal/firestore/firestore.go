package firestore

import (
	"context"
	"log"

	"google.golang.org/api/iterator"
	"riccardotornesello.it/sharedtelemetry/iracing/leagues/config"

	"cloud.google.com/go/firestore"
)

const (
	LeaguesCollection = "iracing_leagues"
)

type Client struct {
	*firestore.Client
}

func NewClient(ctx context.Context, cfg config.Config) (*Client, error) {
	client, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

func (c *Client) GetLeagues(ctx context.Context) ([]League, error) {
	log.Println("Getting leagues from Firestore")

	db := c.Collection(LeaguesCollection)

	var leagues []League
	iter := db.Documents(ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var league League
		doc.DataTo(&league)
		leagues = append(leagues, league)
	}
	return leagues, nil
}
