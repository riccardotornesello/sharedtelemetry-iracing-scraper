package logic

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
	"riccardotornesello.it/sharedtelemetry/iracing/models"
)

func ParseSession(irClient *irapi.IRacingApiClient, subsessionId int, subsessionLaunchAt time.Time, db *gorm.DB) error {
	log.Println("Parsing session", subsessionId)

	// Skip if the session is already in the database
	var count int64
	db.Model(&models.Event{}).Where("subsession_id = ?", subsessionId).Count(&count)
	if count > 0 {
		log.Println("Session", subsessionId, "already parsed")
		return nil
	}

	// Get the whole session results
	results, err := irClient.GetResults(subsessionId)
	if err != nil {
		return fmt.Errorf("error getting results for session %d: %w", subsessionId, err)
	}

	// For each subsession, get the results for each driver
	// results.SessionResults: one for each session (practice, quali...)
	// results.SessionResults[i].Results: one for each driver
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

	// Store the event in the database
	if err = tx.Create(&models.Event{
		SubsessionID: subsessionId,
		LeagueID:     results.LeagueId,
		SeasonID:     results.SeasonId,
		LaunchAt:     subsessionLaunchAt,
		TrackID:      results.Track.TrackId,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Store all the sessions in the database
	sessions := make([]models.EventSession, len(results.SessionResults))
	for i, result := range results.SessionResults {
		sessions[i] = models.EventSession{
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

	// Store the participants of each session in the database
	participants := make([]models.EventSessionParticipant, 0)
	for _, result := range results.SessionResults {
		for _, participant := range result.Results {
			participants = append(participants, models.EventSessionParticipant{
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

	log.Println("Session", subsessionId, "parsed")

	return tx.Commit().Error
}
