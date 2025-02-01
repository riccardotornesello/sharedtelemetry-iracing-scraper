package logic

import (
	"fmt"
	"time"

	"riccardotornesello.it/sharedtelemetry/iracing/api/utils"
	"riccardotornesello.it/sharedtelemetry/iracing/db/events_models"
)

func GenerateSessionsCsv(sessions []*CompetitionSession, drivers []*events_models.CompetitionDriver, allResults map[int]map[int]int) string {
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
				timeString = utils.FormatTime(lapTime)
			}
			csv += fmt.Sprintf("%s,", timeString)
		}
		csv += "\n"
	}

	return csv
}
