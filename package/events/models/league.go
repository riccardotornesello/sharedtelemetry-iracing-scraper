package models

import (
	"time"

	"gorm.io/gorm"
)

type League struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LeagueID int `gorm:"primarykey"`
}
