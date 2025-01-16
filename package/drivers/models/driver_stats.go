package models

import (
	"time"

	"gorm.io/gorm"
)

type carCategory string

type DriverStats struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CustID      int    `gorm:"primaryKey;not null"`
	CarCategory string `gorm:"primaryKey;not null"`

	License string `gorm:"not null"`
	IRating int    `gorm:"not null"`
}
