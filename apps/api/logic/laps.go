package logic

import "github.com/lib/pq"

func IsLapValid(lapNumber int, lapTime int, lapEvents pq.StringArray, incident bool) bool {
	if !(lapNumber > 0 && lapTime > 0 && incident == false) {
		return false
	}

	// Return false if lap.lapEvents contains a blacklisted event
	blacklistedEvents := []string{
		"black flag",
		"car contact",
		"car reset",
		"clock smash",
		"contact",
		"discontinuity",
		"interpolated crossing",
		"invalid",
		"lost control",
		"off track",
		"pitted",
	}

	for _, event := range lapEvents {
		for _, blacklistedEvent := range blacklistedEvents {
			if event == blacklistedEvent {
				return false
			}
		}
	}

	return true
}

func IsLapPitted(lapEvents pq.StringArray) bool {
	for _, event := range lapEvents {
		if event == "pitted" {
			return true
		}
	}

	return false
}
