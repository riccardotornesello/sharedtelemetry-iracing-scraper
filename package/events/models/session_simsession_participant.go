package models

import (
	"time"

	"gorm.io/gorm"
)

// A participant in an iRacing event.
type SessionSimsessionParticipant struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SubsessionID     int `gorm:"primaryKey; not null"`
	SimsessionNumber int `gorm:"primaryKey; not null"`
	CustID           int `gorm:"primaryKey; not null"`

	SessionSimsession SessionSimsession `gorm:"foreignKey:SubsessionID,SimsessionNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CarID int
}
