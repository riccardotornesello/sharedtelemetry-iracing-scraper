package main

import (
	"encoding/csv"
	"encoding/json"
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

const batchSize = 100

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

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type SeasonData struct {
	CarClass string `json:"carClass"`
}

func PubSubHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("io.ReadAll: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	var m PubSubMessage
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	var seasonData SeasonData
	if err := json.Unmarshal(m.Message.Data, &seasonData); err != nil {
		log.Printf("json.Unmarshal data: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	var csvContent io.ReadCloser
	switch seasonData.CarClass {
	case "sports_car":
		csvContent, err = irClient.GetDriverStatsByCategorySportsCar()
	case "oval":
		csvContent, err = irClient.GetDriverStatsByCategoryOval()
	case "formula_car":
		csvContent, err = irClient.GetDriverStatsByCategoryFormulaCar()
	case "road":
		csvContent, err = irClient.GetDriverStatsByCategoryRoad()
	case "dirt_oval":
		csvContent, err = irClient.GetDriverStatsByCategoryDirtOval()
	case "dirt_road":
		csvContent, err = irClient.GetDriverStatsByCategoryDirtRoad()
	default:
		log.Printf("Invalid car class: %v", seasonData.CarClass)
		w.WriteHeader(http.StatusOK)
		return
	}

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

	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" || header[13] != "CLASS" || header[14] != "IRATING" {
		log.Printf("Invalid CSV header: %v", header)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Insert the users in groups of batchSize
	isEof := false

	for !isEof {
		drivers := make([]*models.Driver, batchSize)
		driverStats := make([]*models.DriverStats, batchSize)
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

			irating, err := strconv.Atoi(record[14])
			if err != nil {
				log.Printf("Error converting iRating to int: %v", err)
				w.WriteHeader(http.StatusOK)
				return
			}

			driverStats[n] = &models.DriverStats{
				CustID:      custId,
				CarCategory: "sports_car",
				License:     record[13],
				IRating:     irating,
			}

			n++

			if n == batchSize {
				break
			}
		}

		if n > 0 {
			// Update drivers list
			if err = db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cust_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "location"}),
			}).Create(drivers[:n]).Error; err != nil {
				log.Printf("Error inserting drivers: %v", err)
				w.WriteHeader(http.StatusOK)
				return
			}

			// Save stats
			if err = db.Create(driverStats[:n]).Error; err != nil {
				log.Printf("Error inserting driver stats: %v", err)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
