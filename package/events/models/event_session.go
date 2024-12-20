package models

import (
	"time"

	"gorm.io/gorm"
)

// EventSession represents iRacing's session parts, like practice, qualifying and race.
type EventSession struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SubsessionID     int `gorm:"primaryKey; not null"`
	SimsessionNumber int `gorm:"primaryKey; not null"`

	Event Event `gorm:"foreignKey:SubsessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	SimsessionType int
	SimsessionName string
}
