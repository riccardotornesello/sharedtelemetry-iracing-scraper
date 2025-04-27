package processor

import (
	"fmt"
	"log"
	"strings"

	"riccardotornesello.it/sharedtelemetry/iracing/cars/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/cars/internal/iracing"
)

type ProcessedData struct {
	Cars       map[string]firestore.Car
	CarClasses map[string]firestore.CarClass
}

func Process(data *iracing.IRacingData) (*ProcessedData, error) {
	log.Println("Analyzing cars")
	dbCars := make(map[string]firestore.Car)
	for _, car := range *data.Cars {
		carIdStr := fmt.Sprintf("%d", car.CarId)
		carAsset, ok := (*data.CarAssets)[carIdStr]
		if !ok {
			return nil, fmt.Errorf("car asset not found for car ID %d", car.CarId)
		}

		dbCars[carIdStr] = firestore.Car{
			Name:            car.CarName,
			NameAbbreviated: car.CarNameAbbreviated,
			Brand:           strings.ToUpper(car.CarMake),
			Logo:            carAsset.Logo,
			SmallImage:      carAsset.SmallImage,
			SponsorLogo:     carAsset.SponsorLogo,
		}
	}

	log.Println("Analyzing car classes")
	dbCarClasses := make(map[string]firestore.CarClass)
	for _, carClass := range *data.CarClasses {
		dbCarClass := firestore.CarClass{
			Name:      carClass.Name,
			ShortName: carClass.ShortName,

			Cars: make([]string, 0),
		}

		for _, carInClass := range carClass.CarsInClass {
			dbCarClass.Cars = append(dbCarClass.Cars, fmt.Sprintf("%d", carInClass.CarId))
		}

		dbCarClasses[fmt.Sprintf("%d", carClass.CarClassId)] = dbCarClass
	}

	return &ProcessedData{
		Cars:       dbCars,
		CarClasses: dbCarClasses,
	}, nil
}
