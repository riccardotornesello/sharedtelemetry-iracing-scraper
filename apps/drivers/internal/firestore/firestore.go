package firestore

import (
	"context"
	"fmt"
	"io"
	"log"

	"riccardotornesello.it/sharedtelemetry/iracing/drivers/config"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/internal/processor"

	"cloud.google.com/go/firestore"
)

const (
	DriversCollection = "iracing_drivers"
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

func (c *Client) SaveDrivers(ctx context.Context, driversCsv *processor.DriversCsv, carClass string) error {
	log.Println("Saving drivers to Firestore")

	db := c.Collection(DriversCollection)
	batch := c.BulkWriter(ctx)

	// Insert the drivers
	for {

		record, err := driversCsv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		driverId := fmt.Sprintf("%d", record.CustId)
		driver := Driver{
			Name:     record.Driver,
			Location: record.Location,
			Stats:    make(map[string]DriverStatsDetails),
		}
		driverStats := &DriverStatsDetails{
			License: record.Class,
			IRating: record.Irating,
		}
		driver.Stats[carClass] = *driverStats

		// Store in Firestore
		_, err = batch.Set(db.Doc(driverId), driver.ToMap(), firestore.MergeAll)
		if err != nil {
			return err
		}
	}

	batch.Flush()

	return nil
}
