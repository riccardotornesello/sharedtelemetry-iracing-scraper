package irapi

import "encoding/json"

type TrackAssetsResponse struct {
	Coordinates         string `json:"coordinates"`
	DetailCopy          string `json:"detail_copy"`
	DetailTechspecsCopy string `json:"detail_techspecs_copy"`
	DetailVideo         string `json:"detail_video"`
	Folder              string `json:"folder"`
	GalleryImages       string `json:"gallery_images"`
	GalleryPrefix       string `json:"gallery_prefix"`
	LargeImage          string `json:"large_image"`
	Logo                string `json:"logo"`
	North               string `json:"north"`
	NumSvgImages        int    `json:"num_svg_images"`
	SmallImage          string `json:"small_image"`
	TrackId             int    `json:"track_id"`
	TrackMap            string `json:"track_map"`
	TrackMapLayers      struct {
		Background  string `json:"background"`
		Inactive    string `json:"inactive"`
		Active      string `json:"active"`
		Pitroad     string `json:"pitroad"`
		StartFinish string `json:"start-finish"`
		Turns       string `json:"turns"`
	} `json:"track_map_layers"`
}

type TrackResponse struct {
	AiEnabled              bool    `json:"ai_enabled"`
	AllowPitlaneCollisions bool    `json:"allow_pitlane_collisions"`
	AllowRollingStart      bool    `json:"allow_rolling_start"`
	AllowStandingStart     bool    `json:"allow_standing_start"`
	AwardExempt            bool    `json:"award_exempt"`
	Category               string  `json:"category"`
	CategoryId             int     `json:"category_id"`
	Closes                 string  `json:"closes"`
	ConfigName             string  `json:"config_name"`
	CornersPerLap          int     `json:"corners_per_lap"`
	Created                string  `json:"created"`
	FirstSale              string  `json:"first_sale"`
	FreeWithSubscription   bool    `json:"free_with_subscription"`
	FullyLit               bool    `json:"fully_lit"`
	GridStalls             int     `json:"grid_stalls"`
	HasOptPath             bool    `json:"has_opt_path"`
	HasShortParadeLap      bool    `json:"has_short_parade_lap"`
	HasStartZone           bool    `json:"has_start_zone"`
	HasSvgMap              bool    `json:"has_svg_map"`
	IsDirt                 bool    `json:"is_dirt"`
	IsOval                 bool    `json:"is_oval"`
	IsPsPurchasable        bool    `json:"is_ps_purchasable"`
	LapScoring             int     `json:"lap_scoring"`
	Latitude               float32 `json:"latitude"`
	Location               string  `json:"location"`
	Longitude              float32 `json:"longitude"`
	MaxCars                int     `json:"max_cars"`
	NightLighting          bool    `json:"night_lighting"`
	NominalLapTime         float32 `json:"nominal_lap_time"`
	NumberPitstalls        int     `json:"number_pitstalls"`
	Opens                  string  `json:"opens"`
	PackageId              int     `json:"package_id"`
	PitRoadSpeedLimit      int     `json:"pit_road_speed_limit"`
	Price                  int     `json:"price"`
	PriceDisplay           string  `json:"price_display"`
	Priority               int     `json:"priority"`
	Purchasable            bool    `json:"purchasable"`
	QualifyLaps            int     `json:"qualify_laps"`
	RainEnabled            bool    `json:"rain_enabled"`
	RestartOnLeft          bool    `json:"restart_on_left"`
	Retired                bool    `json:"retired"`
	SearchFilters          string  `json:"search_filters"`
	SiteUrl                string  `json:"site_url"`
	Sku                    int     `json:"sku"`
	SoloLaps               int     `json:"solo_laps"`
	StartOnLeft            bool    `json:"start_on_left"`
	SupportsGripCompound   bool    `json:"supports_grip_compound"`
	TechTrack              bool    `json:"tech_track"`
	TimeZone               string  `json:"time_zone"`
	TrackConfigLength      float32 `json:"track_config_length"`
	TrackDirpath           string  `json:"track_dirpath"`
	TrackId                int     `json:"track_id"`
	TrackName              string  `json:"track_name"`
	TrackTypes             []struct {
		TrackType string `json:"track_type"`
	} `json:"track_types"`
}

func (client *IRacingApiClient) GetTrackAssets() (*map[string]TrackAssetsResponse, error) {
	url := "/data/track/assets"
	respBody, err := client.get(url)
	if err != nil {
		return nil, err
	}

	response := &map[string]TrackAssetsResponse{}
	err = json.NewDecoder(respBody).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *IRacingApiClient) GetTracks() (*[]TrackResponse, error) {
	url := "/data/track/get"
	respBody, err := client.get(url)
	if err != nil {
		return nil, err
	}

	response := &[]TrackResponse{}
	err = json.NewDecoder(respBody).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
