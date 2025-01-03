package models

import (
	"testing"

	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
)

func TestMigration(t *testing.T) {
	_, err := database.Connect("postgres", "example", "127.0.0.1", "5432", "postgres", AllModels, 1, 0)
	if err != nil {
		t.Fatalf("database.Connect: %v", err)
	}
}
