package irapi

import "encoding/json"

type CarClassResponse struct {
	CarClassId  int `json:"car_class_id"`
	CarsInClass []struct {
		CarDirpath  string `json:"car_dirpath"`
		CarId       int    `json:"car_id"`
		RainEnabled bool   `json:"rain_enabled"`
		Retired     bool   `json:"retired"`
	} `json:"cars_in_class"`
	CustId        int    `json:"cust_id"`
	Name          string `json:"name"`
	RainEnabled   bool   `json:"rain_enabled"`
	RelativeSpeed int    `json:"relative_speed"`
	ShortName     string `json:"short_name"`
}

func (client *IRacingApiClient) GetCarClasses() (*[]CarClassResponse, error) {
	url := "/data/carclass/get"
	respBody, err := client.get(url)
	if err != nil {
		return nil, err
	}

	response := &[]CarClassResponse{}
	err = json.NewDecoder(respBody).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
