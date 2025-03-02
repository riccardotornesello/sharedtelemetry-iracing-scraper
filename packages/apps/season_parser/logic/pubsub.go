package logic

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

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
