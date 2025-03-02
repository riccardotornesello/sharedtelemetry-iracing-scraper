package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
	"riccardotornesello.it/sharedtelemetry/iracing/season_parser/logic"
)

const projectID = "sharedtelemetryapp" // TODO: move to env

var irClient *irapi.IRacingApiClient
var firestoreClient *firestore.Client
var firestoreContext context.Context
var pubSubTopic *pubsub.Topic
var pubSubCtx context.Context

func main() {
	var err error

	// Get configuration
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	pubSubProjectId := os.Getenv("PUBSUB_PROJECT")
	pubSubTopicId := os.Getenv("PUBSUB_TOPIC")

	// Initialize database
	firestoreContext = context.Background()
	firebaseConf := &firebase.Config{ProjectID: projectID}
	firebaseApp, err := firebase.NewApp(firestoreContext, firebaseConf)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err = firebaseApp.Firestore(firestoreContext)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()

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

	firstLaunchAt, err := logic.GetFirstSessionLaunchAt(seasonData.LeagueId, seasonData.SeasonId, firestoreClient, firestoreContext)
	if err != nil {
		handlers.ReturnException(w, err, "logic.GetFirstSessionLaunchAt")
		return
	}

	seasonSessionsInfo, err := logic.GetLeagueSeasonSessionsInfo(seasonData.LeagueId, seasonData.SeasonId, firstLaunchAt, irClient)
	if err != nil {
		handlers.ReturnException(w, err, "logic.GetLeagueSeasonSessionsInfo")
		return
	}

	err = logic.SendSessionsToParse(pubSubTopic, pubSubCtx, seasonSessionsInfo)
	if err != nil {
		handlers.ReturnException(w, err, "logic.SendSessionsToParse")
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
