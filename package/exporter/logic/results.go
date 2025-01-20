package logic

import (
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func ExportSessionResults(db *gorm.DB, dates []string) (map[int]map[int]int, []int, map[int]string, []int, error) {
	// Get the valid sessions
	// TODO: validate the track
	// TODO: use UTC
	sessionRows, err := db.Table("sessions").
		Where("(sessions.launch_at AT TIME ZONE 'CET')::date IN (?)", dates).
		Select("subsession_id", "launch_at").
		Order("launch_at").
		Rows()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer sessionRows.Close()

	sessionDates := make(map[int]string)
	sessions := make([]int, 0)
	for sessionRows.Next() {
		var subsessionId int
		var launchAt string
		err := sessionRows.Scan(&subsessionId, &launchAt)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		sessionDates[subsessionId] = launchAt
		sessions = append(sessions, subsessionId)
	}

	// Get the laps
	rows, err := db.Table("laps").
		Joins("JOIN session_simsession_participants ON laps.subsession_id = session_simsession_participants.subsession_id AND laps.simsession_number = session_simsession_participants.simsession_number AND laps.cust_id = session_simsession_participants.cust_id").
		Joins("JOIN session_simsessions ON laps.subsession_id = session_simsessions.subsession_id AND laps.simsession_number = session_simsessions.simsession_number").
		Where("laps.simsession_number = 0 AND session_simsessions.simsession_name = 'QUALIFY' AND laps.subsession_id IN (?)", sessions).
		Select("laps.cust_id", "laps.lap_time", "laps.lap_events", "laps.incident", "laps.lap_number", "laps.subsession_id").
		Order("laps.subsession_id, laps.simsession_number, laps.cust_id, laps.lap_number").
		Rows()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	// TODO: check car and drivers

	// Iterate over the laps and extract the best lap for each driver in each session
	currentSubsessionId := 0
	currentCustId := 0
	stintEnd := false
	stintValidLaps := 0
	stintTimeSum := 0
	subsessionId := 0

	allResults := make(map[int]map[int]int)
	drivers := make(map[int]interface{})

	for rows.Next() {
		// var lap models.Lap
		// db.ScanRows(rows, &lap)

		custId := 0
		lapTime := 0
		lapEvents := pq.StringArray{}
		incident := false
		lapNumber := 0

		err := rows.Scan(&custId, &lapTime, &lapEvents, &incident, &lapNumber, &subsessionId)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		if subsessionId != currentSubsessionId {
			allResults[subsessionId] = make(map[int]int)
			currentSubsessionId = subsessionId
			stintEnd = false
			stintValidLaps = 0
			stintTimeSum = 0
		}

		if custId != currentCustId {
			allResults[subsessionId][custId] = 0
			currentCustId = custId
			drivers[custId] = nil
			stintEnd = false
			stintValidLaps = 0
			stintTimeSum = 0
		}

		if stintEnd {
			continue
		}

		if isLapPitted(lapEvents) {
			if stintValidLaps > 0 {
				stintEnd = true
			}

			continue
		}

		if isLapValid(lapNumber, lapTime, lapEvents, incident) {
			stintValidLaps++
			stintTimeSum += lapTime

			if stintValidLaps == 3 {
				stintEnd = true
				allResults[subsessionId][custId] = stintTimeSum / 3 / 10
			}
		} else {
			stintValidLaps = 0
			stintEnd = true
		}
	}

	driverIds := make([]int, 0)
	for driverId := range drivers {
		driverIds = append(driverIds, driverId)
	}

	return allResults, sessions, sessionDates, driverIds, nil
}

func GenerateSessionsCsv(db *gorm.DB) (string, error) {
	// TODO: variable dates

	// Extract the session results
	results, sessions, sessionDates, drivers, err := ExportSessionResults(db, []string{"2024-09-11", "2024-09-12", "2024-09-14", "2024-09-15"})
	if err != nil {
		return "", err
	}

	// Generate the header
	csv := "Driver ID,"
	for _, sessionId := range sessions {
		csv += fmt.Sprintf("%s,", sessionDates[sessionId])
	}
	csv += "\n"

	// Write the results
	for _, driver := range drivers {
		csv += fmt.Sprintf("%d,", driver)
		for _, sessionId := range sessions {
			timeString := ""
			lapTime, ok := results[sessionId][driver]
			if ok {
				timeString = formatTime(lapTime)
			}
			csv += fmt.Sprintf("%s,", timeString)
		}
		csv += "\n"
	}

	return csv, nil
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

func formatTime(milliseconds int) string {
	// Convert milliseconds to minutes and seconds
	minutes := milliseconds / 60000
	seconds := (milliseconds % 60000) / 1000
	milliseconds = milliseconds % 1000

	// Return always two digits for the seconds and three for the milliseconds
	return fmt.Sprintf("%01d:%02d.%03d", minutes, seconds, milliseconds)
}
