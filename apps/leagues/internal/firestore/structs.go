package firestore

type League struct {
	LeagueID int `firestore:"leagueId"`

	Seasons []*Season `firestore:"seasons"`
}

type Season struct {
	SeasonID int `firestore:"seasonId"`
}
