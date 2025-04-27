package processor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/riccardotornesello/irapi-go/api/results"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/internal/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/internal/iracing"
)

const MAX_WORKERS = 10

type sessionLapTask struct {
	subsessionId     int
	simsessionNumber int
	custId           int
}

type workerResponse struct {
	simsessionNumber int
	custId           int
	laps             []*firestore.Lap
}

func Process(irClient *iracing.Client, sessionResponse *results.ResultsGetResponse, launchAtStr string) (*firestore.Session, error) {
	log.Println("Analyzing data")

	// --- Parse the date object
	launchAt, err := time.Parse(time.RFC3339, launchAtStr)
	if err != nil {
		return nil, err
	}

	// --- Create the data object
	session := firestore.Session{
		LeagueID: sessionResponse.LeagueId,
		SeasonID: sessionResponse.SeasonId,
		LaunchAt: launchAt,
		TrackID:  sessionResponse.Track.TrackId,

		Simsessions: make([]*firestore.SessionSimsession, len(sessionResponse.SessionResults)),
	}

	// --- Count the number of tasks to be done (one for each driver in each simsession)
	tasksCount := 0
	for _, simSessionResult := range sessionResponse.SessionResults {
		tasksCount += len(simSessionResult.Results)
	}

	// --- Start the workers to call the API and generate the lap models
	var workersWg sync.WaitGroup

	tasksChan := make(chan sessionLapTask, tasksCount)
	resultsChan := make(chan *workerResponse)

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	// Limit the number of workers to the number of tasks if less than WORKERS
	workersCount := MAX_WORKERS
	if tasksCount < workersCount {
		workersCount = tasksCount
	}

	for i := 0; i < workersCount; i++ {
		workersWg.Add(1)
		go parseSessionLapsWorker(irClient,
			tasksChan,
			resultsChan,
			ctx,
			&workersWg,
			cancel,
		)
	}

	// --- Start a goroutine to collect the laps in the background
	lapResults := make([]*workerResponse, 0)

	var outputWg sync.WaitGroup
	outputWg.Add(1)

	go func() {
		defer outputWg.Done()
		for result := range resultsChan {
			lapResults = append(lapResults, result)
		}
	}()

	// -- Send the tasks to the workers
	for _, simSessionResult := range sessionResponse.SessionResults {
		for _, participant := range simSessionResult.Results {
			tasksChan <- sessionLapTask{
				subsessionId:     sessionResponse.SubsessionId,
				simsessionNumber: simSessionResult.SimsessionNumber,
				custId:           participant.CustId,
			}
		}
	}
	close(tasksChan) // Signal to workers that no more input will be sent

	// --- Wait for the workers to finish
	workersWg.Wait()
	close(resultsChan) // Signal to the collector that no more output will be sent

	// --- Wait for the outputs collection to finish
	outputWg.Wait()

	// --- In case of error in the workers, return it
	if err = context.Cause(ctx); err != nil {
		return nil, err
	}

	// --- Create the simsessions and participants maps
	simsessions := make(map[int]*firestore.SessionSimsession)
	for i, result := range sessionResponse.SessionResults {
		simsession := firestore.SessionSimsession{
			SimsessionNumber: result.SimsessionNumber,
			SimsessionType:   result.SimsessionType,
			SimsessionName:   result.SimsessionName,

			Participants: make([]*firestore.SessionSimsessionParticipant, len(result.Results)),
		}

		simsessions[result.SimsessionNumber] = &simsession
		session.Simsessions[i] = &simsession
	}

	simsessionParticipants := make(map[int]map[int]*firestore.SessionSimsessionParticipant)
	for _, result := range sessionResponse.SessionResults {
		simsessionParticipants[result.SimsessionNumber] = make(map[int]*firestore.SessionSimsessionParticipant)

		for i, participantResults := range result.Results {
			participant := firestore.SessionSimsessionParticipant{
				CustID: participantResults.CustId,
				CarID:  participantResults.CarId,
			}

			simsessionParticipants[result.SimsessionNumber][participantResults.CustId] = &participant
			simsessions[result.SimsessionNumber].Participants[i] = &participant
		}
	}

	// --- Populate the laps in the participants
	for _, lapResult := range lapResults {
		simsessionParticipants[lapResult.simsessionNumber][lapResult.custId].Laps = lapResult.laps
	}

	return &session, nil
}

func parseSessionLapsWorker(
	irClient *iracing.Client,
	tasksChan <-chan sessionLapTask,
	resultsChan chan<- *workerResponse,
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

			res, err := irClient.FetchDriverResults(task.subsessionId, task.simsessionNumber, task.custId)
			if err != nil {
				cancel(fmt.Errorf("error getting lap data for session %d, simsession %d, cust %d: %w", task.subsessionId, task.simsessionNumber, task.custId, err))
				return
			}

			laps := make([]*firestore.Lap, len(res.Chunks))
			for i, lap := range res.Chunks {
				laps[i] = &firestore.Lap{
					LapEvents: lap.LapEvents,
					Incident:  lap.Incident,
					LapTime:   lap.LapTime,
					LapNumber: lap.LapNumber,
				}
			}

			resultsChan <- &workerResponse{
				simsessionNumber: task.simsessionNumber,
				custId:           task.custId,
				laps:             laps,
			}
		}
	}
}
