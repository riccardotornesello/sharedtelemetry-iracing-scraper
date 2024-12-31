package models

import (
	"time"

	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/common/database"
)

type carCategory string

type DriverStats struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CustID      int           `gorm:"primaryKey;not null"`
	CarCategory string        `gorm:"primaryKey;not null"`
	Date        database.Date `gorm:"primaryKey;not null"`

	License string `gorm:"not null"`
	IRating int    `gorm:"not null"`
}
