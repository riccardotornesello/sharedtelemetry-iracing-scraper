package logic

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

type DriversCsv struct {
	csvReader *csv.Reader
}

type DriversCsvRow struct {
	Driver   string
	CustId   int
	Location string
	Class    string
	Irating  int
}

func GetDriverStatsByCategory(irClient *irapi.IRacingApiClient, carClass string) (io.ReadCloser, error) {
	var csvContent io.ReadCloser
	var err error

	switch carClass {
	case "dirt_oval":
		csvContent, err = irClient.GetDriverStatsByCategoryDirtOval()
	case "dirt_road":
		csvContent, err = irClient.GetDriverStatsByCategoryDirtRoad()
	case "formula_car":
		csvContent, err = irClient.GetDriverStatsByCategoryFormulaCar()
	case "oval":
		csvContent, err = irClient.GetDriverStatsByCategoryOval()
	case "road":
		csvContent, err = irClient.GetDriverStatsByCategoryRoad()
	case "sports_car":
		csvContent, err = irClient.GetDriverStatsByCategorySportsCar()
	default:
		err = fmt.Errorf("Invalid car class: %v", carClass)
	}

	return csvContent, err
}

func CheckCsvReader(header []string) error {
	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" || header[13] != "CLASS" || header[14] != "IRATING" {
		return fmt.Errorf("Invalid header: %v", header)
	}

	return nil
}

func NewDriversCsv(irClient *irapi.IRacingApiClient, carClass string) (*DriversCsv, error) {
	// Get the stats CSV
	csvContent, err := GetDriverStatsByCategory(irClient, carClass)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(csvContent)

	// Check the header
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	err = CheckCsvReader(header)
	if err != nil {
		return nil, err
	}

	// Return the DriversCsv struct
	return &DriversCsv{csvReader: csvReader}, nil
}

func (d *DriversCsv) Read() (*DriversCsvRow, error) {
	record, err := d.csvReader.Read()
	if err != nil {
		return nil, err
	}

	custId, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}

	irating, err := strconv.Atoi(record[14])
	if err != nil {
		return nil, err
	}

	return &DriversCsvRow{
		Driver:   record[0],
		CustId:   custId,
		Location: record[2],
		Class:    record[13],
		Irating:  irating,
	}, nil
}
