package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"riccardotornesello.it/iracing-average-lap/models"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=example dbname=postgres port=5432"), &gorm.Config{}) // TODO: make variable
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.Lap{},
		&models.Event{},
		&models.EventSession{},
	)

	return db, nil
}
