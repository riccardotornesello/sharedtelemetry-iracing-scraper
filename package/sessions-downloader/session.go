package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"riccardotornesello.it/iracing-average-lap/client"
	"riccardotornesello.it/iracing-average-lap/models"
)

func sessionWorker(irClient *client.IRacingApiClient, sessionJobs <-chan *client.LeagueSeasonSession, sessionJobResults chan<- interface{}, db *gorm.DB, saveRequests bool) {
	for session := range sessionJobs {
		// Get the results and make sure the session is a qualify session
		results := irClient.GetResults(session.SubsessionId)

		if saveRequests {
			resultsJson, _ := json.Marshal(results)
			err := os.WriteFile(fmt.Sprintf("downloads/sessions/%d.json", session.SubsessionId), resultsJson, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		// DB: create a new transaction
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		launchAt, error := time.Parse(time.RFC3339, session.LaunchAt)
		if error != nil {
			log.Fatal("Error parsing date")
		}

		tx.Create(&models.Event{
			SubsessionId: session.SubsessionId,
			LaunchAt:     launchAt,
		})

		for _, result := range results.SessionResults {
			eventSesion := models.EventSession{
				EventID:          session.SubsessionId,
				SimsessionNumber: result.SimsessionNumber,
				SimsessionType:   result.SimsessionType,
				SimsessionName:   result.SimsessionName,
			}
			tx.Create(&eventSesion)

			// Start 3 workers to get the lap results for each driver
			numWorkers := 3
			numJobs := len(result.Results)
			lapJobs := make(chan int, numJobs)
			lapJobResults := make(chan interface{}, numJobs)
			for w := 0; w < numWorkers; w++ {
				go lapResultsWorker(session.SubsessionId, irClient, lapJobs, lapJobResults, tx, eventSesion.ID, saveRequests)
			}

			// Get the single laps for each driver to check if the lap is valid
			for _, result := range result.Results {
				lapJobs <- result.CustId
			}
			close(lapJobs)

			for a := 0; a < numJobs; a++ {
				<-lapJobResults
			}
			close(lapJobResults)
		}

		tx.Commit()

		sessionJobResults <- interface{}(nil)
	}
}

func lapResultsWorker(subsessionId int, irClient *client.IRacingApiClient, jobs <-chan int, jobResults chan<- interface{}, tx *gorm.DB, eventSessionId uint, saveRequests bool) {
	for driverId := range jobs {
		lapResults := irClient.GetResultsLapData(subsessionId, 0, driverId)

		if saveRequests {
			lapResultsJson, _ := json.Marshal(lapResults)
			err := os.WriteFile(fmt.Sprintf("downloads/laps/%d_%d.json", subsessionId, driverId), lapResultsJson, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, lap := range lapResults.Laps {
			// noBadEvents := true
			// for _, event := range lap.LapEvents {
			// 	if event == "off track" || event == "pitted" || event == "invalid" {
			// 		noBadEvents = false
			// 	}
			// }
			// if lap.Incident == false && noBadEvents && lap.LapTime > 0 && lap.LapNumber > 0 {
			// 	validLaps = append(validLaps, lap.LapTime)
			// }

			tx.Create(&models.Lap{
				EventSessionID: eventSessionId,
				CustId:         driverId,
				LapEvents:      strings.Join(lap.LapEvents, ","), // TODO: store as array
				Incident:       lap.Incident,
				LapTime:        lap.LapTime,
				LapNumber:      lap.LapNumber,
			})
		}

		jobResults <- interface{}(nil)
	}
}
