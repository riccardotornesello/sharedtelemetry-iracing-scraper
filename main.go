package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"riccardotornesello.it/iracing-average-lap/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	session_days := [2]string{"2024-09-11", "2024-09-14"}

	irClient := client.NewIRacingApiClient(os.Getenv("IRACING_EMAIL"), os.Getenv("IRACING_PASSWORD"))

	seriesId, err := strconv.Atoi(os.Getenv("IRACING_SERIES_ID"))
	if err != nil {
		log.Fatal("Invalid series id")
	}

	sessions := irClient.GetLeagueSeasonSessions(seriesId, 0, true)

	for _, session := range sessions.Sessions {
		startDate := session.LaunchAt[:10]
		if startDate != session_days[0] && startDate != session_days[1] {
			continue
		}

		results := irClient.GetResults(session.SubsessionId)
		qualifyResults := results.SessionResults[1]

		if qualifyResults.SimsessionNumber != 0 {
			log.Fatal("wrong simsession number")
		}
		if qualifyResults.SimsessionName != "QUALIFY" {
			log.Fatal("wrong simsession name")
		}

		drivers := make([]int, 0)
		for _, qResult := range qualifyResults.Results {
			if qResult.LapsComplete >= 3 && qResult.BestLapTime > 0 {
				drivers = append(drivers, qResult.CustId)
			}
		}

		for _, driverId := range drivers {
			lapResults := irClient.GetResultsLapData(session.SubsessionId, 0, driverId)

			validLaps := make([]int, 0)
			for _, lap := range lapResults.Laps {
				noBadEvents := true
				for _, event := range lap.LapEvents {
					if event == "off track" || event == "pitted" || event == "invalid" {
						noBadEvents = false
					}
				}

				if lap.Incident == false && noBadEvents && lap.LapTime > 0 && lap.LapNumber > 0 {
					validLaps = append(validLaps, lap.LapTime)
				}
			}

			if len(validLaps) == 3 {
				average := (validLaps[0] + validLaps[1] + validLaps[2]) / 3 / 10
				log.Printf("Driver %d: %d:%.4f", driverId, average/60000, float64(average%60000)/10000)
			}
		}
	}
}
