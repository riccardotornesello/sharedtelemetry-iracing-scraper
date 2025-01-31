package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"riccardotornesello.it/sharedtelemetry/iracing/api/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	"riccardotornesello.it/sharedtelemetry/iracing/db/events_models"
)

type Turn struct {
	trackId int
	date    string
}

func TestA(t *testing.T) {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	// Initialize database
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, events_models.AllModels, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	// Common data
	competitionId := 2

	// Get the sessions valid for the competition
	sessions, err := logic.GetCompetitionSessions(db, competitionId)
	if err != nil {
		t.Fatalf("Error getting competition sessions: %v", err)
	}

	// Get drivers
	drivers, err := logic.GetCompetitionDrivers(db, competitionId)
	if err != nil {
		t.Fatalf("Error getting competition drivers: %v", err)
	}

	driverCars := make(map[int]int)
	for _, driver := range drivers {
		driverCars[driver.IRacingCustId] = driver.Crew.IRacingCarId
	}

	// Get laps
	var simsessionIds [][]int
	for _, session := range sessions {
		simsessionIds = append(simsessionIds, []int{session.SubsessionId, session.SimsessionNumber})
	}

	laps, err := logic.GetLaps(db, simsessionIds)
	if err != nil {
		t.Fatalf("Error getting laps: %v", err)
	}

	// Analyze
	allResults := make(map[int]map[int]int)

	currentCustId := 0
	currentSubsessionId := 0
	stintEnd := false
	stintValidLaps := 0
	stintTimeSum := 0

	for _, lap := range laps {
		// Check the first key of driverResults
		if lap.CustID != currentCustId {
			allResults[lap.CustID] = make(map[int]int)
			currentCustId = lap.CustID
			stintEnd = false
			stintValidLaps = 0
			stintTimeSum = 0
		}

		if lap.SubsessionID != currentSubsessionId {
			allResults[lap.CustID][lap.SubsessionID] = 0
			currentSubsessionId = lap.SubsessionID
			stintEnd = false
			stintValidLaps = 0
			stintTimeSum = 0
		}

		driverCar, ok := driverCars[lap.CustID]
		if !ok {
			continue
		}

		if driverCar != lap.SessionSimsessionParticipant.CarID {
			continue
		}

		if stintEnd {
			continue
		}

		if isLapPitted(lap.LapEvents) {
			if stintValidLaps > 0 {
				stintEnd = true
			}

			continue
		}

		if isLapValid(lap.LapNumber, lap.LapTime, lap.LapEvents, lap.Incident) {
			stintValidLaps++
			stintTimeSum += lap.LapTime

			if stintValidLaps == 3 {
				stintEnd = true
				allResults[lap.CustID][lap.SubsessionID] = stintTimeSum / 3 / 10
			}
		} else {
			stintValidLaps = 0
			stintEnd = true
		}
	}

	// Generate CSV the header
	loc, _ := time.LoadLocation("Europe/Rome")
	csv := "Driver,"
	for _, session := range sessions {
		csv += fmt.Sprintf("%s,", session.LaunchAt.In(loc).Format("02/01/2006 15:04:05"))
	}
	csv += "\n"

	// Generate CSV rows
	for _, driver := range drivers {
		csv += fmt.Sprintf("%s,", driver.Name)
		for _, session := range sessions {
			timeString := ""
			lapTime, ok := allResults[driver.IRacingCustId][session.SubsessionId]
			if ok {
				timeString = formatTime(lapTime)
			}
			csv += fmt.Sprintf("%s,", timeString)
		}
		csv += "\n"
	}

	// Save CSV
	file, err := os.Create("results.csv")
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(csv)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	/*
		// TODO: fetch competition info

		leagueId := 4403
		seasonId := 0

		// Get the event groups
		eventGroups, err := logic.GetEventGroups(db, competitionId)
		if err != nil {
			t.Fatalf("Error getting event groups: %v", err)
		}

		// Store the results as driver -> group -> date -> session -> average ms
		driverResults := make(map[int]map[uint]map[string]map[int]int)
		for _, driver := range drivers {
			driverResults[driver.IRacingCustId] = make(map[uint]map[string]map[int]int)
			for _, eventGroup := range eventGroups {
				driverResults[driver.IRacingCustId][eventGroup.ID] = make(map[string]map[int]int)
				for _, date := range eventGroup.Dates {
					driverResults[driver.IRacingCustId][eventGroup.ID][date] = make(map[int]int)
				}
			}
		}

		// Loop through event groups and store the results
		for _, eventGroup := range eventGroups {
			for _, startDate := range eventGroup.Dates {
				// Get the valid session ids
				var simsessionIds [][]int
				sessions, err := logic.GetEventGroupSessions(db, eventGroup.IRacingTrackId, startDate, leagueId, seasonId)
				if err != nil {
					t.Fatalf("Error getting event group sessions: %v", err)
				}
				for _, session := range sessions {
					simsessionIds = append(simsessionIds, []int{session.SubsessionID, session.SimsessionNumber})
				}

				// Get laps
				laps, err := logic.GetLaps(db, simsessionIds)
				if err != nil {
					t.Fatalf("Error getting laps: %v", err)
				}


			}
		}
	*/

	/*
		// Print the best result for each driver in each group and date
		for driver, groupResults := range driverResults {
			log.Printf("Driver %d", driver)
			for _, dateResults := range groupResults {
				for date, sessionResults := range dateResults {
					bestDateResult := 0
					for _, result := range sessionResults {
						if result > 0 && (result < bestDateResult || bestDateResult == 0) {
							bestDateResult = result
						}
					}
					log.Printf("Date %s: %s", date, formatTime(bestDateResult))
				}
			}
			log.Println()
		}
	*/
}

func formatTime(milliseconds int) string {
	// Convert milliseconds to minutes and seconds
	minutes := milliseconds / 60000
	seconds := (milliseconds % 60000) / 1000
	milliseconds = milliseconds % 1000

	// Return always two digits for the seconds and three for the milliseconds
	return fmt.Sprintf("%01d:%02d.%03d", minutes, seconds, milliseconds)
}

func isLapValid(lapNumber int, lapTime int, lapEvents pq.StringArray, incident bool) bool {
	if !(lapNumber > 0 && lapTime > 0 && incident == false) {
		return false
	}

	// Return false if lap.lapEvents contains a blacklisted event
	blacklistedEvents := []string{
		"black flag",
		"car contact",
		"car reset",
		"clock smash",
		"contact",
		"discontinuity",
		"interpolated crossing",
		"invalid",
		"lost control",
		"off track",
		"pitted",
	}

	for _, event := range lapEvents {
		for _, blacklistedEvent := range blacklistedEvents {
			if event == blacklistedEvent {
				return false
			}
		}
	}

	return true
}

func isLapPitted(lapEvents pq.StringArray) bool {
	for _, event := range lapEvents {
		if event == "pitted" {
			return true
		}
	}

	return false
}
