package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents iRacing's session.
type Event struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SubsessionID int `gorm:"primarykey"`

	LeagueID int
	SeasonID int

	LaunchAt time.Time `gorm:"index"`
	TrackID  int
}
