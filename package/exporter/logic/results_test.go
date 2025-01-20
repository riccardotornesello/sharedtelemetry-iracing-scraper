package logic

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
)

func TestExportSessionResults(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Get configuration
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	// Initialize database
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, nil, 20, 2)
	if err != nil {
		t.Errorf("database.Connect: %v", err)
	}

	// Get the data
	results, sessions, dates, drivers, err := ExportSessionResults(db, []string{"2024-09-11", "2024-09-12", "2024-09-14", "2024-09-15"})
	if err != nil {
		t.Errorf("ExportSessionResults: %v", err)
	}

	// Save the results in a .csv file. The columns are the session dates and the rows are the drivers
	file, err := os.Create("results.csv")
	if err != nil {
		t.Errorf("os.Create: %v", err)
	}
	defer file.Close()

	// Write the header
	_, err = file.WriteString("Driver ID,")
	if err != nil {
		t.Errorf("file.WriteString: %v", err)
	}
	for _, sessionId := range sessions {
		_, err = file.WriteString(fmt.Sprintf("%s,", dates[sessionId]))
		if err != nil {
			t.Errorf("file.WriteString: %v", err)
		}
	}
	_, err = file.WriteString("\n")
	if err != nil {
		t.Errorf("file.WriteString: %v", err)
	}

	// Write the results
	for _, driver := range drivers {
		_, err = file.WriteString(fmt.Sprintf("%d,", driver))
		if err != nil {
			t.Errorf("file.WriteString: %v", err)
		}
		for _, sessionId := range sessions {
			time := ""
			realTime, ok := results[sessionId][driver]
			if ok {
				fmt.Println(realTime)
				time = formatTime(realTime)
			}
			_, err = file.WriteString(fmt.Sprintf("%s,", time))
			if err != nil {
				t.Errorf("file.WriteString: %v", err)
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			t.Errorf("file.WriteString: %v", err)
		}
	}
}

func formatTime(milliseconds int) string {
	// Convert milliseconds to minutes and seconds
	minutes := milliseconds / 60000
	seconds := (milliseconds % 60000) / 1000
	milliseconds = milliseconds % 1000

	// Return always two digits for the seconds and three for the milliseconds
	return fmt.Sprintf("%01d:%02d.%03d", minutes, seconds, milliseconds)
}
