package firestore

import "encoding/json"

type Driver struct {
	Name     string `firestore:"name"`
	Location string `firestore:"location"`

	Stats map[string]DriverStatsDetails `firestore:"stats"`
}

type DriverStatsDetails struct {
	License string `firestore:"not null"`
	IRating int    `firestore:"not null"`
}

func (c *Driver) ToMap() (res map[string]interface{}) {
	a, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(a, &res)
	return
}
