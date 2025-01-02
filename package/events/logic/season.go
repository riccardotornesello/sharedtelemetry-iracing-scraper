package logic

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/events/models"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

type SessionInfo struct {
	SubsessionId int
	LaunchAt     string
}

func GetMissingSessionIds(leagueId int, seasonId int, irClient *irapi.IRacingApiClient, db *gorm.DB) ([]SessionInfo, error) {
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

	var storedSessions []models.Session
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

func SendSessionsToParse(projectId string, topicId string, sessions []SessionInfo) error {
	if sessions == nil || len(sessions) == 0 {
		return nil
	}

	// Send pub/sub messages to parse the sessions
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicId)

	for _, session := range sessions {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("{\"subsessionId\":" + strconv.Itoa(int(session.SubsessionId)) + ",\"launchAt\":\"" + session.LaunchAt + "\"}"),
		})

		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := res.Get(ctx)
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
