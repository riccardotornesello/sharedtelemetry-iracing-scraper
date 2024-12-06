package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
	"riccardotornesello.it/iracing-average-lap/client"
)

type ParseLapMessage struct {
	Message []client.ResultsLapDataChunk
	Error   error
}

func parseLap(subsessionId int, irClient *client.IRacingApiClient, driverId int, tx *gorm.DB, saveRequests bool) *ParseLapMessage {
	lapResults, err := irClient.GetResultsLapData(subsessionId, 0, driverId)
	if err != nil {
		return &ParseLapMessage{Error: err}
	}

	if saveRequests {
		lapResultsJson, _ := json.Marshal(lapResults)
		err := os.WriteFile(fmt.Sprintf("downloads/laps/%d_%d.json", subsessionId, driverId), lapResultsJson, 0644)
		if err != nil {
			return &ParseLapMessage{Error: err}
		}
	}

	return &ParseLapMessage{Message: lapResults.Laps}

	// noBadEvents := true
	// for _, event := range lap.LapEvents {
	// 	if event == "off track" || event == "pitted" || event == "invalid" {
	// 		noBadEvents = false
	// 	}
	// }
	// if lap.Incident == false && noBadEvents && lap.LapTime > 0 && lap.LapNumber > 0 {
	// 	validLaps = append(validLaps, lap.LapTime)
	// }

	// validLaps = append(validLaps, &models.Lap{
	// 	EventSessionID: eventSessionId,
	// 	CustId:         driverId,
	// 	LapEvents:      strings.Join(lap.LapEvents, ","), // TODO: store as array
	// 	Incident:       lap.Incident,
	// 	LapTime:        lap.LapTime,
	// 	LapNumber:      lap.LapNumber,
	// })
}
