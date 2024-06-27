package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/kiel-live/kiel-live/shared/models"
)

func donkey() ([]*models.Stop, error) {
	stops := make([]*models.Stop, 0)

	// TODO: allow to configure the bounding box
	top := "54.48855"
	left := "9.94689"
	right := "10.30319"
	bottom := "54.19533"
	url := fmt.Sprintf("https://stables.donkey.bike/api/public/nearby?top_right=%s%%2C%s&bottom_left=%s%%2C%s&filter_type=box", top, right, bottom, left)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "donkey/1.0.0")
	req.Header.Set("Accept", "application/com.donkeyrepublic.v7")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	donkeyResp := DonkeyResponse{}
	err = json.Unmarshal(data, &donkeyResp)
	if err != nil {
		return nil, err
	}

	for _, hub := range donkeyResp.Hubs {
		latitude, err := strconv.ParseFloat(hub.Latitude, 32)
		if err != nil {
			return nil, err
		}

		longitude, err := strconv.ParseFloat(hub.Longitude, 32)
		if err != nil {
			return nil, err
		}

		stop := &models.Stop{
			ID:       fmt.Sprintf("donkey-%s", hub.ID),
			Provider: "donkey",
			Name:     hub.Name,
			Type:     "bike-stop",
			Location: &models.Location{
				Latitude:  latitude,
				Longitude: longitude,
				Heading:   0,
			},
		}

		stops = append(stops, stop)
	}

	return stops, nil
}
