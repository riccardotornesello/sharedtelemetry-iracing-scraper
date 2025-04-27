package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/config"
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

func (c *Client) CheckIfSessionExists(ctx context.Context, subsessionId int) (bool, error) {
	// Check if a document with the given subsessionId exists
	id := fmt.Sprintf("%d", subsessionId)
	docRef := c.Collection(SessionsCollection).Doc(id)
	_, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil // Document does not exist
		}

		return false, err // Some other error occurred
	}
	return true, nil // Document exists
}

func (c *Client) AddSession(ctx context.Context, subsessionId int, sessionData *Session) error {
	// Add a new document with the given subsessionId
	id := fmt.Sprintf("%d", subsessionId)
	docRef := c.Collection(SessionsCollection).Doc(id)
	_, err := docRef.Set(ctx, sessionData)
	if err != nil {
		return err
	}
	return nil
}
