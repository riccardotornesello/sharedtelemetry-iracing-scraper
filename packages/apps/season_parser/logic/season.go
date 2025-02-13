package logic

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"riccardotornesello.it/sharedtelemetry/iracing/events_models"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

type SessionInfo struct {
	SubsessionId int
	LaunchAt     string
}

func GetMissingSessionInfo(leagueId int, seasonId int, irClient *irapi.IRacingApiClient, db *gorm.DB) ([]SessionInfo, error) {
	// Extract the sessions list (only the completed ones) for the specified series and league
	sessions, err := irClient.GetLeagueSeasonSessions(leagueId, seasonId, true)
	if err != nil {
		return nil, err
	}

	// Get the sessions which are not already stored in the database
	sessionIds := make([]int, len(sessions.Sessions))
	for i, session := range sessions.Sessions {
		sessionIds[i] = session.SubsessionId
	}

	var storedSessions []events_models.Session
	db.Where("subsession_id IN ?", sessionIds).Find(&storedSessions)
	storedSessionIds := make(map[int]bool)
	for _, storedSession := range storedSessions {
		storedSessionIds[int(storedSession.SubsessionID)] = true
	}

	// Return the sessions which are not already stored in the database
	missingSessions := make([]SessionInfo, 0)
	for _, session := range sessions.Sessions {
		if _, ok := storedSessionIds[int(session.SubsessionId)]; !ok {
			missingSessions = append(missingSessions, SessionInfo{
				SubsessionId: session.SubsessionId,
				LaunchAt:     session.LaunchAt,
			})
		}
	}

	return missingSessions, nil
}

func SendSessionsToParse(pubSubTopic *pubsub.Topic, pubSubCtx context.Context, sessions []SessionInfo) error {
	if sessions == nil || len(sessions) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	var totalErrors uint64

	for _, session := range sessions {
		result := pubSubTopic.Publish(pubSubCtx, &pubsub.Message{
			Data: []byte("{\"subsessionId\":" + strconv.Itoa(int(session.SubsessionId)) + ",\"launchAt\":\"" + session.LaunchAt + "\"}"),
		})

		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := res.Get(pubSubCtx)
			if err != nil {
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, len(sessions))
	}

	return nil
}

func StoreMissingSessions(sessions []SessionInfo, db *gorm.DB) error {
	if sessions == nil || len(sessions) == 0 {
		return nil
	}

	dbSessions := make([]events_models.Session, len(sessions))

	for i, session := range sessions {
		dbSessions[i] = events_models.Session{
			SubsessionID: session.SubsessionId,
		}
	}

	// Create and ignore duplicates on SubsessionID
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "subsession_id"}},
		DoNothing: true,
	}).Create(dbSessions).Error

	return err
}
