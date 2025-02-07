package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/events_models"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
	"riccardotornesello.it/sharedtelemetry/iracing/season_parser/logic"
)

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
	db, err = database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, events_models.AllModels, 2, 2)
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

	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	log.Println("Listening on", listener.Addr())
	if err := http.Serve(listener, nil); err != nil {
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
	LeagueId int `json:"leagueId"`
	SeasonId int `json:"seasonId"`
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

	// TODO: initialize before the handler
	projectId := os.Getenv("PUBSUB_PROJECT")
	topicId := os.Getenv("PUBSUB_TOPIC")

	sessionIds, err := logic.GetMissingSessionIds(seasonData.LeagueId, seasonData.SeasonId, irClient, db)
	if err != nil {
		handlers.ReturnException(w, err, "logic.GetMissingSessionIds")
		return
	}

	err = logic.SendSessionsToParse(projectId, topicId, sessionIds)
	if err != nil {
		handlers.ReturnException(w, err, "logic.SendSessionsToParse")
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
