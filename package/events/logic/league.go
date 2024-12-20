package logic

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/events/models"
	irapi "riccardotornesello.it/sharedtelemetry/iracing/iracing-api"
)

func ParseLeague(leagueId int, seasonId int, irClient *irapi.IRacingApiClient, db *gorm.DB) error {
	log.Println("Parsing league", leagueId, seasonId)

	projectId := os.Getenv("PROJECT_ID")
	topicId := os.Getenv("TOPIC_ID")

	// Extract the sessions list (only the completed ones) for the specified series and league
	sessions, err := irClient.GetLeagueSeasonSessions(leagueId, seasonId, true)
	if err != nil {
		return err
	}

	// Get the sessions which are not already stored in the database
	sessionIds := make([]int, len(sessions.Sessions))
	for i, session := range sessions.Sessions {
		sessionIds[i] = session.SubsessionId
	}

	var storedSessions []models.Event
	db.Where("subsession_id IN ?", sessionIds).Find(&storedSessions)
	storedSessionIds := make(map[int]bool)
	for _, storedSession := range storedSessions {
		storedSessionIds[int(storedSession.SubsessionID)] = true
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

	for _, session := range sessions.Sessions {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("{\"subsessionId\":" + strconv.Itoa(int(session.SubsessionId)) + ",\"launchAt\":\"" + session.LaunchAt + "\"}"),
		})

		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			_, err := res.Get(ctx)
			if err != nil {
				log.Printf("Failed to publish: %v", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, len(sessions.Sessions))
	}

	log.Println("League parsed", leagueId, seasonId, "with", len(sessions.Sessions), "sessions (", len(sessions.Sessions)-len(storedSessions), "new sessions )")

	return nil

	// // Start 3 workers to get the lap results for each driver
	// numWorkers := 3
	// maxNumJobs := len(sessions.Sessions)
	// sessionJobs := make(chan *irapi.LeagueSeasonSession, maxNumJobs)
	// sessionJobResults := make(chan interface{}, maxNumJobs)
	// for w := 0; w < numWorkers; w++ {
	// 	go sessionWorker(irClient, sessionJobs, sessionJobResults, db, saveRequests)
	// }

	// numJobs := 0
	// for _, session := range sessions.Sessions {
	// 	// Check if the session is in the specified days
	// 	// TODO: remove this and store all the sessions
	// 	startDate := session.LaunchAt[:10]
	// 	if startDate != "2024-09-11" && startDate != "2024-09-14" {
	// 		continue
	// 	}

	// 	// Check if the session is already stored
	// 	if _, ok := storedSessionIds[int(session.SubsessionId)]; ok {
	// 		continue
	// 	}

	// 	numJobs++
	// 	sessionJobs <- &session
	// }
	// close(sessionJobs)

	// for a := 0; a < numJobs; a++ {
	// 	<-sessionJobResults
	// }
	// close(sessionJobResults)
}
