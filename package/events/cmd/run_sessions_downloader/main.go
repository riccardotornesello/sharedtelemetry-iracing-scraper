package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
	common "riccardotornesello.it/sharedtelemetry/iracing/common/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/events/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/events/models"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

var db *gorm.DB
var irClient *irapi.IRacingApiClient

func main() {
	var err error

	db, irClient, err = common.InitCloudRun(models.AllModels)
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

type SessionData struct {
	SubsessionId int    `json:"subsessionId"`
	LaunchAt     string `json:"launchAt"`
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

	var sessionData SessionData
	if err := json.Unmarshal(m.Message.Data, &sessionData); err != nil {
		log.Printf("json.Unmarshal data: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	launchAt, err := time.Parse(time.RFC3339, sessionData.LaunchAt)
	if err != nil {
		log.Printf("time.Parse: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := logic.ParseSession(irClient, sessionData.SubsessionId, launchAt, db); err != nil {
		log.Printf("logic.ParseSession: %v", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)

	return
}
