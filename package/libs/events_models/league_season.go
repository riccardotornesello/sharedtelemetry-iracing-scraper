package events_models

import (
	"time"

	"gorm.io/gorm"
)

// iRacing's season
type LeagueSeason struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LeagueID int `gorm:"primarykey"`
	SeasonID int `gorm:"primarykey"`

	League League `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
