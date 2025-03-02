package logic

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"riccardotornesello.it/sharedtelemetry/iracing/firestore"
)

func GetFirstSessionLaunchAt(leagueId int, seasonId int, firestoreClient *firestore.Client, firestoreContext context.Context) (string, error) {
	// Extract the first session launchAt for the specified league and season.
	// It will be the oldest session in the database.
	// If no session is found, an empty string is returned.

	query := firestoreClient.Collection("iracing_sessions").
		Where("leagueId", "==", leagueId).
		Where("seasonId", "==", seasonId).
		OrderBy("launchAt", firestore.Asc).
		Limit(1)

	snapshots := query.Documents(firestoreContext)
	defer snapshots.Stop()

	doc, err := snapshots.Next()
	if err == iterator.Done {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	var session firestore_structs.Session
	err = doc.DataTo(&session)
	if err != nil {
		return "", err
	}

	launchAt := session.LaunchAt.Format(time.RFC3339)
	return launchAt, nil
}
