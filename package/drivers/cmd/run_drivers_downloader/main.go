package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	common "riccardotornesello.it/sharedtelemetry/iracing/common/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers/models"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

// TODO: move logic to dedicated function

const batchSize = 100

var db *gorm.DB
var irClient *irapi.IRacingApiClient

func main() {
	var err error

	// Get configuration
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	// Initialize database
	db, err = database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, models.AllModels, 20, 2)
	if err != nil {
		log.Fatalf("database.Connect: %v", err)
	}

	// Initialize iRacing client
	irClient, err = irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		log.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}

	// Start the HTTP server
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
		common.ReturnException(w, err, "io.ReadAll")
		return
	}

	var m PubSubMessage
	if err := json.Unmarshal(body, &m); err != nil {
		common.ReturnException(w, err, "json.Unmarshal")
		return
	}

	var seasonData SeasonData
	if err := json.Unmarshal(m.Message.Data, &seasonData); err != nil {
		common.ReturnException(w, err, "json.Unmarshal")
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
		common.ReturnException(w, fmt.Errorf("Invalid car class: %v", seasonData.CarClass), "Handler")
		return
	}

	if err != nil {
		common.ReturnException(w, err, "Handler")
		return
	}

	csvReader := csv.NewReader(csvContent)

	// Check the header
	header, err := csvReader.Read()
	if err != nil {
		common.ReturnException(w, err, "csvReader.Read")
		return
	}

	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" || header[13] != "CLASS" || header[14] != "IRATING" {
		common.ReturnException(w, fmt.Errorf("Invalid CSV header: %v", header), "Handler")
		return
	}

	date := database.Date(time.Now().Format("2006-01-02"))

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
				common.ReturnException(w, err, "csvReader.Read")
				return
			}

			custId, err := strconv.Atoi(record[1])
			if err != nil {
				common.ReturnException(w, err, "strconv.Atoi")
				return
			}

			drivers[n] = &models.Driver{
				Name:     record[0],
				CustID:   custId,
				Location: record[2],
			}

			irating, err := strconv.Atoi(record[14])
			if err != nil {
				common.ReturnException(w, err, "strconv.Atoi")
				return
			}

			driverStats[n] = &models.DriverStats{
				CustID:      custId,
				CarCategory: "sports_car",
				Date:        date,
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
				common.ReturnException(w, err, "db.Create")
				return
			}

			// Save stats
			if err = db.Create(driverStats[:n]).Error; err != nil {
				common.ReturnException(w, err, "db.Create")
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
