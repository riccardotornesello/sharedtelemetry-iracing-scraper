package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := Connect("postgres", "example", "127.0.0.1", "5432", "postgres")

	if err != nil {
		t.Fatalf("Connect(...) error: %v", err)
	}
}
