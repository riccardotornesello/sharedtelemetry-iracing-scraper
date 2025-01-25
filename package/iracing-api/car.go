package irapi

import "encoding/json"

type CarAssetsResponse struct {
	CarId    int `json:"car_id"`
	CarRules []struct {
		RuleCategory string `json:"rule_category"`
		Text         string `json:"text"`
	} `json:"car_rules"`
	DetailCopy             string `json:"detail_copy"`
	DetailScreenShotImages string `json:"detail_screen_shot_images"`
	DetailTechspecsCopy    string `json:"detail_techspecs_copy"`
	Folder                 string `json:"folder"`
	GalleryImages          string `json:"gallery_images"`
	LargeImage             string `json:"large_image"`
	Logo                   string `json:"logo"`
	SmallImage             string `json:"small_image"`
	SponsorLogo            string `json:"sponsor_logo"`
	TemplatePath           string `json:"template_path"`
}

type CarResponse struct {
	AiEnabled          bool   `json:"ai_enabled"`
	AllowNumberColors  bool   `json:"allow_number_colors"`
	AllowNumberFont    bool   `json:"allow_number_font"`
	AllowSponsor1      bool   `json:"allow_sponsor1"`
	AllowSponsor2      bool   `json:"allow_sponsor2"`
	AllowWheelColor    bool   `json:"allow_wheel_color"`
	AwardExempt        bool   `json:"award_exempt"`
	CarDirpath         string `json:"car_dirpath"`
	CarId              int    `json:"car_id"`
	CarName            string `json:"car_name"`
	CarNameAbbreviated string `json:"car_name_abbreviated"`
	CarTypes           []struct {
		CarType string `json:"car_type"`
	} `json:"car_types"`
	CarWeight               int      `json:"car_weight"`
	Categories              []string `json:"categories"`
	Created                 string   `json:"created"`
	FirstSale               string   `json:"first_sale"`
	ForumUrl                string   `json:"forum_url"`
	FreeWithSubscription    bool     `json:"free_with_subscription"`
	HasHeadlights           bool     `json:"has_headlights"`
	HasMultipleDryTireTypes bool     `json:"has_multiple_dry_tire_types"`
	HasRainCapableTireTypes bool     `json:"has_rain_capable_tire_types"`
	Hp                      int      `json:"hp"`
	IsPsPurchasable         bool     `json:"is_ps_purchasable"`
	MaxPowerAdjustPct       int      `json:"max_power_adjust_pct"`
	MaxWeightPenaltyKg      int      `json:"max_weight_penalty_kg"`
	MinPowerAdjustPct       int      `json:"min_power_adjust_pct"`
	PackageId               int      `json:"package_id"`
	Patterns                int      `json:"patterns"`
	Price                   float32  `json:"price"`
	PriceDisplay            string   `json:"price_display"`
	RainEnabled             bool     `json:"rain_enabled"`
	Retired                 bool     `json:"retired"`
	SearchFilters           string   `json:"search_filters"`
	Sku                     int      `json:"sku"`
}

func (client *IRacingApiClient) GetCarAssets() (*map[string]CarAssetsResponse, error) {
	url := "/data/car/assets"
	respBody, err := client.get(url)
	if err != nil {
		return nil, err
	}

	response := &map[string]CarAssetsResponse{}
	err = json.NewDecoder(respBody).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *IRacingApiClient) GetCars() (*[]CarResponse, error) {
	url := "/data/car/get"
	respBody, err := client.get(url)
	if err != nil {
		return nil, err
	}

	response := &[]CarResponse{}
	err = json.NewDecoder(respBody).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
