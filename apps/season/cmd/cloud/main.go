package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"riccardotornesello.it/sharedtelemetry/iracing/season/internal/app"
)

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

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server Cloud Run in ascolto su :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
		return
	}

	var m PubSubMessage
	if err := json.Unmarshal(body, &m); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling PubSub message: %v", err), http.StatusBadRequest)
		return
	}

	var seasonData SeasonData
	if err := json.Unmarshal(m.Message.Data, &seasonData); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling season data: %v", err), http.StatusBadRequest)
		return
	}

	err = app.Run(r.Context(), seasonData.LeagueId, seasonData.SeasonId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Done")
}
