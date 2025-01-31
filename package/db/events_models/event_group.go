package events_models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type EventGroup struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CompetitionID uint        `gorm:"not null"`
	Competition   Competition `gorm:"foreignKey:CompetitionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Name           string         `gorm:"not null"`
	IRacingTrackId int            `gorm:"not null"`
	Dates          pq.StringArray `gorm:"type:text[]"`
}
