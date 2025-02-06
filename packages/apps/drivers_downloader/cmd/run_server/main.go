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

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/drivers_models"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

// TODO: move logic to dedicated function

const batchSize = 100

var db *gorm.DB
var irClient *irapi.IRacingApiClient

func main() {
	var err error

	// Get configuration
	godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	// Initialize database
	db, err = database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, drivers_models.AllModels, 20, 2)
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
		handlers.ReturnException(w, err, "io.ReadAll")
		return
	}

	var m PubSubMessage
	if err := json.Unmarshal(body, &m); err != nil {
		handlers.ReturnException(w, err, "json.Unmarshal")
		return
	}

	var seasonData SeasonData
	if err := json.Unmarshal(m.Message.Data, &seasonData); err != nil {
		handlers.ReturnException(w, err, "json.Unmarshal")
		return
	}

	// First check if the stats for the chosen car class have been already downloaded in the current day
	var count int64
	if err := db.Model(&drivers_models.DriverStats{}).Where("car_category = ? AND created_at >= ?", seasonData.CarClass, time.Now().Format("2006-01-02")).Count(&count).Error; err != nil {
		handlers.ReturnException(w, err, "db.Model")
		return
	}

	if count > 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get the stats CSV
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
		handlers.ReturnException(w, fmt.Errorf("Invalid car class: %v", seasonData.CarClass), "Handler")
		return
	}

	if err != nil {
		handlers.ReturnException(w, err, "Handler")
		return
	}

	csvReader := csv.NewReader(csvContent)

	// Check the header
	header, err := csvReader.Read()
	if err != nil {
		handlers.ReturnException(w, err, "csvReader.Read")
		return
	}

	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" || header[13] != "CLASS" || header[14] != "IRATING" {
		handlers.ReturnException(w, fmt.Errorf("Invalid CSV header: %v", header), "Handler")
		return
	}

	now := time.Now()

	// Insert the users in groups of batchSize
	isEof := false

	for !isEof {
		drivers := make([]*drivers_models.Driver, batchSize)
		driverStats := make([]*drivers_models.DriverStats, batchSize)
		n := 0

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				isEof = true
				break
			}
			if err != nil {
				handlers.ReturnException(w, err, "csvReader.Read")
				return
			}

			custId, err := strconv.Atoi(record[1])
			if err != nil {
				handlers.ReturnException(w, err, "strconv.Atoi")
				return
			}

			drivers[n] = &drivers_models.Driver{
				Name:     record[0],
				CustID:   custId,
				Location: record[2],
			}

			irating, err := strconv.Atoi(record[14])
			if err != nil {
				handlers.ReturnException(w, err, "strconv.Atoi")
				return
			}

			driverStats[n] = &drivers_models.DriverStats{
				CustID:      custId,
				CarCategory: "sports_car",
				License:     record[13],
				IRating:     irating,
				CreatedAt:   now,
				UpdatedAt:   now,
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
				handlers.ReturnException(w, err, "db.Create")
				return
			}

			// Update stats
			if err = db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cust_id"}, {Name: "car_category"}},
				DoUpdates: clause.AssignmentColumns([]string{"license", "i_rating", "updated_at"}),
			}).Create(driverStats[:n]).Error; err != nil {
				handlers.ReturnException(w, err, "db.Create")
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
