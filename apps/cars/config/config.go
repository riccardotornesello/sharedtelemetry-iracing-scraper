package config

import "os"

type Config struct {
	FirestoreProjectID string
	IRacingEmail       string
	IRacingPassword    string
}

func Load() Config {
	return Config{
		FirestoreProjectID: os.Getenv("FIRESTORE_PROJECT_ID"),
		IRacingEmail:       os.Getenv("IRACING_EMAIL"),
		IRacingPassword:    os.Getenv("IRACING_PASSWORD"),
	}
}
