package sessions_downloader

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

func parseSession(irClient *client.IRacingApiClient, session *client.LeagueSeasonSession, db *gorm.DB, saveRequests bool) error {
	log.Println("Parsing session", session.SubsessionId)

	// Get the results and make sure the session is a qualify session
	results, err := irClient.GetResults(session.SubsessionId)
	if err != nil {
		return err
	}

	if saveRequests {
		resultsJson, _ := json.Marshal(results)
		err := os.WriteFile(fmt.Sprintf("downloads/sessions/%d.json", session.SubsessionId), resultsJson, 0644)
		if err != nil {
			return err
		}
	}

	// DB: create a new transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	launchAt, err := time.Parse(time.RFC3339, session.LaunchAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(&models.Event{
		SubsessionId: session.SubsessionId,
		LaunchAt:     launchAt,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, result := range results.SessionResults {
		eventSesion := models.EventSession{
			EventID:          session.SubsessionId,
			SimsessionNumber: result.SimsessionNumber,
			SimsessionType:   result.SimsessionType,
			SimsessionName:   result.SimsessionName,
		}
		err = tx.Create(&eventSesion).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// Start 3 workers to get the lap results for each driver
		numWorkers := 3
		numJobs := len(result.Results)
		lapJobsInput := make(chan int, numJobs)
		lapJobsOutput := make(chan *ParseLapMessage, numJobs)
		for w := 0; w < numWorkers; w++ {
			go lapResultsWorker(session.SubsessionId, irClient, lapJobsInput, lapJobsOutput, tx, saveRequests)
		}

		// Get the single laps for each driver to check if the lap is valid
		for _, result := range result.Results {
			lapJobsInput <- result.CustId
		}

		for a := 0; a < numJobs; a++ {
			parseLapMessage := <-lapJobsOutput
			if parseLapMessage.Error != nil {
				// Empty the input channel to stop pending workers
				cleanedJobs := 0
				for len(lapJobsInput) > 0 {
					<-lapJobsInput
					cleanedJobs++
				}

				// Empty the output channel to stop pending workers
				for i := 0; i < numJobs-a-cleanedJobs-1; i++ {
					<-lapJobsOutput
				}

				close(lapJobsOutput)
				close(lapJobsInput)

				tx.Rollback()
				return parseLapMessage.Error
			}

			// Store the laps in the database
			if len(parseLapMessage.Message) > 0 {
				dbData := make([]models.Lap, len(parseLapMessage.Message))
				for i, lap := range parseLapMessage.Message {
					dbData[i] = models.Lap{
						EventSessionID: eventSesion.ID,
						CustId:         lap.CustId,
						LapEvents:      strings.Join(lap.LapEvents, ","), // TODO: store as array
						Incident:       lap.Incident,
						LapTime:        lap.LapTime,
						LapNumber:      lap.LapNumber,
					}
				}

				err = tx.Create(dbData).Error
				if err != nil {
					// Empty the input channel to stop pending workers
					cleanedJobs := 0
					for len(lapJobsInput) > 0 {
						<-lapJobsInput
						cleanedJobs++
					}

					// Empty the output channel to stop pending workers
					for i := 0; i < numJobs-a-cleanedJobs-1; i++ {
						<-lapJobsOutput
					}

					close(lapJobsOutput)
					close(lapJobsInput)

					tx.Rollback()
					return err
				}
			}
		}
		close(lapJobsOutput)
		close(lapJobsInput)
	}

	return tx.Commit().Error
}

func lapResultsWorker(subsessionId int, irClient *client.IRacingApiClient, lapJobsInput <-chan int, lapJobsOutput chan<- *ParseLapMessage, tx *gorm.DB, saveRequests bool) {
	for driverId := range lapJobsInput {
		message := parseLap(subsessionId, irClient, driverId, tx, saveRequests)
		lapJobsOutput <- message
	}
}
