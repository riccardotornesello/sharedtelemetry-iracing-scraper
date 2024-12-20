package logic

import (
	"os"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

func InitCloudRun(models []interface{}) (*gorm.DB, *irapi.IRacingApiClient, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, models)
	if err != nil {
		return nil, nil, err
	}

	irClient, err := irapi.NewIRacingApiClient(os.Getenv("IRACING_EMAIL"), os.Getenv("IRACING_PASSWORD"))
	if err != nil {
		return nil, nil, err
	}

	return db, irClient, nil
}
