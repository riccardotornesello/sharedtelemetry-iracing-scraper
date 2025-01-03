package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Lap struct {
	gorm.Model

	SubsessionID     int `gorm:"not null"`
	SimsessionNumber int `gorm:"not null"`
	CustID           int `gorm:"not null"`

	SessionSimsessionParticipant SessionSimsessionParticipant `gorm:"foreignKey:SubsessionID,SimsessionNumber,CustID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LapEvents pq.StringArray `gorm:"type:text[]"`
	Incident  bool
	LapTime   int
	LapNumber int
}
