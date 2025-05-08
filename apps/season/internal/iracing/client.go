package iracing

import (
	"log"

	irapi "github.com/riccardotornesello/irapi-go"
	"github.com/riccardotornesello/irapi-go/api/league"
	"riccardotornesello.it/sharedtelemetry/iracing/season/config"
)

func FetchData(cfg config.Config, leagueId int, seasonId int) (*league.LeagueSeasonSessionsResponse, error) {
	// Authenticate
	log.Println("Authenticating with iRacing API")
	irClient, err := irapi.NewIRacingApiClient(cfg.IRacingEmail, cfg.IRacingPassword)
	if err != nil {
		return nil, err
	}

	// Get the data
	log.Println("Fetching season")
	resultsOnly := true
	seasonSessions, err := irClient.League.GetLeagueSeasonSessions(league.LeagueSeasonSessionsParams{
		LeagueId:    leagueId,
		SeasonId:    seasonId,
		ResultsOnly: &resultsOnly,
	})
	if err != nil {
		return nil, err
	}

	return seasonSessions, nil
}
