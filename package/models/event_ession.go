package models

import (
	"gorm.io/gorm"
)

// EventSession represents iRacing's session parts, like practice, qualifying and race.
type EventSession struct {
	gorm.Model

	EventID int
	Event   Event

	SimsessionNumber int
	SimsessionType   int
	SimsessionName   string
}
