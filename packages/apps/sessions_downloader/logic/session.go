package logic

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	firestore_structs "riccardotornesello.it/sharedtelemetry/iracing/firestore"
	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

type workerResponse struct {
	simsessionNumber int
	custId           int
	laps             []*firestore_structs.Lap
}

func ParseSession(irClient *irapi.IRacingApiClient, subsessionId int, subsessionLaunchAt time.Time, firestoreClient *firestore.Client, firestoreContext context.Context, workers int) error {
	db := firestoreClient.Collection("iracing_sessions")

	// Skip if already in the database
	dbSession := firestore_structs.Session{}
	dbSessionDoc := db.Doc(fmt.Sprintf("%d", subsessionId))
	dbSessionSnap, err := dbSessionDoc.Get(firestoreContext)
	if err != nil {
		// TODO: handle
	} else {
		err = dbSessionSnap.DataTo(&dbSession)
		if err != nil {
			return fmt.Errorf("error parsing session %d from the database: %w", subsessionId, err)
		}
	}

	if dbSession.Parsed {
		log.Printf("Session %d already parsed", subsessionId)
		return nil
	}

	// Get the whole session results to extract simsessions and participants
	results, err := irClient.GetResults(subsessionId)
	if err != nil {
		return fmt.Errorf("error getting results for session %d: %w", subsessionId, err)
	}

	session := firestore_structs.Session{
		Parsed: true,

		LeagueID: results.LeagueId,
		SeasonID: results.SeasonId,
		LaunchAt: subsessionLaunchAt, // TODO: populate in league parser
		TrackID:  results.Track.TrackId,

		Simsessions: make([]*firestore_structs.SessionSimsession, len(results.SessionResults)),
	}

	// For each simsession, get the results for each driver.
	// results.SessionResults: one for each simsession (practice, quali...)
	// results.SessionResults[i].Results: one for each driver

	// Count the number of tasks to be done (one for each driver in each simsession)
	tasksCount := 0
	for _, simSessionResult := range results.SessionResults {
		tasksCount += len(simSessionResult.Results)
	}

	tasksChan := make(chan sessionLapTask, tasksCount)
	resultsChan := make(chan *workerResponse, 0)
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

	// Collect the laps in background
	lapResults := make([]*workerResponse, 0)
	var outputWg sync.WaitGroup
	outputWg.Add(1)
	go func() {
		defer outputWg.Done()
		for result := range resultsChan {
			lapResults = append(lapResults, result)
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

	// In case of error in the workers, return it
	if err = context.Cause(ctx); err != nil {
		return err
	}

	// Create the simsessions and participants maps
	simsessions := make(map[int]*firestore_structs.SessionSimsession)
	for i, result := range results.SessionResults {
		simsession := firestore_structs.SessionSimsession{
			SimsessionNumber: result.SimsessionNumber,
			SimsessionType:   result.SimsessionType,
			SimsessionName:   result.SimsessionName,

			Participants: make([]*firestore_structs.SessionSimsessionParticipant, len(result.Results)),
		}

		simsessions[result.SimsessionNumber] = &simsession
		session.Simsessions[i] = &simsession
	}

	simsessionParticipants := make(map[int]map[int]*firestore_structs.SessionSimsessionParticipant)
	for _, result := range results.SessionResults {
		simsessionParticipants[result.SimsessionNumber] = make(map[int]*firestore_structs.SessionSimsessionParticipant)

		for i, participantResults := range result.Results {
			participant := firestore_structs.SessionSimsessionParticipant{
				CustID: participantResults.CustId,
				CarID:  participantResults.CarId,
			}

			simsessionParticipants[result.SimsessionNumber][participantResults.CustId] = &participant
			simsessions[result.SimsessionNumber].Participants[i] = &participant
		}
	}

	// Populate the laps in the participants
	for _, lapResult := range lapResults {
		simsessionParticipants[lapResult.simsessionNumber][lapResult.custId].Laps = lapResult.laps
	}

	// Save the session in the database
	_, err = db.Doc(fmt.Sprintf("%d", subsessionId)).Set(firestoreContext, session)
	if err != nil {
		return fmt.Errorf("error updating session %d in the database: %w", subsessionId, err)
	}

	return nil
}

type sessionLapTask struct {
	subsessionId     int
	simsessionNumber int
	custId           int
}

func parseSessionLapsWorker(irClient *irapi.IRacingApiClient,
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

			res, err := irClient.GetResultsLapData(task.subsessionId, task.simsessionNumber, task.custId)
			if err != nil {
				cancel(fmt.Errorf("error getting lap data for session %d, simsession %d, cust %d: %w", task.subsessionId, task.simsessionNumber, task.custId, err))
				return
			}

			laps := make([]*firestore_structs.Lap, len(res.Laps))
			for i, lap := range res.Laps {
				laps[i] = &firestore_structs.Lap{
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
