package main

import (
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"riccardotornesello.it/sharedtelemetry/iracing/api/logic"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
	models "riccardotornesello.it/sharedtelemetry/iracing/db/events_models"
)

type RankingResponse struct {
	Ranking     []*Rank                           `json:"ranking"`
	Drivers     map[int]*models.CompetitionDriver `json:"drivers"`
	EventGroups []*models.EventGroup              `json:"eventGroups"`
}

type Rank struct {
	CustId  int                     `json:"custId"`
	Sum     int                     `json:"sum"`
	Results map[uint]map[string]int `json:"results"`
}

func main() {
	// Get configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	// Initialize database
	db, err := database.Connect(dbUser, dbPass, dbHost, dbPort, dbName, models.AllModels, 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/competitions/:id/ranking", func(c *gin.Context) {
		competitionIdParam := c.Param("id")
		competitionId, err := strconv.Atoi(competitionIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition id"})
			return
		}

		// Get the sessions valid for the competition
		sessions, sessionsMap, err := logic.GetCompetitionSessions(db, competitionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition sessions"})
			return
		}

		// Get event groups
		eventGroups, err := logic.GetEventGroups(db, competitionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting event groups"})
			return
		}

		// Get drivers
		drivers, driversMap, err := logic.GetCompetitionDrivers(db, competitionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition drivers"})
			return
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting laps"})
			return
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

			if logic.IsLapPitted(lap.LapEvents) {
				if stintValidLaps > 0 {
					stintEnd = true
				}

				continue
			}

			if logic.IsLapValid(lap.LapNumber, lap.LapTime, lap.LapEvents, lap.Incident) {
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

		// Generate the ranking
		ranking := make([]*Rank, 0)
		for _, driver := range drivers {
			driverRank := &Rank{
				CustId:  driver.IRacingCustId,
				Sum:     0,
				Results: bestResults[driver.IRacingCustId],
			}

			driverBestResults, ok := bestResults[driver.IRacingCustId]
			if !ok {
				ranking = append(ranking, driverRank)
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

			if isValid {
				driverRank.Sum = sum
			}

			ranking = append(ranking, driverRank)
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

		response := RankingResponse{
			Ranking:     ranking,
			EventGroups: eventGroups,
			Drivers:     driversMap,
		}

		c.JSON(http.StatusOK, response)
	})

	r.GET("/competitions/:id/csv", func(c *gin.Context) {
		competitionIdParam := c.Param("id")
		competitionId, err := strconv.Atoi(competitionIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition id"})
			return
		}

		// Get the sessions valid for the competition
		sessions, sessionsMap, err := logic.GetCompetitionSessions(db, competitionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition sessions"})
			return
		}

		// Get drivers
		drivers, _, err := logic.GetCompetitionDrivers(db, competitionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition drivers"})
			return
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting laps"})
			return
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

			if logic.IsLapPitted(lap.LapEvents) {
				if stintValidLaps > 0 {
					stintEnd = true
				}

				continue
			}

			if logic.IsLapValid(lap.LapNumber, lap.LapTime, lap.LapEvents, lap.Incident) {
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

		// Generate CSV
		csv := logic.GenerateSessionsCsv(sessions, drivers, allResults)

		c.String(http.StatusOK, csv)
	})

	r.Run()
}
