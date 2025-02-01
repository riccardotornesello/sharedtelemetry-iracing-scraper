package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"riccardotornesello.it/sharedtelemetry/iracing/api/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	"riccardotornesello.it/sharedtelemetry/iracing/db/events_models"
)

type Rank struct {
	CustId int
	Sum    int
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
	sessions, sessionsMap, err := logic.GetCompetitionSessions(db, competitionId)
	if err != nil {
		t.Fatalf("Error getting competition sessions: %v", err)
	}

	// Get event groups
	eventGroups, err := logic.GetEventGroups(db, competitionId)
	if err != nil {
		t.Fatalf("Error getting competition event groups: %v", err)
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
	bestResults := make(map[int]map[uint]map[string]int) // Customer ID, Group, Date, average ms

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

				averageTime := stintTimeSum / 3 / 10

				// Store the average time of the session for the driver (only valid stints)
				allResults[lap.CustID][lap.SubsessionID] = averageTime

				// Store the best result of the driver for the date in the event group (only valid stints)
				sessionDetails := sessionsMap[lap.SubsessionID]
				// 1. Add the customer to the map if it does not exist
				if _, ok := bestResults[lap.CustID]; !ok {
					bestResults[lap.CustID] = make(map[uint]map[string]int)
				}
				// 2. Add the event group to the map if it does not exist
				if _, ok := bestResults[lap.CustID][sessionDetails.EventGroupId]; !ok {
					bestResults[lap.CustID][sessionDetails.EventGroupId] = make(map[string]int)
				}
				// 3. Add the result to the date if it does not exist or if it is better than the previous one
				if oldResult, ok := bestResults[lap.CustID][sessionDetails.EventGroupId][sessionDetails.Date]; !ok {
					bestResults[lap.CustID][sessionDetails.EventGroupId][sessionDetails.Date] = averageTime
				} else {
					if oldResult > averageTime {
						bestResults[lap.CustID][sessionDetails.EventGroupId][sessionDetails.Date] = averageTime
					}
				}
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

	// Generate the ranking
	ranking := make([]Rank, 0)
	for _, driver := range drivers {
		driverBestResults, ok := bestResults[driver.IRacingCustId]
		if !ok {
			ranking = append(ranking, Rank{driver.IRacingCustId, 0})
			continue
		}

		sum := 0
		isValid := true
		for _, eventGroup := range eventGroups {
			if driverBestGroupResults, ok := driverBestResults[eventGroup.ID]; !ok {
				// If the driver did not participate in the event group, the result is 0
				isValid = false
				break
			} else {
				// Check if the driver has at least a result in one date of the event group and in case add the best result
				bestResult := 0
				for _, result := range driverBestGroupResults {
					if bestResult == 0 || result < bestResult {
						bestResult = result
					}
				}

				if bestResult > 0 {
					sum += bestResult
				} else {
					isValid = false
					break
				}
			}
		}

		if !isValid {
			ranking = append(ranking, Rank{driver.IRacingCustId, 0})
		} else {
			ranking = append(ranking, Rank{driver.IRacingCustId, sum})
		}
	}

	// Sort the ranking by sum. If the sum is 0, put the driver at the end of the ranking
	sort.Slice(ranking, func(i, j int) bool {
		if ranking[i].Sum == 0 {
			return false
		}

		if ranking[j].Sum == 0 {
			return true
		}

		return ranking[i].Sum < ranking[j].Sum
	})

	// Generate ranking CSV
	rankingCsv := "Driver,Sum\n"
	for _, rank := range ranking {
		driver, ok := drivers[rank.CustId]
		if !ok {
			continue
		}

		rankingCsv += fmt.Sprintf("%s,%s\n", driver.Name, formatTime(rank.Sum))
	}

	file, err = os.Create("ranking.csv")
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}

	_, err = file.WriteString(rankingCsv)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
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
