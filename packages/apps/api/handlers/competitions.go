package handlers

import (
	"errors"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/api/logic"
)

type RankingResponse struct {
	Classes     []*ClassInfo        `json:"classes"`
	Ranking     []*Rank             `json:"ranking"`
	Drivers     map[int]*DriverInfo `json:"drivers"`
	EventGroups []*EventGroupInfo   `json:"eventGroups"`
	Competition *CompetitionInfo    `json:"competition"`
}

type Rank struct {
	Pos     int                     `json:"pos"`
	CustId  int                     `json:"custId"`
	Sum     int                     `json:"sum"`
	IsValid bool                    `json:"isValid"`
	Results map[uint]map[string]int `json:"results"`
}

type TeamInfo struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type CrewInfo struct {
	Id           uint     `json:"id"`
	Name         string   `json:"name"`
	CarId        int      `json:"carId"`
	Team         TeamInfo `json:"team"`
	ClassId      uint     `json:"classId"`
	CarModel     string   `json:"carModel"`
	CarBrandIcon string   `json:"carBrandIcon"`
}

type ClassInfo struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Index int    `json:"index"`
}

type DriverInfo struct {
	CustId    int      `json:"custId"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Crew      CrewInfo `json:"crew"`
}

type EventGroupInfo struct {
	Id      uint     `json:"id"`
	Name    string   `json:"name"`
	TrackId int      `json:"trackId"`
	Dates   []string `json:"dates"`
}

type CompetitionInfo struct {
	Id               uint   `json:"id"`
	Name             string `json:"name"`
	CrewDriversCount int    `json:"crewDriversCount"`
}

func CompetitionRankingHandler(c *gin.Context, eventsDb *gorm.DB, carsDb *gorm.DB) {
	// Get the competition
	competition, err := logic.GetCompetitionBySlug(eventsDb, c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Competition not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition"})
			return
		}
	}

	// Get the sessions valid for the competition
	sessions, sessionsMap, err := logic.GetCompetitionSessions(eventsDb, competition.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition sessions"})
		return
	}

	// Get event groups
	eventGroups, err := logic.GetEventGroups(eventsDb, competition.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting event groups"})
		return
	}

	// Get drivers
	drivers, _, err := logic.GetCompetitionDrivers(eventsDb, competition.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition drivers"})
		return
	}

	driverCars := make(map[int]int)
	allowedCars := make(map[int]bool)
	for _, driver := range drivers {
		driverCars[driver.IRacingCustId] = driver.Crew.IRacingCarId
		allowedCars[driver.Crew.IRacingCarId] = true
	}

	// Get laps
	var simsessionIds [][]int
	for _, session := range sessions {
		simsessionIds = append(simsessionIds, []int{session.SubsessionId, session.SimsessionNumber})
	}

	laps, err := logic.GetLaps(eventsDb, simsessionIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting laps"})
		return
	}

	// Get cars
	allwedCarIds := make([]int, 0)
	for carId := range allowedCars {
		allwedCarIds = append(allwedCarIds, carId)
	}

	carBrands, err := logic.GetCarBrands(carsDb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting car brands"})
		return
	}

	carModels, err := logic.GetCarModelsById(carsDb, allwedCarIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting car models"})
		return
	}

	// Get classes
	classes, err := logic.GetCompetitionClasses(eventsDb, competition.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition classes"})
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
			IsValid: true,
			Results: bestResults[driver.IRacingCustId], // TODO: add default value, it might be null
		}

		driverBestResults, ok := bestResults[driver.IRacingCustId]
		if ok {
			for _, eventGroup := range eventGroups {
				if driverBestGroupResults, ok := driverBestResults[eventGroup.ID]; !ok {
					// If the driver did not participate in the event group, the result is 0
					driverRank.IsValid = false
				} else {
					// Check if the driver has at least a result in one date of the event group and in case add the best result
					bestResult := 0
					for _, result := range driverBestGroupResults {
						if bestResult == 0 || result < bestResult {
							bestResult = result
						}
					}

					if bestResult > 0 {
						driverRank.Sum += bestResult
					} else {
						driverRank.IsValid = false
					}
				}
			}
		}

		if driverRank.Sum == 0 {
			driverRank.IsValid = false
		}

		ranking = append(ranking, driverRank)
	}

	// Sort the ranking by sum. First the valid ones, then the invalid ones and the ones with 0 sum
	sort.Slice(ranking, func(i, j int) bool {
		if ranking[i].IsValid != ranking[j].IsValid {
			return ranking[i].IsValid
		}
		if ranking[i].Sum == 0 {
			return false
		}
		if ranking[j].Sum == 0 {
			return true
		}
		return ranking[i].Sum < ranking[j].Sum
	})

	for i, driver := range ranking {
		driver.Pos = i + 1

		if !driver.IsValid {
			driver.Sum = 0
		}
	}

	// Return the response
	driversInfo := make(map[int]*DriverInfo)
	for _, driver := range drivers {
		carModel := ""
		carBrandIcon := ""

		car, ok := carModels[driver.Crew.IRacingCarId]
		if ok {
			carModel = car.Name

			brand, ok := carBrands[car.Brand]
			if ok {
				carBrandIcon = brand.Icon
			}
		}

		driverInfo := &DriverInfo{
			CustId:    driver.IRacingCustId,
			FirstName: driver.FirstName,
			LastName:  driver.LastName,
			Crew: CrewInfo{
				Id:           driver.Crew.ID,
				Name:         driver.Crew.Name,
				CarId:        driver.Crew.IRacingCarId,
				CarModel:     carModel,
				CarBrandIcon: carBrandIcon,
				ClassId:      driver.Crew.ClassID,
				Team: TeamInfo{
					Id:      driver.Crew.Team.ID,
					Name:    driver.Crew.Team.Name,
					Picture: driver.Crew.Team.Picture,
				},
			},
		}

		driversInfo[driver.IRacingCustId] = driverInfo
	}

	eventGroupsInfo := make([]*EventGroupInfo, 0)
	for _, eventGroup := range eventGroups {
		eventGroupInfo := &EventGroupInfo{
			Id:      eventGroup.ID,
			Name:    eventGroup.Name,
			TrackId: eventGroup.IRacingTrackId,
			Dates:   eventGroup.Dates,
		}

		eventGroupsInfo = append(eventGroupsInfo, eventGroupInfo)
	}

	competitionInfo := &CompetitionInfo{
		Id:               competition.ID,
		Name:             competition.Name,
		CrewDriversCount: competition.CrewDriversCount,
	}

	classesInfo := make([]*ClassInfo, len(classes))
	for i, class := range classes {
		classesInfo[i] = &ClassInfo{
			Id:    class.ID,
			Name:  class.Name,
			Color: class.Color,
			Index: class.Index,
		}
	}

	response := RankingResponse{
		Classes:     classesInfo,
		Ranking:     ranking,
		EventGroups: eventGroupsInfo,
		Drivers:     driversInfo,
		Competition: competitionInfo,
	}

	c.JSON(http.StatusOK, response)
}

func CompetitionCsvHandler(c *gin.Context, eventsDb *gorm.DB) {
	// Get the competition
	competition, err := logic.GetCompetitionBySlug(eventsDb, c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Competition not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition"})
			return
		}
	}

	// Get the sessions valid for the competition
	sessions, sessionsMap, err := logic.GetCompetitionSessions(eventsDb, competition.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting competition sessions"})
		return
	}

	// Get drivers
	drivers, _, err := logic.GetCompetitionDrivers(eventsDb, competition.ID)
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

	laps, err := logic.GetLaps(eventsDb, simsessionIds)
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

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=sessions.csv")
	c.Data(http.StatusOK, "text/csv", []byte(csv))
}
