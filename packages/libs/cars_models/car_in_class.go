package cars_models

import (
	"time"
)

type CarInClass struct {
	ID *int `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	CarID      int      `json:"carID"`
	CarClassID int      `json:"carClassID"`
	Car        Car      `json:"car" gorm:"foreignKey:CarID"`
	CarClass   CarClass `json:"carClass" gorm:"foreignKey:CarClassID"`
}
