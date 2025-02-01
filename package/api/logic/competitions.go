package logic

import (
	"time"

	"gorm.io/gorm"
	models "riccardotornesello.it/sharedtelemetry/iracing/db/events_models"
)

type CompetitionSession struct {
	EventGroupId     uint
	Date             string
	LaunchAt         time.Time
	SubsessionId     int
	SimsessionNumber int
}

func GetCompetitionSessions(db *gorm.DB, competitionId int) ([]*CompetitionSession, map[int]*CompetitionSession, error) {
	var sessions []*CompetitionSession
	err := db.
		Table("session_simsessions").
		Select("event_groups.id as event_group_id, text(date(sessions.launch_at)) as date, sessions.launch_at, session_simsessions.subsession_id, session_simsessions.simsession_number").
		Joins("join sessions on session_simsessions.subsession_id = sessions.subsession_id").
		Joins("join event_groups on sessions.track_id = event_groups.i_racing_track_id and text(date(sessions.launch_at)) = ANY(event_groups.dates)").
		Joins("join competitions on competitions.id = event_groups.competition_id").
		Where("event_groups.competition_id = ?", competitionId).
		Where("session_simsessions.simsession_name = ?", "QUALIFY").
		Where("sessions.league_id = competitions.league_id").
		Where("sessions.season_id = competitions.season_id").
		Order("event_groups.id, sessions.launch_at").
		Find(&sessions).
		Error
	if err != nil {
		return nil, nil, err
	}

	sessionsMap := make(map[int]*CompetitionSession)
	for _, session := range sessions {
		sessionsMap[session.SubsessionId] = session
	}

	return sessions, sessionsMap, nil
}

func GetEventGroupSessions(db *gorm.DB, trackId int, sessionDate string, leagueId int, seasonId int) ([]*models.SessionSimsession, error) {
	simsessionName := "QUALIFY"
	simsessionNumber := 0

	var simsessions []*models.SessionSimsession
	err := db.
		Joins("Session").
		Where("session_simsessions.simsession_name = ?", simsessionName).
		Where("session_simsessions.simsession_number = ?", simsessionNumber).
		Where("\"Session\".track_id = ?", trackId).
		Where("date(\"Session\".launch_at) = ?", sessionDate).
		Where("\"Session\".league_id = ?", leagueId).
		Where("\"Session\".season_id = ?", seasonId).
		Find(&simsessions).
		Error
	if err != nil {
		return nil, err
	}

	return simsessions, nil
}

func GetCompetitionDrivers(db *gorm.DB, competitionId int) (map[int]*models.CompetitionDriver, error) {
	var competitionDrivers []*models.CompetitionDriver
	err := db.
		Joins("Crew").
		Joins("Crew.Team").
		Where("\"Crew__Team\".competition_id = ?", competitionId).
		Order("name").
		Find(&competitionDrivers).
		Error
	if err != nil {
		return nil, err
	}

	competitionDriversMap := make(map[int]*models.CompetitionDriver)
	for _, driver := range competitionDrivers {
		competitionDriversMap[driver.IRacingCustId] = driver
	}

	return competitionDriversMap, nil
}

func GetLaps(db *gorm.DB, simsessionIds [][]int) ([]*models.Lap, error) {
	var laps []*models.Lap
	err := db.
		Joins("SessionSimsessionParticipant").
		Where("(laps.subsession_id, laps.simsession_number) IN ?", simsessionIds).
		Order("laps.cust_id, laps.subsession_id, laps.simsession_number, laps.lap_number").
		Find(&laps).
		Error
	if err != nil {
		return nil, err
	}

	return laps, nil
}

func GetEventGroups(db *gorm.DB, competitionId int) ([]*models.EventGroup, error) {
	var groups []*models.EventGroup
	err := db.
		Where("competition_id = ?", competitionId).
		Find(&groups).
		Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func GetSessionParticipants(db *gorm.DB, subsessionId int, simsessionNumber int) ([]*models.SessionSimsessionParticipant, error) {
	var participants []*models.SessionSimsessionParticipant
	err := db.
		Where("subsession_id = ?", subsessionId).
		Where("simsesion_number = ?", simsessionNumber).
		Find(&participants).
		Error
	if err != nil {
		return nil, err
	}

	return participants, nil
}
