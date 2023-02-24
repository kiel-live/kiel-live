package api

import (
	"encoding/json"
	"net/url"
)

type Platform struct {
	StopPoint string `json:"stopPoint"`
	ShortName string `json:"shortName"`
	Label     string `json:"label"`
}

type platforms struct {
	Platforms []Platform `json:"stopPoints"`
}

func GetPlatforms(stopShortName string) ([]Platform, error) {
	var platforms platforms
	data := url.Values{}
	data.Set("top", "324000000")
	data.Set("bottom", "-324000000")
	data.Set("left", "-648000000")
	data.Set("right", "648000000")

	resp, _ := post(platformsURL, data)
	if err := json.Unmarshal(resp, &platforms); err != nil {
		return nil, err
	}

	var res []Platform
	for _, platform := range platforms.Platforms {
		if platform.ShortName == stopShortName {
			res = append(res, platform)
		}
	}

	return res, nil
}

func GetPlatformDepartures(stopPoint string) (*StopDepartures, error) {
	var platformDepartures StopDepartures
	resp, _ := post(platformURL, url.Values{"stopPoint": {stopPoint}})
	if err := json.Unmarshal(resp, &platformDepartures); err != nil {
		return nil, err
	}
	return &platformDepartures, nil
}
