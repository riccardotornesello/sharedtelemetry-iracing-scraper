package logic

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/firestore"
	firestore_structs "riccardotornesello.it/sharedtelemetry/iracing/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func UpdateDriverStatsByCategory(firestoreClient *firestore.Client, firestoreContext context.Context, irClient *irapi.IRacingApiClient, carClass string) error {
	// Get the stats CSV
	log.Println("Fetching drivers stats for car class", carClass)
	driversCsv, err := NewDriversCsv(irClient, carClass)
	if err != nil {
		return err
	}
	log.Println("Drivers stats fetched")

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
		driver := firestore_structs.Driver{
			Name:     record.Driver,
			Location: record.Location,
			Stats:    firestore_structs.DriverStats{},
		}
		driverStats := &firestore_structs.DriverStatsDetails{
			License: record.Class,
			IRating: record.Irating,
		}

		switch carClass {
		case "dirt_oval":
			driver.Stats.DirtOval = driverStats
		case "dirt_road":
			driver.Stats.DirtRoad = driverStats
		case "formula_car":
			driver.Stats.FormulaCar = driverStats
		case "oval":
			driver.Stats.Oval = driverStats
		case "road":
			driver.Stats.Road = driverStats
		case "sports_car":
			driver.Stats.SportsCar = driverStats
		default:
			return fmt.Errorf("Invalid car class: %v", carClass)
		}

		// Store in Firestore
		_, err = firestoreClient.Collection(firestore_structs.DriversCollection).Doc(driverId).Set(firestoreContext, driver, firestore.MergeAll)
		if err != nil {
			return err
		}
	}

	return nil
}
