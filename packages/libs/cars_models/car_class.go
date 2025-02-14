package cars_models

import (
	"time"
)

type CarClass struct {
	ID *int `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}
