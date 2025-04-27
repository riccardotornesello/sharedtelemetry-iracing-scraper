package firestore

import "time"

type Session struct {
	LeagueID int       `firestore:"leagueId"`
	SeasonID int       `firestore:"seasonId"`
	LaunchAt time.Time `firestore:"launchAt"`
	TrackID  int       `firestore:"trackId"`

	Simsessions []*SessionSimsession `firestore:"simsessions"`
}

type SessionSimsession struct {
	SimsessionNumber int    `firestore:"simsessionNumber"`
	SimsessionType   int    `firestore:"simsessionType"`
	SimsessionName   string `firestore:"simsessionName"`

	Participants []*SessionSimsessionParticipant `firestore:"participants"`
}

type SessionSimsessionParticipant struct {
	CustID int `firestore:"custId"`
	CarID  int `firestore:"carId"`

	Laps []*Lap `firestore:"laps"`
}

type Lap struct {
	LapEvents []string `firestore:"lapEvents"`
	Incident  bool     `firestore:"incident"`
	LapTime   int      `firestore:"lapTime"`
	LapNumber int      `firestore:"lapNumber"`
}
