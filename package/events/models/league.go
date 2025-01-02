package models

import (
	"time"

	"gorm.io/gorm"
)

// iRacing's league
// NOTE: if the platform will allow anyone to add any league, the league ID should not be unique
type League struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LeagueID int `gorm:"primarykey"`
}
