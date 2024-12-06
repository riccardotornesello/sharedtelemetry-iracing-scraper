package client

import (
	"encoding/json"
	"strconv"
)

type leagueGetResponse struct {
	LeagueId        int    `json:"league_id"`
	OwnerId         int    `json:"owner_id"`
	LeagueName      string `json:"league_name"`
	Created         string `json:"created"`
	Hidden          bool   `json:"hidden"`
	Message         string `json:"message"`
	About           string `json:"about"`
	Url             string `json:"url"`
	Recruiting      bool   `json:"recruiting"`
	PrivateWall     bool   `json:"private_wall"`
	PrivateRoster   bool   `json:"private_roster"`
	PrivateSchedule bool   `json:"private_schedule"`
	PrivateResults  bool   `json:"private_results"`
	IsOwner         bool   `json:"is_owner"`
	IsAdmin         bool   `json:"is_admin"`
	RosterCount     int    `json:"roster_count"`
	Owner           struct {
		CustId      int    `json:"cust_id"`
		DisplayName string `json:"display_name"`
		Helmet      struct {
			Pattern    int    `json:"pattern"`
			Color1     string `json:"color1"`
			Color2     string `json:"color2"`
			Color3     string `json:"color3"`
			FaceType   int    `json:"face_type"`
			HelmetType int    `json:"helmet_type"`
		} `json:"helmet"`
		CarNumber string `json:"car_number"`
		NickName  string `json:"nick_name"`
	} `json:"owner"`
	Image struct {
		SmallLogo string `json:"small_logo"`
		LargeLogo string `json:"large_logo"`
	} `json:"image"`
	Tags struct {
		Categorized    []string `json:"categorized"`
		NotCategorized []string `json:"not_categorized"`
	} `json:"tags"`
	LeagueApplications []string `json:"league_applications"`
	PendingRequests    []string `json:"pending_requests"`
	IsMember           bool     `json:"is_member"`
	IsApplicant        bool     `json:"is_applicant"`
	IsInvite           bool     `json:"is_invite"`
	IsIgnored          bool     `json:"is_ignored"`
	Roster             []struct {
		CustId      int    `json:"cust_id"`
		DisplayName string `json:"display_name"`
		Helmet      struct {
			Pattern    int    `json:"pattern"`
			Color1     string `json:"color1"`
			Color2     string `json:"color2"`
			Color3     string `json:"color3"`
			FaceType   int    `json:"face_type"`
			HelmetType int    `json:"helmet_type"`
		} `json:"helmet"`
		Licenses []struct {
			CategoryId    int     `json:"category_id"`
			Category      string  `json:"category"`
			CategoryName  string  `json:"category_name"`
			LicenseLevel  int     `json:"license_level"`
			SafetyRating  float32 `json:"safety_rating"`
			Cpi           float32 `json:"cpi"`
			Irating       int     `json:"irating"`
			TtRating      int     `json:"tt_rating"`
			MprNumRaces   int     `json:"mpr_num_races"`
			Color         string  `json:"color"`
			GroupName     string  `json:"group_name"`
			GroupId       int     `json:"group_id"`
			ProPromotable bool    `json:"pro_promotable"`
			Seq           int     `json:"seq"`
			MprNumTts     int     `json:"mpr_num_tts"`
		} `json:"licenses"`
		Owner             bool   `json:"owner"`
		Admin             bool   `json:"admin"`
		LeagueMailOptOut  bool   `json:"league_mail_opt_out"`
		LeaguePmOptOut    bool   `json:"league_pm_opt_out"`
		LeagueMemberSince string `json:"league_member_since"`
		CarNumber         string `json:"car_number"`
		NickName          string `json:"nick_name"`
	} `json:"roster"`
}

type leagueSeasonsResponse struct {
	Subscribed bool `json:"subscribed"`
	Seasons    []struct {
		LeagueId                int    `json:"league_id"`
		SeasonId                int    `json:"season_id"`
		PointsSystemId          int    `json:"points_system_id"`
		SeasonName              string `json:"season_name"`
		Active                  bool   `json:"active"`
		Hidden                  bool   `json:"hidden"`
		NumDrops                int    `json:"num_drops"`
		NoDropsOnOrAfterRaceNum int    `json:"no_drops_on_or_after_race_num"`
		PointsCars              []struct {
			CarId   int    `json:"car_id"`
			CarName string `json:"car_name"`
		} `json:"points_cars"`
		DriverPointsCarClasses []struct {
			CarClassId  int    `json:"car_class_id"`
			Name        string `json:"name"`
			CarsInClass []struct {
				CarId   int    `json:"car_id"`
				CarName string `json:"car_name"`
			} `json:"cars_in_class"`
		} `json:"driver_points_car_classes"`
		TeamPointsCarClasses []struct {
			CarClassId  int    `json:"car_class_id"`
			Name        string `json:"name"`
			CarsInClass []struct {
				CarId   int    `json:"car_id"`
				CarName string `json:"car_name"`
			} `json:"cars_in_class"`
		} `json:"team_points_car_classes"`
		PointsSystemName string `json:"points_system_name"`
		PointsSystemDesc string `json:"points_system_desc"`
	} `json:"seasons"`
	Success  bool `json:"success"`
	Retired  bool `json:"retired"`
	LeagueId int  `json:"league_id"`
}

