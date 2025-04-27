package firestore

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/iterator"
	"riccardotornesello.it/sharedtelemetry/iracing/season/config"

	"cloud.google.com/go/firestore"
)

const (
	SessionsCollection = "iracing_sessions"
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

func (c *Client) GetLastSessionLaunchAt(ctx context.Context, leagueId int, seasonId int) (*time.Time, error) {
	log.Println("Getting last session launchAt")

	db := c.Collection(SessionsCollection)

	query := db.
		Where("leagueId", "==", leagueId).
		Where("seasonId", "==", seasonId).
		OrderBy("launchAt", firestore.Desc).
		Limit(1)

	snapshots := query.Documents(ctx)
	defer snapshots.Stop()

	doc, err := snapshots.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	data := doc.Data()
	launchAt, ok := data["launchAt"].(time.Time)
	if !ok {
		return nil, nil
	}

	return &launchAt, nil
}
