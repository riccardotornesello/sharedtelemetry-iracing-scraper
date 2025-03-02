package logic

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	firestore_structs "riccardotornesello.it/sharedtelemetry/iracing/firestore"
)

func GetActiveLeagues(firestoreClient *firestore.Client, firestoreContext context.Context) ([]firestore_structs.League, error) {
	var leagues []firestore_structs.League
	iter := firestoreClient.Collection(firestore_structs.LeaguesCollection).Documents(firestoreContext)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		var league firestore_structs.League
		doc.DataTo(&league)
		leagues = append(leagues, league)
	}
	return leagues, nil
}
