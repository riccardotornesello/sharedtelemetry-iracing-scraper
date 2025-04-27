package iracing

import (
	"log"

	irapi "github.com/riccardotornesello/irapi-go"
	"github.com/riccardotornesello/irapi-go/api/car"
	"github.com/riccardotornesello/irapi-go/api/carclass"
	"riccardotornesello.it/sharedtelemetry/iracing/cars/config"
)

type IRacingData struct {
	Cars       *car.CarGetResponse
	CarAssets  *car.CarAssetsResponse
	CarClasses *carclass.CarclassGetResponse
}

func FetchData(cfg config.Config) (*IRacingData, error) {
	// Authenticate
	log.Println("Authenticating with iRacing API")
	irClient, err := irapi.NewIRacingApiClient(cfg.IRacingEmail, cfg.IRacingPassword)
	if err != nil {
		return nil, err
	}

	// Get the data
	log.Println("Fetching cars")
	cars, err := irClient.Car.GetCar()
	if err != nil {
		return nil, err
	}

	log.Println("Fetching car assets")
	carAssets, err := irClient.Car.GetCarAssets()
	if err != nil {
		return nil, err
	}

	log.Println("Fetching car classes")
	carClasses, err := irClient.Carclass.GetCarclass()
	if err != nil {
		return nil, err
	}

	return &IRacingData{
		Cars:       cars,
		CarAssets:  carAssets,
		CarClasses: carClasses,
	}, nil
}
