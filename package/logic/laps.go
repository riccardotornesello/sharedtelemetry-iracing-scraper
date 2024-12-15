package logic

import (
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

type ParseLapMessage struct {
	Message []irapi.ResultsLapDataChunk
	Error   error
}

func parseLap(subsessionId int, irClient *irapi.IRacingApiClient, driverId int, tx *gorm.DB, saveRequests bool) *ParseLapMessage {
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
}
