package models

import (
	"gorm.io/gorm"
)

// A participant in an iRacing event.
type SessionSimsessionParticipant struct {
	gorm.Model

	SubsessionID     int `gorm:"primaryKey; not null"`
	SimsessionNumber int `gorm:"primaryKey; not null"`
	CustID           int `gorm:"primaryKey; not null"`

	SessionSimsession SessionSimsession `gorm:"foreignKey:SubsessionID,SimsessionNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CarID int
}
