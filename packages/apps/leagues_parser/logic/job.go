package logic

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"gorm.io/gorm"
)

func FeedLeagues(db *gorm.DB, pubSubTopic *pubsub.Topic, pubSubCtx context.Context) error {
	// Get active league seasons
	seasonInfos, err := GetActiveLeagueSeasonIds(db)
	if err != nil {
		return err
	}

	// Send pub/sub messages to parse the season
	var wg sync.WaitGroup
	var totalErrors uint64

	for _, season := range seasonInfos {
		result := pubSubTopic.Publish(pubSubCtx, &pubsub.Message{
			Data: []byte("{\"leagueId\":" + strconv.Itoa(season.LeagueId) + ",\"seasonId\":" + strconv.Itoa(season.SeasonId) + "}"),
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
		return fmt.Errorf("Failed to send %d pub/sub messages", totalErrors)
	}

	return nil
}
