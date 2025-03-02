package logic

import (
	"io"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers_models"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func UpdateDriverStatsByCategory(db *gorm.DB, irClient *irapi.IRacingApiClient, carClass string, batchSize int) error {
	now := time.Now()

	// Get the stats CSV
	log.Println("Fetching drivers stats for car class", carClass)
	driversCsv, err := NewDriversCsv(irClient, carClass)
	if err != nil {
		return err
	}
	log.Println("Drivers stats fetched")

	// Insert the users in groups of size batchSize
	isEof := false

	for !isEof {
		log.Println("Processing batch")
		drivers := make([]*drivers_models.Driver, batchSize)
		driverStats := make([]*drivers_models.DriverStats, batchSize)
		n := 0

		for {
			record, err := driversCsv.Read()
			if err == io.EOF {
				isEof = true
				break
			}
			if err != nil {
				return err
			}

			drivers[n] = &drivers_models.Driver{
				Name:     record.Driver,
				CustID:   record.CustId,
				Location: record.Location,
			}

			driverStats[n] = &drivers_models.DriverStats{
				CustID:      record.CustId,
				CarCategory: carClass,
				License:     record.Class,
				IRating:     record.Irating,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			n++

			if n == batchSize {
				break
			}
		}

		if n > 0 {
			// Update drivers list
			if err = db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cust_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "location"}),
			}).Create(drivers[:n]).Error; err != nil {
				return err
			}

			// Update stats
			if err = db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cust_id"}, {Name: "car_category"}},
				DoUpdates: clause.AssignmentColumns([]string{"license", "i_rating", "updated_at"}),
			}).Create(driverStats[:n]).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
