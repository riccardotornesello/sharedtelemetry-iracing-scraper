package firestore

import (
	"context"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/cars/config"

	"cloud.google.com/go/firestore"
)

const (
	CarsCollection       = "iracing_cars"
	CarClassesCollection = "iracing_car_classes"
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

func (c *Client) SaveCars(ctx context.Context, cars map[string]Car) error {
	log.Println("Saving cars to Firestore")

	db := c.Collection(CarsCollection)

	for carId, car := range cars {
		_, err := db.Doc(carId).Set(ctx, car)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) SaveCarClasses(ctx context.Context, carClasses map[string]CarClass) error {
	log.Println("Saving car classes to Firestore")

	db := c.Collection(CarClassesCollection)

	for carClassId, carClass := range carClasses {
		_, err := db.Doc(carClassId).Set(ctx, carClass)
		if err != nil {
			return err
		}
	}
	return nil
}
