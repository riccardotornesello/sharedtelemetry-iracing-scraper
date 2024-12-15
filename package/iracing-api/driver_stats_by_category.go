package irapi

import "log"

func (client *IRacingApiClient) GetDriverStatsByCategoryFormulaCar() []byte {
	url := "/data/driver_stats_by_category/formula_car"
	body, err := client.Get(url)
	if err != nil {
		log.Fatal("Query failed")
	}

	return body
}
