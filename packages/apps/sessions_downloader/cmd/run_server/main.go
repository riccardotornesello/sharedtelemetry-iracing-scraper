package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils/handlers"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions_downloader/logic"
)

var irClient *irapi.IRacingApiClient
var firestoreClient *firestore.Client
var firestoreContext context.Context

const projectID = "sharedtelemetryapp" // TODO: move to env

func main() {
	var err error

	// Get configuration
	godotenv.Load()

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

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

type SessionData struct {
	SubsessionId int    `json:"subsessionId"`
	LaunchAt     string `json:"launchAt"`
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

	var sessionData SessionData
	if err := json.Unmarshal(m.Message.Data, &sessionData); err != nil {
		handlers.ReturnException(w, err, "json.Unmarshal")
		return
	}

	launchAt, err := time.Parse(time.RFC3339, sessionData.LaunchAt)
	if err != nil {
		handlers.ReturnException(w, err, "time.Parse")
		return
	}

	if err := logic.ParseSession(irClient, sessionData.SubsessionId, launchAt, firestoreClient, firestoreContext, 10); err != nil {
		handlers.ReturnException(w, err, "logic.ParseSession")
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
