package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Lap struct {
	gorm.Model

	SubsessionID     int
	SimsessionNumber int
	EventSession     EventSession `gorm:"foreignKey:SubsessionID,SimsessionNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CustID int

	LapEvents pq.StringArray `gorm:"type:text[]"`
	Incident  bool
	LapTime   int
	LapNumber int
}
