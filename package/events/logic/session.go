package logic

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/events/models"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

func ParseSession(irClient *irapi.IRacingApiClient, subsessionId int, subsessionLaunchAt time.Time, db *gorm.DB) error {
	// Skip if the session is already in the database
	var count int64
	db.Model(&models.Session{}).Where("subsession_id = ?", subsessionId).Count(&count)
	if count > 0 {
		slog.Info("Session", subsessionId, "already parsed")
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
	// NOTE: the laps download can be parallelized but the API has a rate limit and we are already parsing other sessions in parallel
	laps := make([]models.Lap, 0)
	for _, simSessionResult := range results.SessionResults {
		for _, participant := range simSessionResult.Results {
			res, err := irClient.GetResultsLapData(results.SubsessionId, simSessionResult.SimsessionNumber, participant.CustId)
			if err != nil {
				return fmt.Errorf("error getting lap data for session %d, simsession %d, cust %d: %w", results.SubsessionId, simSessionResult.SimsessionNumber, participant.CustId, err)
			}

			for _, lap := range res.Laps {
				laps = append(laps, models.Lap{
					SubsessionID:     results.SubsessionId,
					SimsessionNumber: simSessionResult.SimsessionNumber,
					CustID:           lap.CustId,
					LapEvents:        lap.LapEvents,
					Incident:         lap.Incident,
					LapTime:          lap.LapTime,
					LapNumber:        lap.LapNumber,
				})
			}
		}
	}

	// DB: create a new transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Store the session in the database
	if err = tx.Create(&models.Session{
		SubsessionID: subsessionId,
		LeagueID:     results.LeagueId,
		SeasonID:     results.SeasonId,
		LaunchAt:     subsessionLaunchAt,
		TrackID:      results.Track.TrackId,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Store all the simsessions in the database
	sessions := make([]models.SessionSimsession, len(results.SessionResults))
	for i, result := range results.SessionResults {
		sessions[i] = models.SessionSimsession{
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
	participants := make([]models.SessionSimsessionParticipant, 0)
	for _, result := range results.SessionResults {
		for _, participant := range result.Results {
			participants = append(participants, models.SessionSimsessionParticipant{
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
