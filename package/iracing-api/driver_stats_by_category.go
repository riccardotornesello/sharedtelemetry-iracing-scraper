package irapi

import (
	"io"
)

func (client *IRacingApiClient) GetDriverStatsByCategoryOval() (io.ReadCloser, error) {
	url := "/data/driver_stats_by_category/oval"
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *IRacingApiClient) GetDriverStatsByCategorySportsCar() (io.ReadCloser, error) {
	url := "/data/driver_stats_by_category/sports_car"
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *IRacingApiClient) GetDriverStatsByCategoryFormulaCar() (io.ReadCloser, error) {
	url := "/data/driver_stats_by_category/formula_car"
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *IRacingApiClient) GetDriverStatsByCategoryRoad() (io.ReadCloser, error) {
	url := "/data/driver_stats_by_category/road"
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *IRacingApiClient) GetDriverStatsByCategoryDirtOval() (io.ReadCloser, error) {
	url := "/data/driver_stats_by_category/dirt_oval"
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *IRacingApiClient) GetDriverStatsByCategoryDirtRoad() (io.ReadCloser, error) {
	url := "/data/driver_stats_by_category/dirt_road"
	body, err := client.get(url)
	if err != nil {
		return nil, err
	}

	return body, nil
}
