package drivers_models

import (
	"time"

	"gorm.io/gorm"
)

type Driver struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CustID int `gorm:"primarykey"`

	Name     string
	Location string
}
