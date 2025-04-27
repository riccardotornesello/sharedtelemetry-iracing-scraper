package processor

import (
	"log"
	"time"

	"github.com/riccardotornesello/irapi-go/api/league"
)

type ProcessedData struct {
	SubsesionID int    `json:"subsessionId"`
	LaunchAt    string `json:"launchAt"`
}

func Process(data *league.LeagueSeasonSessionsResponse, lastSessionLaunchAt *time.Time) ([]ProcessedData, error) {
	log.Println("Analyzing data")

	pubSubData := make([]ProcessedData, 0)
	for _, session := range data.Sessions {
		if lastSessionLaunchAt != nil {
			launchAt, err := time.Parse(time.RFC3339, session.LaunchAt)
			if err != nil {
				return nil, err
			}
			if launchAt.Before(*lastSessionLaunchAt) {
				continue
			}
		}

		pubSubData = append(pubSubData, ProcessedData{
			SubsesionID: session.SubsessionId,
			LaunchAt:    session.LaunchAt,
		})

	}

	return pubSubData, nil
}
