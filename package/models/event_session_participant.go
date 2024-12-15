package models

import (
	"gorm.io/gorm"
)

// EventSessionParticipant represents a participant in an iRacing event.
type EventSessionParticipant struct {
	gorm.Model

	EventSessionID uint
	EventSession   EventSession `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CustId int
	CarId  int
}
