package events_models

import (
	"time"

	"gorm.io/gorm"
)

// iRacing's session parts, like practice, qualifying and race.
type SessionSimsession struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SubsessionID     int `gorm:"primaryKey; not null"`
	SimsessionNumber int `gorm:"primaryKey; not null"`

	Session Session `gorm:"foreignKey:SubsessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	SimsessionType int
	SimsessionName string
}
