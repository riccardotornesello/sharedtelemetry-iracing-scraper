package logic

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/events_models"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

func ParseSession(irClient *irapi.IRacingApiClient, subsessionId int, subsessionLaunchAt time.Time, db *gorm.DB, workers int) error {
	// Check the info already in the database
	var dbSession events_models.Session
	err := db.Where("subsession_id = ?", subsessionId).First(&dbSession).Error
	if err != nil {
		return err
	}

	// If the session is already parsed, return
	// TODO: check by parse date
	if dbSession.TrackID != 0 {
		slog.Info("Session %d already parsed", subsessionId)
		return nil
	}

	// Get the whole session results to extract simsessions and participants
	results, err := irClient.GetResults(subsessionId)
	if err != nil {
		return fmt.Errorf("error getting results for session %d: %w", subsessionId, err)
	}

	// For each simsession, get the results for each driver
	// results.SessionResults: one for each simsession (practice, quali...)
	// results.SessionResults[i].Results: one for each driver
	tasksCount := 0
	for _, simSessionResult := range results.SessionResults {
		tasksCount += len(simSessionResult.Results)
	}

	tasksChan := make(chan sessionLapTask, tasksCount)
	resultsChan := make(chan *events_models.Lap, 0)
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	// Start the workers to call the API and generate the lap models
	var workersWg sync.WaitGroup
	for i := 0; i < workers; i++ {
		workersWg.Add(1)
		go parseSessionLapsWorker(irClient,
			tasksChan,
			resultsChan,
			ctx,
			&workersWg,
			cancel,
		)
	}

	// Collect the laps
	laps := make([]*events_models.Lap, 0)

	var outputWg sync.WaitGroup
	outputWg.Add(1)
	go func() {
		defer outputWg.Done()
		for lap := range resultsChan {
			laps = append(laps, lap)
		}
	}()

	// Send the tasks to the workers
	for _, simSessionResult := range results.SessionResults {
		for _, participant := range simSessionResult.Results {
			tasksChan <- sessionLapTask{
				subsessionId:     results.SubsessionId,
				simsessionNumber: simSessionResult.SimsessionNumber,
				custId:           participant.CustId,
			}
		}
	}
	close(tasksChan) // Signal to workers that no more input will be sent

	// Wait for the workers to finish
	workersWg.Wait()
	close(resultsChan) // Signal to the collector that no more output will be sent

	// Wait for the outputs collection to finish
	outputWg.Wait()

	// In case of error, return it
	if err = context.Cause(ctx); err != nil {
		return err
	}

	// DB: create a new transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update the session in the database.
	// If the session is already parsed, return an error.
	// TODO: check if the session is already parsed by the launch date.
	result := tx.Model(&events_models.Session{}).Where("subsession_id = ? AND track_id = 0", subsessionId).Updates(events_models.Session{
		LeagueID: results.LeagueId,
		SeasonID: results.SeasonId,
		LaunchAt: subsessionLaunchAt,
		TrackID:  results.Track.TrackId,
	})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("session %d already parsed", subsessionId)
	}
	if result.Error != nil {
		tx.Rollback()
		return err
	}

	// Store all the simsessions in the database
	sessions := make([]events_models.SessionSimsession, len(results.SessionResults))
	for i, result := range results.SessionResults {
		sessions[i] = events_models.SessionSimsession{
			SubsessionID:     subsessionId,
			SimsessionNumber: result.SimsessionNumber,
			SimsessionType:   result.SimsessionType,
			SimsessionName:   result.SimsessionName,
		}
	}

	if len(sessions) > 0 {
		if err = tx.Create(sessions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store the participants of each simsession in the database
	participants := make([]events_models.SessionSimsessionParticipant, 0)
	for _, result := range results.SessionResults {
		for _, participant := range result.Results {
			participants = append(participants, events_models.SessionSimsessionParticipant{
				SubsessionID:     subsessionId,
				SimsessionNumber: result.SimsessionNumber,
				CustID:           participant.CustId,
				CarID:            participant.CarId,
			})
		}
	}

	if len(participants) > 0 {
		if err = tx.Create(participants).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store the laps
	if len(laps) > 0 {
		if err = tx.Create(laps).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

type sessionLapTask struct {
	subsessionId     int
	simsessionNumber int
	custId           int
}

func parseSessionLapsWorker(irClient *irapi.IRacingApiClient,
	tasksChan <-chan sessionLapTask,
	resultsChan chan<- *events_models.Lap,
	ctx context.Context,
	wg *sync.WaitGroup,
	cancel context.CancelCauseFunc,
) {
	defer wg.Done() // Ensure the wait group counter is decremented when the worker exits

	for {
		select {
		case <-ctx.Done():
			// Another worker has already failed
			return

		case task, ok := <-tasksChan:
			if !ok {
				// The input channel is closed
				return
			}

			res, err := irClient.GetResultsLapData(task.subsessionId, task.simsessionNumber, task.custId)
			if err != nil {
				cancel(fmt.Errorf("error getting lap data for session %d, simsession %d, cust %d: %w", task.subsessionId, task.simsessionNumber, task.custId, err))
				return
			}

			for _, lap := range res.Laps {
				resultsChan <- &events_models.Lap{
					SubsessionID:     task.subsessionId,
					SimsessionNumber: task.simsessionNumber,
					CustID:           lap.CustId,
					LapEvents:        lap.LapEvents,
					Incident:         lap.Incident,
					LapTime:          lap.LapTime,
					LapNumber:        lap.LapNumber,
				}
			}
		}
	}
}
