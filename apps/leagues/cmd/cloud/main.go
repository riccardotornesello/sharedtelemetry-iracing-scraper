package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"riccardotornesello.it/sharedtelemetry/iracing/leagues/internal/app"
)

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
	err := app.Run(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Successo.")
}
