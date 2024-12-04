package models

import (
	"gorm.io/gorm"
)

type Lap struct {
	gorm.Model

	EventSessionID uint
	EventSession   EventSession

	CustId int

	LapEvents string
	Incident  bool
	LapTime   int
	LapNumber int

	// TODO: car
}
