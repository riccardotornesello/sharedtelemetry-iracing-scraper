package events_models

import (
	"time"

	"gorm.io/gorm"
)

type Competition struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LeagueID     int          `gorm:"not null"`
	SeasonID     int          `gorm:"not null"`
	LeagueSeason LeagueSeason `gorm:"foreignKey:LeagueID,SeasonID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`

	Name             string `gorm:"not null"`
	Slug             string `gorm:"not null;unique"`
	CrewDriversCount int    `gorm:"not null;default:1"`
}
