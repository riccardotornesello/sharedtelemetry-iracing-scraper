package irapi

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestCall(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	iRacingEmail := os.Getenv("IRACING_EMAIL")
	iRacingPassword := os.Getenv("IRACING_PASSWORD")

	client, err := NewIRacingApiClient(iRacingEmail, iRacingPassword)
	if err != nil {
		t.Fatalf("irapi.NewIRacingApiClient: %v", err)
	}

	_, err = client.GetLeague(4403, false)
	if err != nil {
		t.Fatalf("client.GetLeague: %v", err)
	}
}
