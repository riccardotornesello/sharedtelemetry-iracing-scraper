package config

import "os"

type Config struct {
	FirestoreProjectID string
	PubSubProjectID    string
	PubSubTopicID      string
}

func Load() Config {
	return Config{
		FirestoreProjectID: os.Getenv("FIRESTORE_PROJECT_ID"),
		PubSubProjectID:    os.Getenv("PUBSUB_PROJECT_ID"),
		PubSubTopicID:      os.Getenv("PUBSUB_TOPIC_ID"),
	}
}
