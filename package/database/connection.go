package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/models"
)

func Connect(user string, pass string, host string, port string, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, pass, name)
	if port != "" {
		dsn = fmt.Sprintf("%s port=%s", dsn, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	// TODO: do not migrate the schema in production
	db.AutoMigrate(
		&models.Lap{},
		&models.Event{},
		&models.EventSession{},
	)

	return db, nil
}
