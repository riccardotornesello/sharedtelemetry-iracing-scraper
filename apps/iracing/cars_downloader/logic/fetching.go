package logic

import (
	"fmt"
	"log"
	"strings"

	firestore_structs "riccardotornesello.it/sharedtelemetry/iracing/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func FetchCars(irClient *irapi.IRacingApiClient) (map[string]firestore_structs.Car, map[string]firestore_structs.CarClass, error) {
	// Get the data
	log.Println("Fetching cars")
	cars, err := irClient.GetCars()
	if err != nil {
		return nil, nil, err
	}

	log.Println("Fetching car assets")
	carAssets, err := irClient.GetCarAssets()
	if err != nil {
		return nil, nil, err
	}

	log.Println("Fetching car classes")
	carClasses, err := irClient.GetCarClasses()
	if err != nil {
		return nil, nil, err
	}

	log.Println("Analyzing data")
	dbCars := make(map[string]firestore_structs.Car)
	for _, car := range *cars {
		dbCars[fmt.Sprintf("%d", car.CarId)] = firestore_structs.Car{
			Name:            car.CarName,
			NameAbbreviated: car.CarNameAbbreviated,
			Brand:           strings.ToUpper(car.CarMake),
			Logo:            carAssets[car.CarId].Logo,
			SmallImage:      carAssets[car.CarId].SmallImage,
			SponsorLogo:     carAssets[car.CarId].SponsorLogo,
		}
	}

	dbCarClasses := make(map[string]firestore_structs.CarClass)
	for _, carClass := range *carClasses {
		dbCarClass := firestore_structs.CarClass{
			Name:      carClass.Name,
			ShortName: carClass.ShortName,

			Cars: make([]string, 0),
		}

		for _, carInClass := range carClass.CarsInClass {
			dbCarClass.Cars = append(dbCarClass.Cars, fmt.Sprintf("%d", carInClass.CarId))
		}

		dbCarClasses[fmt.Sprintf("%d", carClass.CarClassId)] = dbCarClass
	}

	return dbCars, dbCarClasses, nil
}
