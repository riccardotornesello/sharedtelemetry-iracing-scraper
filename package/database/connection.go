package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"riccardotornesello.it/iracing-average-lap/models"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("db.sql"), &gorm.Config{}) // TODO: make variable
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
