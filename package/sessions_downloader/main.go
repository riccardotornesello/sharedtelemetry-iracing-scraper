package sessions_downloader

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"riccardotornesello.it/iracing-average-lap/client"
	"riccardotornesello.it/iracing-average-lap/database"
	"riccardotornesello.it/iracing-average-lap/models"
)

func Run() {
	////////////////////////////////////////
	// INITIALIZATION
	////////////////////////////////////////

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	seriesId, err := strconv.Atoi(os.Getenv("IRACING_SERIES_ID"))
	if err != nil {
		log.Fatal("Invalid series id")
	}

	leagueId, err := strconv.Atoi(os.Getenv("IRACING_LEAGUE_ID"))
	if err != nil {
		log.Fatal("Invalid league id")
	}

	saveRequests := os.Getenv("SAVE_REQUESTS") == "true"

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Error connecting to the database")
	}

	// Initialize iRacing API client
	irClient, err := client.NewIRacingApiClient(os.Getenv("IRACING_EMAIL"), os.Getenv("IRACING_PASSWORD"))
	if err != nil {
		log.Fatal("Error initializing iRacing API client")
	}

	if saveRequests {
		// Initialize the directory structure
		// Create (if not exists) the downloads folder. In that folder, create the folders sessions and laps
		os.Mkdir("downloads", 0755)
		os.Mkdir("downloads/sessions", 0755)
		os.Mkdir("downloads/laps", 0755)
	}

	////////////////////////////////////////
	// PROCESSING
	////////////////////////////////////////

	// Extract the sessions list (only the completed ones) for the specified series and league
	sessions, err := irClient.GetLeagueSeasonSessions(seriesId, leagueId, true)
	if err != nil {
		log.Fatal("Error getting league season sessions")
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
		storedSessionIds[int(storedSession.SubsessionId)] = true
	}

	// Start 3 workers to get the lap results for each driver
	numWorkers := 3
	maxNumJobs := len(sessions.Sessions)
	sessionJobs := make(chan *client.LeagueSeasonSession, maxNumJobs)
	sessionJobResults := make(chan interface{}, maxNumJobs)
	for w := 0; w < numWorkers; w++ {
		go sessionWorker(irClient, sessionJobs, sessionJobResults, db, saveRequests)
	}

	numJobs := 0
	for _, session := range sessions.Sessions {
		// Check if the session is in the specified days
		// TODO: remove this and store all the sessions
		startDate := session.LaunchAt[:10]
		if startDate != "2024-09-11" && startDate != "2024-09-14" {
			continue
		}

		// Check if the session is already stored
		if _, ok := storedSessionIds[int(session.SubsessionId)]; ok {
			continue
		}

		numJobs++
		sessionJobs <- &session
	}
	close(sessionJobs)

	for a := 0; a < numJobs; a++ {
		<-sessionJobResults
	}
	close(sessionJobResults)
}

func sessionWorker(irClient *client.IRacingApiClient, sessionJobs <-chan *client.LeagueSeasonSession, sessionJobResults chan<- interface{}, db *gorm.DB, saveRequests bool) {
	for session := range sessionJobs {
		err := parseSession(irClient, session, db, saveRequests)
		if err != nil {
			log.Println("Error parsing session", session.SubsessionId)
			log.Println(err)
		}

		sessionJobResults <- interface{}(nil)
	}
}
