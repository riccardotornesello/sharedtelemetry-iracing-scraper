package logic

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firestore_structs "riccardotornesello.it/sharedtelemetry/iracing/firestore"
)

func StoreCars(cars map[string]firestore_structs.Car, firestoreClient *firestore.Client, firestoreContext context.Context) error {
	db := firestoreClient.Collection(firestore_structs.CarsCollection)

	for carId, car := range cars {
		_, err := db.Doc(carId).Set(firestoreContext, car)
		if err != nil {
			return fmt.Errorf("error updating car %s in the database: %w", carId, err)
		}
	}

	return nil
}

func StoreCarClasses(carClasses map[string]firestore_structs.CarClass, firestoreClient *firestore.Client, firestoreContext context.Context) error {
	db := firestoreClient.Collection(firestore_structs.CarClassesCollection)

	for carClassId, carClass := range carClasses {
		_, err := db.Doc(carClassId).Set(firestoreContext, carClass)
		if err != nil {
			return fmt.Errorf("error updating car class %s in the database: %w", carClassId, err)
		}
	}

	return nil
}
