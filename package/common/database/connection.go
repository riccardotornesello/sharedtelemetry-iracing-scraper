package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(user string, pass string, host string, port string, name string, models []interface{}, maxOpenConns int, maxIdleConns int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, pass, name)
	if port != "" {
		dsn = fmt.Sprintf("%s port=%s", dsn, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// Migrate the schema
	// TODO: do not migrate the schema in production
	if models != nil && len(models) > 0 {
		if err = db.AutoMigrate(
			models...,
		); err != nil {
			return nil, err
		}
	}

	return db, nil
}
