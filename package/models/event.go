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

	SubsessionId int `gorm:"primarykey"`

	LaunchAt time.Time `gorm:"index"`
	TrackId  int

	// TODO: league and season
}
