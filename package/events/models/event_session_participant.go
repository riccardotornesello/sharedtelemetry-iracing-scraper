package models

import (
	"gorm.io/gorm"
)

// EventSessionParticipant represents a participant in an iRacing event.
type EventSessionParticipant struct {
	gorm.Model

	SubsessionID     int `gorm:"primaryKey; not null"`
	SimsessionNumber int `gorm:"primaryKey; not null"`
	CustID           int `gorm:"primaryKey; not null"`

	EventSession EventSession `gorm:"foreignKey:SubsessionID,SimsessionNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CarID int
}
