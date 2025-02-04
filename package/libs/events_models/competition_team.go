package events_models

import (
	"time"

	"gorm.io/gorm"
)

type CompetitionTeam struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CompetitionID uint        `gorm:"not null"`
	Competition   Competition `gorm:"foreignKey:CompetitionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Name string `gorm:"not null"`
}
