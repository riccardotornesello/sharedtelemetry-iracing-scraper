package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/database"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

func BenchmarkDriversParse(b *testing.B) {
	err := godotenv.Load("../../.env")
	if err != nil {
		b.Fatalf("Error loading .env file: %v", err)
	}

	_, err = database.Connect("postgres", "example", "127.0.0.1", "5432", "postgres")
	if err != nil {
		b.Fatalf("Connect(...) error: %v", err)
	}

	irClient, err := irapi.NewIRacingApiClient(os.Getenv("IRACING_EMAIL"), os.Getenv("IRACING_PASSWORD"))
	if err != nil {
		b.Fatalf("NewIRacingApiClient(...) error: %v", err)
	}

	csvContent, err := irClient.GetDriverStatsByCategorySportsCar()
	if err != nil {
		b.Fatalf("GetDriverStatsByCategorySportsCar() error: %v", err)
	}

	csvReader := csv.NewReader(csvContent)

	// Check the header
	header, err := csvReader.Read()
	if err != nil {
		b.Fatalf("Error reading CSV header: %v", err)
	}

	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" {
		b.Fatalf("Invalid CSV header: %v", header)
	}

	// Count the rows and how many duplicates
	drivers := make(map[int]struct{})
	duplicates := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}

		if len(drivers) < 10 {
			b.Logf("Record: %v", record[0])
		}

		custId, err := strconv.Atoi(record[1])
		if err != nil {
			b.Fatalf("Error converting custId to int: %v", err)
		}

		if _, ok := drivers[custId]; ok {
			duplicates++
		} else {
			drivers[custId] = struct{}{}
		}
	}

	b.Logf("Drivers: %d, Duplicates: %d", len(drivers), duplicates)

	// defer respBody.Close()
	// body, err := io.ReadAll(respBody)
	// if err != nil {
	// 	b.Fatalf("Error reading response body: %v", err)
	// }

	// b.Logf("CSV size: %d", len(body))

	// csvReader := csv.NewReader(strings.NewReader(string(csvContent)))

	// // Check the header
	// header, err := csvReader.Read()
	// if err != nil {
	// 	b.Fatalf("Error reading CSV header: %v", err)
	// }

	// if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" {
	// 	b.Fatalf("Invalid CSV header: %v", header)
	// }

	// drivers := make([]*models.Driver, 0)
	// for {
	// 	record, err := csvReader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		b.Fatalf("Error reading CSV record: %v", err)
	// 		return
	// 	}

	// 	custId, err := strconv.Atoi(record[1])
	// 	if err != nil {
	// 		b.Fatalf("Error converting custId to int: %v", err)
	// 		return
	// 	}

	// 	drivers = append(drivers, &models.Driver{
	// 		Name:     record[0],
	// 		CustID:   custId,
	// 		Location: record[2],
	// 	})
	// }
}
