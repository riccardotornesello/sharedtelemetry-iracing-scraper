package logic

import (
	"time"

	"riccardotornesello.it/sharedtelemetry/iracing/irapi"
)

type SessionInfo struct {
	SubsessionId int
	LaunchAt     string
}

func GetLeagueSeasonSessionsInfo(leagueId int, seasonId int, maxLaunchatStr string, irClient *irapi.IRacingApiClient) ([]SessionInfo, error) {
	// Extract the sessions list (only the completed ones) for the specified series and league
	sessions, err := irClient.GetLeagueSeasonSessions(leagueId, seasonId, true)
	if err != nil {
		return nil, err
	}

	// Parse the launchAt date and time
	if maxLaunchatStr == "" {
		maxLaunchatStr = "0001-01-01T00:00:00Z"
	}
	maxLaunchAt, err := time.Parse(time.RFC3339, maxLaunchatStr)
	if err != nil {
		return nil, err
	}

	// Extract the session info
	sessionsInfo := make([]SessionInfo, 0)
	for _, session := range sessions.Sessions {
		sessionLaunchAt, err := time.Parse(time.RFC3339, session.LaunchAt)
		if err != nil {
			return nil, err
		}

		if sessionLaunchAt.After(maxLaunchAt) {
			continue
		}

		sessionsInfo = append(sessionsInfo, SessionInfo{
			SubsessionId: session.SubsessionId,
			LaunchAt:     session.LaunchAt,
		})
	}

	return sessionsInfo, nil
}