type LeagueSeasonSessionsResponse struct {
	Sessions []LeagueSeasonSession `json:"sessions"`
}

type LeagueSeasonSession struct {
	Cars []struct {
		CarId        int    `json:"car_id"`
		CarName      string `json:"car_name"`
		CarClassId   int    `json:"car_class_id"`
		CarClassName string `json:"car_class_name"`
	} `json:"cars"`
	DriverChanges     bool   `json:"driver_changes"`
	EntryCount        int    `json:"entry_count"`
	HasResults        bool   `json:"has_results"`
	LaunchAt          string `json:"launch_at"`
	LeagueId          int    `json:"league_id"`
	LeagueSeasonId    int    `json:"league_season_id"`
	LoneQualify       bool   `json:"lone_qualify"`
	PaceCarClassId    int    `json:"pace_car_class_id"`
	PaceCarId         int    `json:"pace_car_id"`
	PasswordProtected bool   `json:"password_protected"`
	PracticeLength    int    `json:"practice_length"`
	PrivateSessionId  int    `json:"private_session_id"`
	QualifyLaps       int    `json:"qualify_laps"`
	QualifyLength     int    `json:"qualify_length"`
	RaceLaps          int    `json:"race_laps"`
	RaceLength        int    `json:"race_length"`
	SessionId         int    `json:"session_id"`
	Status            int    `json:"status"`
	SubsessionId      int    `json:"subsession_id"`
	TeamEntryCount    int    `json:"team_entry_count"`
	TimeLimit         int    `json:"time_limit"`
	Track             struct {
		ConfigName string `json:"config_name"`
		TrackId    int    `json:"track_id"`
		TrackName  string `json:"track_name"`
	} `json:"track"`
	TrackState struct {
		LeaveMarbles         bool `json:"leave_marbles"`
		PracticeGripCompound int  `json:"practice_grip_compound"`
		PracticeRubber       int  `json:"practice_rubber"`
		QualifyGripCompound  int  `json:"qualify_grip_compound"`
		QualifyRubber        int  `json:"qualify_rubber"`
		RaceGripCompound     int  `json:"race_grip_compound"`
		RaceRubber           int  `json:"race_rubber"`
		WarmupGripCompound   int  `json:"warmup_grip_compound"`
		WarmupRubber         int  `json:"warmup_rubber"`
	} `json:"track_state"`
	Weather struct {
		AllowFog       bool `json:"allow_fog"`
		Fog            int  `json:"fog"`
		PrecipOption   int  `json:"precip_option"`
		RelHumidity    int  `json:"rel_humidity"`
		Skies          int  `json:"skies"`
		TempUnits      int  `json:"temp_units"`
		TempValue      int  `json:"temp_value"`
		TrackWater     int  `json:"track_water"`
		Type           int  `json:"type"`
		Version        int  `json:"version"`
		WeatherSummary struct {
			MaxPrecipRateDesc string `json:"max_precip_rate_desc"`
			PrecipChance      int    `json:"precip_chance"`
		} `json:"weather_summary"`
		WeatherVarInitial int `json:"weather_var_initial"`
		WeatherVarOngoing int `json:"weather_var_ongoing"`
		WindDir           int `json:"wind_dir"`
		WindUnits         int `json:"wind_units"`
		WindValue         int `json:"wind_value"`
	} `json:"weather"`
	WinnerId   int    `json:"winner_id"`
	WinnerName string `json:"winner_name"`
}

func (client *IRacingApiClient) GetLeague(leagueId int, include_licenses bool) (*leagueGetResponse, error) {
	url := "/data/league/get?league_id=" + strconv.Itoa(leagueId) + "&include_licenses=" + strconv.FormatBool(include_licenses)
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	response := &leagueGetResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *IRacingApiClient) GetLeagueSeasons(leagueId int, retired bool) (*leagueSeasonsResponse, error) {
	url := "/data/league/seasons?league_id=" + strconv.Itoa(leagueId) + "&retired=" + strconv.FormatBool(retired)
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	response := &leagueSeasonsResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *IRacingApiClient) GetLeagueSeasonSessions(leagueId int, seasonId int, resultsOnly bool) (*LeagueSeasonSessionsResponse, error) {
	resultsOnlyStr := "false"
	if resultsOnly {
		resultsOnlyStr = "true"
	}

	url := "/data/league/season_sessions?league_id=" + strconv.Itoa(leagueId) + "&season_id=" + strconv.Itoa(seasonId) + "&results_only=" + resultsOnlyStr
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	response := &LeagueSeasonSessionsResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
