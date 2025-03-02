package cars_models

import (
	"time"
)

type Brand struct {
	Name string `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Icon string `json:"icon"`
}
