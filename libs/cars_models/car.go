package cars_models

import (
	"time"
)

type Car struct {
	ID *int `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Name            string `json:"name"`
	NameAbbreviated string `json:"nameAbbreviated"`
	Brand           string `json:"brand"`

	Logo        string `json:"logo"`
	SmallImage  string `json:"smallImage"`
	SponsorLogo string `json:"sponsorLogo"`
}
