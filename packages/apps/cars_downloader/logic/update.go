package logic

import (
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/cars_models"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func UpdateCarsDb(db *gorm.DB, irClient *irapi.IRacingApiClient) error {
	now := time.Now()

	// Get the data
	log.Println("Fetching cars")
	cars, err := irClient.GetCars()
	if err != nil {
		return err
	}

	log.Println("Fetching car assets")
	carAssets, err := irClient.GetCarAssets()
	if err != nil {
		return err
	}

	log.Println("Fetching car classes")
	carClasses, err := irClient.GetCarClasses()
	if err != nil {
		return err
	}

	log.Println("Analyzing data")
	dbCars := make([]cars_models.Car, len(*cars))
	for i, car := range *cars {
		dbCars[i] = cars_models.Car{
			ID:              &car.CarId,
			CreatedAt:       now,
			UpdatedAt:       now,
			Name:            car.CarName,
			NameAbbreviated: car.CarNameAbbreviated,
			Brand:           strings.ToUpper(car.CarMake),
			Logo:            carAssets[car.CarId].Logo,
			SmallImage:      carAssets[car.CarId].SmallImage,
			SponsorLogo:     carAssets[car.CarId].SponsorLogo,
		}
	}

	dbCarClasses := make([]cars_models.CarClass, len(*carClasses))
	dbCarsInClass := make([]cars_models.CarInClass, 0)
	for i, carClass := range *carClasses {
		dbCarClasses[i] = cars_models.CarClass{
			ID:        &carClass.CarClassId,
			CreatedAt: now,
			UpdatedAt: now,
			Name:      carClass.Name,
			ShortName: carClass.ShortName,
		}

		for _, carInClass := range carClass.CarsInClass {
			dbCarsInClass = append(dbCarsInClass, cars_models.CarInClass{
				CarID:      carInClass.CarId,
				CarClassID: carClass.CarClassId,
			})
		}
	}

	log.Println("Saving data")
	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("1=1").Delete(&cars_models.CarInClass{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("1=1").Delete(&cars_models.Car{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("1=1").Delete(&cars_models.CarClass{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Create(&dbCars).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Create(&dbCarClasses).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Create(&dbCarsInClass).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
