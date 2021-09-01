package api

import (
	"encoding/json"
	"log"
	"net/url"
)

type tripStops struct {
	Stop       stop   `json:"stop"`
	Status     string `json:"status"`
	ActualTime string `json:"actualTime"`
}

type trip struct {
	Stops         []tripStops `json:"actual"`
	OldStops      []tripStops `json:"old"`
	DirectionText string      `json:"directionText"`
	RouteName     string      `json:"routeName"`
}

func GetTrip(tripID string) trip {
	data := url.Values{}
	data.Set("tripId", tripID)

	resp, _ := post(tripURL, data)
	var trip trip
	if err := json.Unmarshal(resp, &trip); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}
	return trip
}
