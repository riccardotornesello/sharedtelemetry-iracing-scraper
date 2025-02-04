package logic

import (
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/events_models"
)

type SeasonInfo struct {
	LeagueId int
	SeasonId int
}

func GetActiveLeagueSeasonIds(db *gorm.DB) ([]SeasonInfo, error) {
	var leagueSeasons []events_models.LeagueSeason
	err := db.Model(&events_models.LeagueSeason{}).Find(&leagueSeasons).Error
	if err != nil {
		return nil, err
	}

	seasonInfos := make([]SeasonInfo, len(leagueSeasons))
	for i, leagueSeason := range leagueSeasons {
		seasonInfos[i] = SeasonInfo{
			LeagueId: leagueSeason.LeagueID,
			SeasonId: leagueSeason.SeasonID,
		}
	}

	return seasonInfos, nil
}
