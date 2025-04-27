package iracing

import (
	"fmt"
	"io"
	"log"

	irapi "github.com/riccardotornesello/irapi-go"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/config"
)

func FetchData(cfg config.Config, carClass string) (io.ReadCloser, error) {
	var csvContent io.ReadCloser

	// Authenticate
	log.Println("Authenticating with iRacing API")
	irClient, err := irapi.NewIRacingApiClient(cfg.IRacingEmail, cfg.IRacingPassword)
	if err != nil {
		return nil, err
	}

	// Get the raw csv
	log.Println("Fetching csv")

	switch carClass {
	case "dirt_oval":
		csvContent, err = irClient.DriverStatsByCategory.GetDriverStatsByCategoryDirtOval()
	case "dirt_road":
		csvContent, err = irClient.DriverStatsByCategory.GetDriverStatsByCategoryDirtRoad()
	case "formula_car":
		csvContent, err = irClient.DriverStatsByCategory.GetDriverStatsByCategoryFormulaCar()
	case "oval":
		csvContent, err = irClient.DriverStatsByCategory.GetDriverStatsByCategoryOval()
	case "road":
		csvContent, err = irClient.DriverStatsByCategory.GetDriverStatsByCategoryRoad()
	case "sports_car":
		csvContent, err = irClient.DriverStatsByCategory.GetDriverStatsByCategorySportsCar()
	default:
		err = fmt.Errorf("Invalid car class: %v", carClass)
	}

	return csvContent, err
}
