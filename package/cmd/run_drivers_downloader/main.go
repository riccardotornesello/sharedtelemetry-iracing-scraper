package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
	"riccardotornesello.it/sharedtelemetry/iracing/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/models"
)

var db *gorm.DB
var irClient *irapi.IRacingApiClient

func main() {
	var err error

	db, irClient, err = logic.InitHandler()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", PubSubHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func PubSubHandler(w http.ResponseWriter, r *http.Request) {
	csvContent, err := irClient.GetDriverStatsByCategorySportsCar()
	if err != nil {
		log.Printf("Error getting driver stats by category: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	csvReader := csv.NewReader(csvContent)

	// Check the header
	header, err := csvReader.Read()
	if err != nil {
		log.Printf("Error reading CSV header: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" {
		log.Printf("Invalid CSV header: %v", header)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Insert the users in groups of 100
	isEof := false

	for !isEof {
		drivers := make([]*models.Driver, 100)
		n := 0

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				isEof = true
				break
			}
			if err != nil {
				log.Printf("Error reading CSV record: %v", err)
				w.WriteHeader(http.StatusOK)
				return
			}

			custId, err := strconv.Atoi(record[1])
			if err != nil {
				log.Printf("Error converting custId to int: %v", err)
				w.WriteHeader(http.StatusOK)
				return
			}

			drivers[n] = &models.Driver{
				Name:     record[0],
				CustID:   custId,
				Location: record[2],
			}

			n++

			if n == 100 {
				break
			}
		}

		if n > 0 {
			if err = db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cust_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "location"}),
			}).Create(drivers[:n]).Error; err != nil {
				log.Printf("Error inserting drivers: %v", err)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
