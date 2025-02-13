package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/gorm_utils/database"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
	"riccardotornesello.it/sharedtelemetry/iracing/season_parser/logic"
)

var db *gorm.DB
var irClient *irapi.IRacingApiClient
var pubSubTopic *pubsub.Topic
var pubSubCtx context.Context

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

	pubSubProjectId := os.Getenv("PUBSUB_PROJECT")
	pubSubTopicId := os.Getenv("PUBSUB_TOPIC")

	// Initialize database
	db, err = database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, 2, 2)
	if err != nil {
		log.Fatalf("database.Connect: %v", err)
	}

	// Initialize iRacing client
	irClient, err = irapi.NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		log.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}

	// Initialize Pub/Sub client
	pubSubCtx = context.Background()
	client, err := pubsub.NewClient(pubSubCtx, pubSubProjectId)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	pubSubTopic = client.Topic(pubSubTopicId)

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

	sessionInfo, err := logic.GetMissingSessionInfo(seasonData.LeagueId, seasonData.SeasonId, irClient, db)
	if err != nil {
		handlers.ReturnException(w, err, "logic.GetMissingSessionInfo")
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = logic.StoreMissingSessions(sessionInfo, tx)

	err = logic.SendSessionsToParse(pubSubTopic, pubSubCtx, sessionInfo)
	if err != nil {
		tx.Rollback()
		handlers.ReturnException(w, err, "logic.SendSessionsToParse")
		return
	}

	err = tx.Commit().Error
	if err != nil {
		handlers.ReturnException(w, err, "tx.Commit")
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
