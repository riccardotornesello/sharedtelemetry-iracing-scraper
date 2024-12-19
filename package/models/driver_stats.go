package models

import (
	"time"

	"gorm.io/gorm"
)

type carCategory string

type DriverStats struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CustID      int    `gorm:"index"`
	CarCategory string `gorm:"index"`
	License     string
	IRating     int
}
