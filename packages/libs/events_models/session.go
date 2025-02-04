package events_models

import (
	"time"

	"gorm.io/gorm"
)

// iRacing's session/subsession.
type Session struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SubsessionID int `gorm:"primarykey"`

	LeagueID int
	SeasonID int

	LaunchAt time.Time `gorm:"index"`
	TrackID  int
}
