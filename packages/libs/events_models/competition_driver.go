package events_models

import (
	"time"

	"gorm.io/gorm"
)

type CompetitionDriver struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CrewID uint            `gorm:"not null"`
	Crew   CompetitionCrew `gorm:"foreignKey:CrewID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	IRacingCustId int    `gorm:"not null"`
	FirstName     string `gorm:"not null"`
	LastName      string
}
